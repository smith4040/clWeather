package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	colour "github.com/fatih/color"
	flag "github.com/ogier/pflag"
	dataModel "github.com/smith4040/clWeather/datamodel"
)

func makeURL(s string) string {
	const (
		apiURL          = "https://api.weather.gov/stations/"
		observationType = "/observations/latest"
	)

	url := apiURL + s + observationType
	return url
}

func requestWeather(stationID string) dataModel.Response {
	url := makeURL(stationID)
	fmt.Println(url)
	weatherResponse, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error reading data: %s\n", err)
	}

	defer weatherResponse.Body.Close()

	responseData, err := ioutil.ReadAll(weatherResponse.Body)

	if err != nil {
		log.Fatalf("Error reading data: %s\n", err)
	}

	var station dataModel.Response
	json.Unmarshal(responseData, &station)

	if weatherResponse.StatusCode >= 400 {
		colour.Red(stationID + ":" + " Weather observation for this station is currently unavialable. Check spelling and try again later.")
		log.Print("Response status code: ", weatherResponse.StatusCode)
		return station
	}

	return station
}

var (
	station string
)

func init() {
	flag.StringVarP(&station, "station", "s", "", "Station ID")
}

func main() {
	flag.Parse()

	// if user does not supply flags, print usage
	if flag.NFlag() == 0 {
		printUsage()
	}

	stations := strings.Split(station, ",")
	fmt.Printf("Searching station(s): %s\n", stations)
	fmt.Println("")

	for _, s := range stations {
		result := requestWeather(s)
		colour.HiGreen(result.Properties.RawMessage)
		fmt.Println("")

	}
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
