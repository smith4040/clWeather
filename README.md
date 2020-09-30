# clWeather
clWeather is a command line application that queries the National Weather Service API (https://api.weather.gov) and receives the latest weather observation for the requested station(s).  

## How to Use the Application 
These instructions require your work environment to be setup in accordance with https://golang.org/doc/gopath_code.html

* Clone or fork the repository and then execute the command `go install github.com/{your username}/clWeather`  
* Run the command `clWeather -s {four letter station identifier}` 

For example: `clWeather -s kord` returns the latest weather observation from Chicago O'Hare International Airport.

The program also accepts POSIX/GNU-style --flags, ie: `--station` in place of `-s`.


## How to Find the Identity of a Station
* Go to https://www.faa.gov/air_traffic/weather/asos/
* Select your state and click "Go"
* Click on one of the pins closest to the location where you want to check the weather
* The four-letter station identifier will be in the top left corner of the view that opens.


