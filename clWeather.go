package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	flag "github.com/ogier/pflag"
	"github.com/smith4040/clWeather/dataModel"
)

const (
	apiURL          = "https://api.weather.gov"
	stationEndpoint = "/stations/"
	observationType = "/observations/latest"
)

func requestWeather(stationID string) dataModel.Response {
	url := apiURL + stationEndpoint + stationID + observationType

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

	return station
}

var (
	statation string
)

func init() {
	flag.StringVarP(&statation, "station", "s", "", "Station ID")
}
func main() {
	flag.Parse()

	// if user does not supply flags, print usage
	if flag.NFlag() == 0 {
		printUsage()
	}

	stations := strings.Split(statation, ",")
	fmt.Printf("Searching station(s): %s\n", stations)
	fmt.Println("")

	for _, s := range stations {
		result := requestWeather(s)
		fmt.Println(result.Properties.RawMessage)
		fmt.Println("")

	}
}

func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
