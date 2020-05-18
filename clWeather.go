package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Context    string     `json:"@context"`
	ID         string     `json:"id"`
	NWSType    string     `json:"type"`
	Geometry   string     `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	ID         string `json:"@id"`
	NWSType    string `json:"@type"`
	RawMessage string `json:"rawMessage"`
}

func main() {
	response, err := http.Get("https://api.weather.gov/stations/kfwb/observations/latest")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(responseData))
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	fmt.Println(responseObject.ID)
	fmt.Println(responseObject.Properties.RawMessage)
}

// type Response struct {
//     Name    string    `json:"name"`
//     Pokemon []Pokemon `json:"pokemon_entries"`
// }

// // A Pokemon Struct to map every pokemon to.
// type Pokemon struct {
//     EntryNo int            `json:"entry_number"`
//     Species PokemonSpecies `json:"pokemon_species"`
// }

// // A struct to map our Pokemon's Species which includes it's name
// type PokemonSpecies struct {
//     Name string `json:"name"`
// }
