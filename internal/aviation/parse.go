package aviation

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type Response struct {
	RawText     string
	StationID   string
	Time        string
	Temperature string
	Dewpoint    string
	WindDir     string
	WindSpeed   string
	WindGust    string
	Visibility  string
	Clouds      []Cloud
	Weather     []string
	Barometer   string
	Elevation   string
	FltCat      string // Flight Category: VFR, IFR, MVFR, LIFR
	// TAF-specific
	ValidFrom string
	ValidTo   string
	Forecast  []ForecastPeriod
}

type Cloud struct {
	Type string // e.g., FEW, SCT, BKN
	Base int    // feet
}

type ForecastPeriod struct {
	Type       string // forecast or tempo
	StartTime  string
	EndTime    string
	WindDir    string
	WindSpeed  string
	Visibility string
	Clouds     []Cloud
	Weather    []string
}

func ParseMETAR(data []byte) (Response, error) {
	jsonStr := string(data)
	if !gjson.Valid(jsonStr) {
		return Response{}, fmt.Errorf("invalid JSON")
	}

	path := "data.METAR.0"
	if !gjson.Get(jsonStr, path).Exists() {
		return Response{}, fmt.Errorf("no METAR data")
	}

	r := Response{
		RawText:     gjson.Get(jsonStr, path+".raw_text").String(),
		StationID:   gjson.Get(jsonStr, path+".station_id").String(),
		Time:        gjson.Get(jsonStr, path+".observation_time").String(),
		Temperature: gjson.Get(jsonStr, path+".temperature").String() + "°C",
		Dewpoint:    gjson.Get(jsonStr, path+".dewpoint").String() + "°C",
		WindDir:     gjson.Get(jsonStr, path+".wind.direction").String(),
		WindSpeed:   gjson.Get(jsonStr, path+".wind.speed").String() + " KT",
		WindGust:    gjson.Get(jsonStr, path+".wind.gust").String() + " KT",
		Visibility:  gjson.Get(jsonStr, path+".visibility.statute_miles").String() + " SM",
		Barometer:   gjson.Get(jsonStr, path+".barometer.value").String() + " inHg",
		Elevation:   gjson.Get(jsonStr, path+".elevation.value").String() + " ft",
		FltCat:      gjson.Get(jsonStr, path+".fltCat").String(),
	}

	// Clouds array
	clouds := gjson.Get(jsonStr, path+".clouds.#")
	for i := 0; i < int(clouds.Num); i++ {
		cPath := fmt.Sprintf("%s.clouds.%d", path, i)
		r.Clouds = append(r.Clouds, Cloud{
			Type: gjson.Get(jsonStr, cPath+".type").String(),
			Base: int(gjson.Get(jsonStr, cPath+".base").Num),
		})
	}

	// Weather array
	weather := gjson.Get(jsonStr, path+".weather.#")
	for i := 0; i < int(weather.Num); i++ {
		wPath := fmt.Sprintf("%s.weather.%d", path, i)
		r.Weather = append(r.Weather, gjson.Get(jsonStr, wPath+".type").String())
	}

	if r.StationID == "" {
		r.StationID = "Unknown"
	}

	return r, nil
}

func ParseTAF(data []byte) (Response, error) {
	jsonStr := string(data)
	if !gjson.Valid(jsonStr) {
		return Response{}, fmt.Errorf("invalid JSON")
	}

	path := "data.TAF.0"
	if !gjson.Get(jsonStr, path).Exists() {
		return Response{}, fmt.Errorf("no TAF data")
	}

	r := Response{
		RawText:   gjson.Get(jsonStr, path+".raw_text").String(),
		StationID: gjson.Get(jsonStr, path+".station_id").String(),
		Time:      gjson.Get(jsonStr, path+".issue_time").String(),
		ValidFrom: gjson.Get(jsonStr, path+".valid_time_from").String(),
		ValidTo:   gjson.Get(jsonStr, path+".valid_time_to").String(),
	}

	// Forecast periods
	forecasts := gjson.Get(jsonStr, path+".forecast.#")
	for i := 0; i < int(forecasts.Num); i++ {
		fPath := fmt.Sprintf("%s.forecast.%d", path, i)
		fp := ForecastPeriod{
			Type:       gjson.Get(jsonStr, fPath+".type").String(),
			StartTime:  gjson.Get(jsonStr, fPath+".start_time").String(),
			EndTime:    gjson.Get(jsonStr, fPath+".end_time").String(),
			WindDir:    gjson.Get(jsonStr, fPath+".wind.direction").String(),
			WindSpeed:  gjson.Get(jsonStr, fPath+".wind.speed").String() + " KT",
			Visibility: gjson.Get(jsonStr, fPath+".visibility.statute_miles").String() + " SM",
		}

		// Clouds in forecast
		fClouds := gjson.Get(jsonStr, fPath+".clouds.#")
		for j := 0; j < int(fClouds.Num); j++ {
			cPath := fmt.Sprintf("%s.clouds.%d", fPath, j)
			fp.Clouds = append(fp.Clouds, Cloud{
				Type: gjson.Get(jsonStr, cPath+".type").String(),
				Base: int(gjson.Get(jsonStr, cPath+".base").Num),
			})
		}

		// Weather in forecast
		fWeather := gjson.Get(jsonStr, fPath+".weather.#")
		for j := 0; j < int(fWeather.Num); j++ {
			wPath := fmt.Sprintf("%s.weather.%d", fPath, j)
			fp.Weather = append(fp.Weather, gjson.Get(jsonStr, wPath+".type").String())
		}

		r.Forecast = append(r.Forecast, fp)
	}

	if r.StationID == "" {
		r.StationID = "Unknown"
	}

	return r, nil
}
