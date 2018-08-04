package main

import (
	"log"
	"os"

	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/wheresalice/influx-trains/rtt"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
)

var config struct {
	RttUsername string
	RttPassword string
}

const (
	MyDB     = "trains"
	InfluxDB = "http://influxdb:8086"
)

func main() {
	config.RttPassword = os.Getenv("RTT_PASSWORD")
	config.RttUsername = os.Getenv("RTT_USERNAME")

	//// Get latest
	//latest := getRtt("LDS")
	//batchPoints := rttToPoints(latest, time.Now())
	//writeBatch(batchPoints)

	// Batch import old data
	files, err := ioutil.ReadDir("/data")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		log.Printf("found file: %s", f.Name())
		fileData := fileToRtt(f.Name())
		points := rttToPoints(fileData, filenameToDate(f.Name()))
		writeBatch(points)
	}
}

func filenameToDate(filename string) time.Time {
	filename = strings.TrimPrefix(filename, "LDS-")
	filename = strings.TrimSuffix(filename, ".json")
	log.Printf("date to parse: %s", filename)
	dt, _ := time.Parse("2006-01-02", filename)
	log.Printf("parsed date: %s", dt)
	return dt
}

func writeBatch(bp client.BatchPoints) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: InfluxDB,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

func getRtt(station string) rtt.Station {
	log.Printf("getting rtt data for %s", station)

	var stationData rtt.Station
	rtt_url := fmt.Sprintf("https://%s:%s@api.rtt.io/api/v1/json/search/%s", config.RttUsername, config.RttPassword, station)
	log.Print(rtt_url)
	response, err := http.Get(rtt_url)
	if err != nil {
		log.Fatal(err)
	}
	buf, _ := ioutil.ReadAll(response.Body)
	stationData, err = rtt.UnmarshalStation(buf)
	if err != nil {
		log.Fatalf("failed parsing station data: %v", err)
	}
	return stationData
}

func rttToPoints(station rtt.Station, date time.Time) client.BatchPoints {
	// setup the batchpoint
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	tags := map[string]string{"service": "service"}

	// loop through services

	for _, v := range station.Services {
		// @todo handle delays crossing midnight
		var arrivalDelay int
		if v.LocationDetail.RealtimeArrival == nil || v.LocationDetail.GbttBookedArrival == nil {
			arrivalDelay = 0
		} else {
			realArrival, _ := time.Parse("1504", *v.LocationDetail.RealtimeArrival)
			bookedArrival, _ := time.Parse("1504", *v.LocationDetail.GbttBookedArrival)
			arrivalDelay = int(realArrival.Sub(bookedArrival).Minutes())
			log.Printf("booked arrival: %s, real arrival: %s, arrival delay: %v", bookedArrival, realArrival, arrivalDelay)
		}

		realDeparture, _ := time.Parse("1504", v.LocationDetail.RealtimeDeparture)
		bookedDeparture, _ := time.Parse("1504", v.LocationDetail.GbttBookedDeparture)
		departureDelay := int(realDeparture.Sub(bookedDeparture).Minutes())
		log.Printf("departure delay: %v", departureDelay)

		fields := map[string]interface{}{
			"arrival_delay":     arrivalDelay,
			"departure_delay":   departureDelay,
			"platform":          v.LocationDetail.Platform,
			"origin":            v.LocationDetail.Origin[0].Tiploc,
			"destination":       v.LocationDetail.Destination[0].Tiploc,
			"operator":          v.AtocCode,
			"cancellation_code": v.LocationDetail.CancelReasonCode,
		}

		rttTime, _ := time.Parse("1504", v.LocationDetail.RealtimeDeparture)
		t := date
		y := t.Year()
		mon := t.Month()
		d := t.Day()
		h := rttTime.Hour()
		m := rttTime.Minute()
		logTime := time.Date(y, mon, d, h, m, 0, 0, time.Local)

		pt, err := client.NewPoint("services", tags, fields, logTime)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	return bp
}

func fileToRtt(filename string) rtt.Station {
	jsonFile, err := os.Open("/data/" + filename)
	if err != nil {
		log.Println(err)
	}
	log.Printf("sucessfuly opened %s", filename)
	defer jsonFile.Close()

	buf, _ := ioutil.ReadAll(jsonFile)
	stationData, err := rtt.UnmarshalStation(buf)
	if err != nil {
		log.Fatalf("failed parsing station data %s: %v", filename, err)
	}
	return stationData
}
