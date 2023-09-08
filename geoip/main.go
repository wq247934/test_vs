package main

import (
	"fmt"
	"github.com/ip2location/ip2location-go"
)

func main() {
	db, err := ip2location.OpenDB("./IP-COUNTRY-REGION-CITY-SAMPLE.BIN")
	if err != nil {
		fmt.Print(err)
		return
	}
	ip := "49.232.13.131"
	results, err := db.Get_all(ip)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("ip: %s\n", ip)
	fmt.Printf("country_short: %s\n", results.Country_short)
	fmt.Printf("country_long: %s\n", results.Country_long)
	fmt.Printf("region: %s\n", results.Region)
	fmt.Printf("city: %s\n", results.City)
	db.Close()
}
