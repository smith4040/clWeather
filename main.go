// clWeather is a Golang command line tool for querying the aviationweather.gov
// API for current weather information for a specified station/stations

package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	datamodel "github.com/smith4040/clWeather/datamodel"
)

// makeURL builds the URL for the endpoint to query
func makeURL(s string, t string) string {
	url := "https://aviationweather.gov/api/data/metar?ids=" + s + "&format=json&taf=" + t + "&hours=1.5"

	return url
}

// celsiusToFahrenheit converts celsius to fahrenheit
func celsiusToFahrenheit(c float64) float64 {
	value := (c * 9 / 5) + 32
	return value
}

// processData prepares JSON
func processData(d []byte) (datamodel.AvWxResponse, error) {
	var station datamodel.AvWxResponse
	err := json.Unmarshal(d, &station)
	if err != nil {
		return station, err
	}
	return station, nil
}

func requestObservation(stationID string, taf string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	url := makeURL(stationID, taf)

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error setting up cookie jar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper), // Disable HTTP/2
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading 400 response body: %v\n", err)
			return
		}

		// Attempt to unmarshal the JSON error response
		var errorResponse datamodel.ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			fmt.Printf("Error unmarshaling 400 error JSON: %v\n", err)
			fmt.Printf("Raw 400 response body: %s\n", body) // Log raw body if unmarshaling fails
			return
		}

		fmt.Printf("Received 400 Bad Request error:\n")
		fmt.Printf("  Status: %s\n", errorResponse.Status)
		fmt.Printf("  Error: %s\n", errorResponse.Error)
	} else if resp.StatusCode != http.StatusOK {
		// Handle other non-200 status codes as needed
		fmt.Printf("Received unexpected status code: %d\n", resp.StatusCode)
	} else {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading data: %s\n", err)
		}

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		ch <- string(responseData)
	}
}

// presentResults is called to display the weather on command line
func presentResults(stations []string, taf string) {
	ch := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(len(stations))

	for _, s := range stations {
		go requestObservation(s, taf, ch, &wg)
	}

	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		p, err := processData([]byte(result))
		if err != nil {
			log.Printf("error processing data: %s\n", err)
			return
		}

		name := p[0].Name
		if name != "" {
			fmt.Println(name)
		}

		metar := p[0].RawOb
		fCat := p[0].FltCat
		if metar != "" {
			switch {
			case fCat == "VFR":
				fmt.Println(green(metar))
			case fCat == "MVFR":
				fmt.Println(blue(metar))
			case fCat == "IFR":
				fmt.Println(red(metar))
			case fCat == "LIFR":
				fmt.Println(magenta(metar))
			default:
				fmt.Println(metar)
			}

		}

		if taf == "true" {
			rt := p[0].RawTaf
			if rt != "" {
				fmt.Println(rt)
			}
		}

		if p[0].Temp != nil {
			t := p[0].Temp
			f := celsiusToFahrenheit(*t)
			s := fmt.Sprintf("%.2f", f)
			fmt.Println(teal("Temperature is " + s + "°F"))
		} else {
			fmt.Println(yellow("Temperature is currently unavailable, please try again later."))
		}

		if p[0].Wdir != nil {
			fmt.Println(teal("Wind is from ", p[0].Wdir, "\u00B0 at ", p[0].Wspd, "kts\n"))
		} else {
			println("\n")
		}
	}
	// this isnt true, fix
	fmt.Println("All requests complete.")
}

// printUsage displays flags to the user if none are presented
func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	start := time.Now()
	var station string

	flag.StringVar(&station, "s", "", "Station ID")
	taf := flag.Bool("t", false, "Show Terminal Area Forecast")
	flag.Parse()
	if flag.NFlag() == 0 {
		printUsage() // if user does not supply flags, print usage
	}

	str := strconv.FormatBool(*taf)

	stations := strings.Split(station, ",")
	fmt.Printf("Searching station(s): %s\n\n", stations)
	presentResults(stations, str)

	defer func() {
		fmt.Println("Execution Time: ", time.Since(start))
	}()
}
