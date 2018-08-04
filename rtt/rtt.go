// To parse and unparse this JSON data, add this code to your project and do:
//
//    station, err := UnmarshalStation(bytes)
//    bytes, err = station.Marshal()

package rtt

import "encoding/json"

func UnmarshalStation(data []byte) (Station, error) {
	var r Station
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Station) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Station struct {
	Location Location    `json:"location"`
	Filter   interface{} `json:"filter"`
	Services []Service   `json:"services"`
}

type Location struct {
	Name   string `json:"name"`
	CRS    string `json:"crs"`
	Tiploc string `json:"tiploc"`
}

type Service struct {
	LocationDetail  LocationDetail `json:"locationDetail"`
	ServiceUid      string         `json:"serviceUid"`
	RunDate         string         `json:"runDate"`
	TrainIdentity   string         `json:"trainIdentity"`
	RunningIdentity string         `json:"runningIdentity"`
	AtocCode        AtocCode       `json:"atocCode"`
	AtocName        AtocName       `json:"atocName"`
	ServiceType     ServiceType    `json:"serviceType"`
	IsPassenger     bool           `json:"isPassenger"`
	Origin          []Destination  `json:"origin"`
	Destination     []Destination  `json:"destination"`
}

type Destination struct {
	Tiploc      string `json:"tiploc"`
	Description string `json:"description"`
	WorkingTime string `json:"workingTime"`
	PublicTime  string `json:"publicTime"`
}

type LocationDetail struct {
	RealtimeActivated          bool          `json:"realtimeActivated"`
	Tiploc                     Tiploc        `json:"tiploc"`
	CRS                        CRS           `json:"crs"`
	Description                Description   `json:"description"`
	GbttBookedArrival          *string       `json:"gbttBookedArrival,omitempty"`
	GbttBookedArrivalNextDay   *bool         `json:"gbttBookedArrivalNextDay,omitempty"`
	GbttBookedDeparture        string        `json:"gbttBookedDeparture"`
	GbttBookedDepartureNextDay *bool         `json:"gbttBookedDepartureNextDay,omitempty"`
	Origin                     []Destination `json:"origin"`
	Destination                []Destination `json:"destination"`
	IsCall                     bool          `json:"isCall"`
	IsPublicCall               bool          `json:"isPublicCall"`
	RealtimeArrival            *string       `json:"realtimeArrival,omitempty"`
	RealtimeArrivalActual      *bool         `json:"realtimeArrivalActual,omitempty"`
	RealtimeArrivalNextDay     *bool         `json:"realtimeArrivalNextDay,omitempty"`
	RealtimeDeparture          string        `json:"realtimeDeparture"`
	RealtimeDepartureActual    bool          `json:"realtimeDepartureActual"`
	RealtimeDepartureNextDay   *bool         `json:"realtimeDepartureNextDay,omitempty"`
	Platform                   string        `json:"platform"`
	PlatformConfirmed          bool          `json:"platformConfirmed"`
	PlatformChanged            bool          `json:"platformChanged"`
	DisplayAs                  DisplayAs     `json:"displayAs"`
	Associations               []Association `json:"associations"`
	CancelReasonCode           *string       `json:"cancelReasonCode,omitempty"`
	CancelReasonShortText      *string       `json:"cancelReasonShortText,omitempty"`
	CancelReasonLongText       *string       `json:"cancelReasonLongText,omitempty"`
	ServiceLocation            *string       `json:"serviceLocation,omitempty"`
}

type Association struct {
	Type              Type   `json:"type"`
	AssociatedUid     string `json:"associatedUid"`
	AssociatedRunDate string `json:"associatedRunDate"`
}

type AtocCode string
const (
	Em AtocCode = "EM"
	Gr AtocCode = "GR"
	NT AtocCode = "NT"
	Tp AtocCode = "TP"
	Xc AtocCode = "XC"
)

type AtocName string
const (
	CrossCountry AtocName = "CrossCountry"
	EastMidlandsTrains AtocName = "East Midlands Trains"
	FirstTranspennineExpress AtocName = "First Transpennine Express"
	Northern AtocName = "Northern"
	VirginTrainsEastCoast AtocName = "Virgin Trains (East Coast)"
)

type Type string
const (
	Divide Type = "divide"
	Join Type = "join"
	Next Type = "next"
)

type CRS string
const (
	Lds CRS = "LDS"
)

type Description string
const (
	Leeds Description = "Leeds"
)

type DisplayAs string
const (
	Call DisplayAs = "CALL"
	CancelledCall DisplayAs = "CANCELLED_CALL"
	Origin DisplayAs = "ORIGIN"
	Terminates DisplayAs = "TERMINATES"
)

type Tiploc string
const (
	TiplocLEEDS Tiploc = "LEEDS"
)

type ServiceType string
const (
	Train ServiceType = "train"
)
