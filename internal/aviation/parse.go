// internal/aviation/parse.go
package aviation

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

type Response struct {
	RawText     string
	StationID   string
	Time        string // human-readable or Unix
	Temperature string
	Dewpoint    string
	WindDir     string
	WindSpeed   string
	WindGust    string
	Visibility  string
	Clouds      []Cloud
	Barometer   string
	FltCat      string
	// TAF
	ValidFrom string
	ValidTo   string
	Forecast  []ForecastPeriod
}

type Cloud struct {
	Type string // FEW, SCT, BKN, OVC
	Base int    // feet
}

type ForecastPeriod struct {
	Type       string // FM, TEMPO, BECMG, PROB
	Start      string
	End        string
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
	path := "0"
	if !gjson.Get(jsonStr, path).Exists() {
		return Response{}, fmt.Errorf("no METAR data")
	}

	r := Response{
		RawText:     gjson.Get(jsonStr, path+".rawOb").String(),
		StationID:   gjson.Get(jsonStr, path+".icaoId").String(),
		Time:        formatUnixTime(gjson.Get(jsonStr, path+".obsTime").Int()),
		Temperature: fmt.Sprintf("%.1f°C", gjson.Get(jsonStr, path+".temp").Float()),
		Dewpoint:    fmt.Sprintf("%.1f°C", gjson.Get(jsonStr, path+".dewp").Float()),
		WindDir:     gjson.Get(jsonStr, path+".wdir").String(),
		WindSpeed:   fmt.Sprintf("%dKT", int(gjson.Get(jsonStr, path+".wspd").Float())),
		Visibility:  formatVis(gjson.Get(jsonStr, path+".visib").String()),
		Barometer:   fmt.Sprintf("%.2finHg", gjson.Get(jsonStr, path+".altim").Float()/33.8639),
		FltCat:      gjson.Get(jsonStr, path+".fltCat").String(),
	}

	// Clouds
	gjson.Get(jsonStr, path+".clouds.#").ForEach(func(_, v gjson.Result) bool {
		r.Clouds = append(r.Clouds, Cloud{
			Type: v.Get("cover").String(),
			Base: int(v.Get("base").Float()),
		})
		return true
	})

	if r.StationID == "" {
		r.StationID = "UNKNOWN"
	}
	return r, nil
}

func ParseTAF(data []byte) (Response, error) {
	jsonStr := string(data)
	if !gjson.Valid(jsonStr) {
		return Response{}, fmt.Errorf("invalid JSON")
	}
	path := "0"
	if !gjson.Get(jsonStr, path).Exists() {
		return Response{}, fmt.Errorf("no TAF data")
	}

	issueTime := gjson.Get(jsonStr, path+".issueTime").String()
	validFrom := gjson.Get(jsonStr, path+".validTimeFrom").Int()
	validTo := gjson.Get(jsonStr, path+".validTimeTo").Int()

	r := Response{
		RawText:   gjson.Get(jsonStr, path+".rawTAF").String(),
		StationID: gjson.Get(jsonStr, path+".icaoId").String(),
		Time:      formatISO(issueTime), // e.g., 2025-11-02T02:21:00.000Z
		ValidFrom: formatUnixTime(validFrom),
		ValidTo:   formatUnixTime(validTo),
	}

	// Parse each forecast period
	gjson.Get(jsonStr, path+".fcsts.#").ForEach(func(_, f gjson.Result) bool {
		fp := ForecastPeriod{
			Type:       f.Get("fcstChange").String(), // FM, TEMPO, etc.
			Start:      formatUnixTime(f.Get("timeFrom").Int()),
			End:        formatUnixTime(f.Get("timeTo").Int()),
			WindDir:    f.Get("wdir").String(),
			WindSpeed:  fmt.Sprintf("%dKT", int(f.Get("wspd").Float())),
			Visibility: formatVis(f.Get("visib").String()),
		}

		// Clouds
		f.Get("clouds.#").ForEach(func(_, c gjson.Result) bool {
			fp.Clouds = append(fp.Clouds, Cloud{
				Type: c.Get("cover").String(),
				Base: int(c.Get("base").Float()),
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

// ---------------------------------------------------------------
// HELPERS
// ---------------------------------------------------------------
func formatUnixTime(ts int64) string {
	if ts == 0 {
		return "N/A"
	}
	return time.Unix(ts, 0).UTC().Format("2006-01-02T15:04Z")
}

func formatISO(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	return t.UTC().Format("2006-01-02T15:04Z")
}

func formatVis(v string) string {
	if v == "" || v == "6+" {
		return "10+SM"
	}
	return v + "SM"
}
