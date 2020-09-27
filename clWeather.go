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
	weatherResponse, err := http.Get("https://api.weather.gov/stations/kfwb/observations/latest")

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
