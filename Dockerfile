FROM golang:1.10.0 AS builder
WORKDIR /go/src/github.com/wheresalice/influx-trains
ADD . .
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o app

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/wheresalice/influx-trains/app .
ENTRYPOINT [ "/app" ]