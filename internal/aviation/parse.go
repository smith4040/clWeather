package aviation

import (
	"fmt"
	"time"

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
	Barometer   string
	FltCat      string
	// TAF
	ValidFrom string
	ValidTo   string
	Forecast  []ForecastPeriod
}

type Cloud struct {
	Type string
	Base int
}

type ForecastPeriod struct {
	Type        string
	Start       string
	End         string
	WindDir     string
	WindSpeed   string
	WindGust    string
	Visibility  string
	Clouds      []Cloud
	Weather     []string
	Probability string
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
		WindDir:     fmt.Sprintf("%d", int(gjson.Get(jsonStr, path+".wdir").Float())),
		WindSpeed:   fmt.Sprintf("%dKT", int(gjson.Get(jsonStr, path+".wspd").Float())),
		Visibility:  formatVis(gjson.Get(jsonStr, path+".visib").String()),
		Barometer:   fmt.Sprintf("%.2finHg", gjson.Get(jsonStr, path+".altim").Float()/33.8639),
		FltCat:      gjson.Get(jsonStr, path+".fltCat").String(),
	}

	if gust := gjson.Get(jsonStr, path+".wgst").Float(); gust > 0 {
		r.WindGust = fmt.Sprintf("%dKT", int(gust))
	}

	gjson.Get(jsonStr, path+".clouds.#").ForEach(func(_, v gjson.Result) bool {
		base := int(v.Get("base").Float())
		if base == 0 {
			base = 100 // Avoid 000
		}
		r.Clouds = append(r.Clouds, Cloud{
			Type: v.Get("cover").String(),
			Base: base,
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

	r := Response{
		RawText:   gjson.Get(jsonStr, path+".rawTAF").String(),
		StationID: gjson.Get(jsonStr, path+".icaoId").String(),
		Time:      formatISO(gjson.Get(jsonStr, path+".issueTime").String()),
		ValidFrom: formatUnixTime(gjson.Get(jsonStr, path+".validTimeFrom").Int()),
		ValidTo:   formatUnixTime(gjson.Get(jsonStr, path+".validTimeTo").Int()),
	}

	gjson.Get(jsonStr, path+".fcsts.#").ForEach(func(_, f gjson.Result) bool {
		fp := ForecastPeriod{
			Type:        f.Get("fcstChange").String(),
			Probability: fmt.Sprintf("%d", int(f.Get("probability").Float())),
			Start:       formatUnixTime(f.Get("timeFrom").Int()),
			End:         formatUnixTime(f.Get("timeTo").Int()),
			Visibility:  formatVis(f.Get("visib").String()),
		}

		wdir := f.Get("wdir").String()
		if wdir == "VRB" {
			fp.WindDir = "VRB"
		} else if wdir != "" {
			fp.WindDir = wdir
		}

		wspd := int(f.Get("wspd").Float())
		if wspd > 0 {
			fp.WindSpeed = fmt.Sprintf("%dKT", wspd)
		}

		wgst := int(f.Get("wgst").Float())
		if wgst > 0 {
			fp.WindGust = fmt.Sprintf("%dKT", wgst)
		}

		if wx := f.Get("wxString").String(); wx != "" {
			fp.Weather = append(fp.Weather, wx)
		}

		f.Get("clouds.#").ForEach(func(_, c gjson.Result) bool {
			base := int(c.Get("base").Float())
			if base == 0 {
				base = 100
			}
			fp.Clouds = append(fp.Clouds, Cloud{
				Type: c.Get("cover").String(),
				Base: base,
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

func formatUnixTime(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format("2006-01-02T15:04Z")
}

func formatISO(s string) string {
	if s == "" {
		return ""
	}
	t, _ := time.Parse(time.RFC3339, s)
	return t.UTC().Format("2006-01-02T15:04Z")
}

func formatVis(v string) string {
	if v == "" {
		return "10+SM"
	}
	if v == "6+" {
		return "10+SM"
	}
	return v + "SM"
}
