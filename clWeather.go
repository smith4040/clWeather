package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type response struct {
	Context    string     `json:"@context"`
	ID         string     `json:"id"`
	NWSType    string     `json:"type"`
	Geometry   string     `json:"geometry"`
	Properties properties `json:"properties"`
}

type properties struct {
	ID         string `json:"@id"`
	NWSType    string `json:"@type"`
	RawMessage string `json:"rawMessage"`
}

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

	var responseObject response
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.ID)
	fmt.Println(responseObject.Properties.RawMessage)
}

func main() {
	requestWeather()
}
