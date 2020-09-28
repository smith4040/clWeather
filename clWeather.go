package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/smith4040/clWeather/dataModel"
)

func requestWeather() {
	var stationID string = "kfwb"
	url := "https://api.weather.gov/stations/" + (stationID) + "/observations/latest"

	weatherResponse, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(weatherResponse.Body)

	if err != nil {
		log.Fatal(err)
	}

	var responseObject dataModel.Response
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.ID)
	fmt.Println(responseObject.Properties.RawMessage)
}

func main() {
	requestWeather()
}
