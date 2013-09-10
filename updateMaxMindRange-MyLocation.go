package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	maxLat = 180
	maxLon = 360
)

var (
	c = [][]string{}
)

func initc() {
	c = make([][]string, maxLat+1)
	for i := range c {
		c[i] = make([]string, maxLon+1)
		for j := range c[i] {
			c[i][j] = ""
		}
	}
}

func main() {
	initc()
	getNewLocationData()
	updateRangeLocation()
}

func updateRangeLocation() {

	filename := os.Args[1:]
	filenamer := "maxmind_city_ipv4_blocks.csv"
	filenamew := "maxmind_range_ipv4_location.csv"

	if len(filename) == 1 {
		filenamer = filename[0]
	} else if len(filename) == 2 {
		filenamer = filename[0]
		filenamew = filename[1]
	}

	f1, errr := os.Open(filenamer)
	f2, errw := os.Create(filenamew)
	defer f1.Close()
	defer f2.Close()
	if errr != nil {
		panic(errr)
	}
	if errw != nil {
		panic(errw)
	}

	r := csv.NewReader(f1)
	var record []string
	var start, end, locId string
	var added bool
	n := 0
	for {
		record, errr = r.Read()
		if errr == io.EOF {
			break
		}
		if len(record) == 0 {
			continue
		}

		start = record[0]
		end = record[1]
		locId = record[2]

		added = false
		for i := 0; i <= maxLat && !added; i++ {
			for j := 0; j <= maxLon && !added; j++ {
				if c[i][j] == "" {
					continue
				} else {
					split := strings.Split(c[i][j], ",")
					for _, s := range split {
						if s == "" {
							continue
						}
						if s == locId {
							fmt.Fprintf(f2, "%s,%s,%d\n", start, end, i*1000+j)
							added = true
						}
						if added {
							break
						}
					}
				}
			}
		}
		if !added {
			fmt.Printf("ERROR: Range not added %d,%d\n", start, end)
			return
		}
		n++
		fmt.Printf("Added %d record\n", n)
	}
}
func getNewLocationData() {
	// Open stream to Blocks CSV db
	f, err := os.Open("GeoLiteCity-Location.csv")
	defer f.Close()
	if err != nil {
		panic(err)
	}

	// Create CSV reader
	r := csv.NewReader(f)
	r.TrailingComma = true
	r.FieldsPerRecord = -1

	var record []string
	var locId string
	var myLat, myLon int
	var lat, lon float64
	n := 0

	for {
		// Read until EOF
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		if len(record) == 0 {
			continue
		}
		// parse locId
		locId = record[0]
		// parse Lat
		lat, err = strconv.ParseFloat(record[5], 64)
		if err != nil {
			fmt.Println("Error parsing lat %v\n", err)
			continue
		}
		myLat = int(lat + maxLat/2)

		// parse Lon
		lon, err = strconv.ParseFloat(record[6], 64)
		if err != nil {
			fmt.Println("Error parsing lon %v\n", err)
			continue
		}
		myLon = int(lon + maxLon/2)
		if myLat > 180 || myLon > 360 || myLat < 0 || myLon < 0 {
			fmt.Printf("geolocation out of rannge %d, %d\n", myLat, myLon)
		}
		c[myLat][myLon] += locId + ","
		n++

	}
	fmt.Printf("Done compressing %d locations\n", n)
	getCount()
}

func getCount() {

	count := 0
	for i := 0; i <= maxLat; i++ {
		for j := 0; j <= maxLon; j++ {
			if c[i][j] == "" {
				continue
			} else {
				count++
			}
		}
	}
	fmt.Printf("There are %d locations\n", count)
}
