package main

import (
	"encoding/json"
	"github.com/jasonmoo/lambda_proc"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

func main() {
	lambda_proc.Run(func(context *lambda_proc.Context, eventJSON json.RawMessage) (interface{}, error) {
		var v map[string]interface{}
		if err := json.Unmarshal(eventJSON, &v); err != nil {
			return nil, err
		}
		return getLocation(v["source-ip"].(string))
	})
}

func getLocation(ip string) (*geoip2.City, error) {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	parsedIp := net.ParseIP(ip)
	record, err := db.City(parsedIp)
	if err != nil {
		log.Fatal(err)
	}

	return record, err
}
