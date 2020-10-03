# clWeather
clWeather is a command line application that queries the National Weather Service API (https://api.weather.gov) and receives the latest weather observation for the requested station(s).  

## How to Use the Application 
These instructions require your work environment to be setup in accordance with https://golang.org/doc/gopath_code.html

### Installation

* `go get github.com/smith4040/clWeather`
* `go build` or `go install` the application

### Usage

* Run the command `clWeather -s {four letter station identifier}` 

For example: `clWeather -s kord` returns the latest weather observation from Chicago O'Hare International Airport.

Multiple stations can be queried by adding additional stations seperated by commas, i.e.: `clWeather -s kord,kjfk,kbos`

The program also accepts POSIX/GNU-style --flags, i.e.: `--station` in place of `-s`.


## How to Find the Identity of a Station
* Go to https://www.faa.gov/air_traffic/weather/asos/
* Select your state and click "Go"
* Click on one of the pins closest to the location where you want to check the weather
* The four-letter station identifier will be in the top left corner of the view that opens.

The current output of the response is in METAR format. See https://www.weather.gov/media/wrh/mesowest/metar_decode_key.pdf for help reading and decoding the response. The current plan is to modify the output to be more readable/decoded and add a flag for getting the raw response. 

##
This repository is a work in progress. New features will be added soon, and a more complete test suite is in progress. 
