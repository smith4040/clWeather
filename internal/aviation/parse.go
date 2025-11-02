package aviation

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// ---------------------------------------------------------------------
// NOTE: The AWC API returns a **plain JSON array** for a single station:
//   [{"raw_text":"KJFK 011200Z ...","station_id":"KJFK",...}]
// ---------------------------------------------------------------------

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
	FltCat      string // VFR / MVFR / IFR / LIFR
	// TAF-specific
	ValidFrom string
	ValidTo   string
	Forecast  []ForecastPeriod
}

type Cloud struct {
	Type string // FEW, SCT, BKN, OVC
	Base int    // feet AGL
}

type ForecastPeriod struct {
	Type       string
	StartTime  string
	EndTime    string
	WindDir    string
	WindSpeed  string
	Visibility string
	Clouds     []Cloud
	Weather    []string
}

// ---------------------------------------------------------------------
// METAR
// ---------------------------------------------------------------------
func ParseMETAR(data []byte) (Response, error) {
	jsonStr := string(data)
	if !gjson.Valid(jsonStr) {
		return Response{}, fmt.Errorf("invalid JSON")
	}

	// First element of the array (index 0)
	path := "0"
	if !gjson.Get(jsonStr, path).Exists() {
		return Response{}, fmt.Errorf("no METAR data")
	}

	r := Response{
		RawText:     gjson.Get(jsonStr, path+".raw_text").String(),
		StationID:   gjson.Get(jsonStr, path+".station_id").String(),
		Time:        gjson.Get(jsonStr, path+".observation_time").String(),
		Temperature: gjson.Get(jsonStr, path+".temp_c").String() + "°C",
		Dewpoint:    gjson.Get(jsonStr, path+".dewp_c").String() + "°C",
		WindDir:     gjson.Get(jsonStr, path+".wdir").String(),
		WindSpeed:   gjson.Get(jsonStr, path+".wspd").String() + "KT",
		WindGust:    gjson.Get(jsonStr, path+".wgust").String() + "KT",
		Visibility:  gjson.Get(jsonStr, path+".visib").String() + "SM",
		Barometer:   gjson.Get(jsonStr, path+".altim").String() + "inHg",
		Elevation:   gjson.Get(jsonStr, path+".elev").String() + "ft",
		FltCat:      gjson.Get(jsonStr, path+".flight_category").String(),
	}

	// Clouds
	gjson.Get(jsonStr, path+".sky_condition.#").ForEach(func(_, v gjson.Result) bool {
		r.Clouds = append(r.Clouds, Cloud{
			Type: v.Get("cover").String(),
			Base: int(v.Get("base_agl").Num),
		})
		return true
	})

	// Weather
	gjson.Get(jsonStr, path+".wx_string").String() // fallback to raw if needed
	// The API doesn't always give an array – just use raw_text for now

	if r.StationID == "" {
		r.StationID = "UNKNOWN"
	}
	return r, nil
}

// ---------------------------------------------------------------------
// TAF
// ---------------------------------------------------------------------
func ParseTAF(data []byte) (Response, error) {
	jsonStr := string(data)
	if !gjson.Valid(jsonStr) {
		return Response{}, fmt.Errorf("invalid JSON")
	}

	path := "0"
	if !gjson.Get(jsonStr, path).Exists() {
		return Response{}, fmt.Errorf("no TAF data")
	}

	r := Response{
		RawText:   gjson.Get(jsonStr, path+".raw_text").String(),
		StationID: gjson.Get(jsonStr, path+".station_id").String(),
		Time:      gjson.Get(jsonStr, path+".issue_time").String(),
		ValidFrom: gjson.Get(jsonStr, path+".valid_from").String(),
		ValidTo:   gjson.Get(jsonStr, path+".valid_to").String(),
	}

	// Forecast periods
	gjson.Get(jsonStr, path+".forecast.#").ForEach(func(_, f gjson.Result) bool {
		fp := ForecastPeriod{
			Type:       f.Get("change_indicator").String(),
			StartTime:  f.Get("time_becoming").String(),
			WindDir:    f.Get("wind_dir_degrees").String(),
			WindSpeed:  f.Get("wind_speed_kt").String() + "KT",
			Visibility: f.Get("visibility_statute_mi").String() + "SM",
		}

		// Clouds in forecast
		f.Get("sky_condition.#").ForEach(func(_, c gjson.Result) bool {
			fp.Clouds = append(fp.Clouds, Cloud{
				Type: c.Get("cover").String(),
				Base: int(c.Get("base_agl").Num),
			})
			return true
		})

		r.Forecast = append(r.Forecast, fp)
		return true
	})

	if r.StationID == "" {
		r.StationID = "UNKNOWN"
	}
	return r, nil
}
