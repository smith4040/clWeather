package format

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/smith4040/clWeather/internal/aviation"
)

type OutputFormat string

const (
	Human OutputFormat = "human"
	Raw   OutputFormat = "raw"
	JSON  OutputFormat = "json"
)

func PrintAviation(metar, taf aviation.Response, format OutputFormat) {
	switch format {
	case JSON:
		out := map[string]interface{}{}
		if metar.StationID != "" {
			out["metar"] = metar
		}
		if taf.StationID != "" {
			out["taf"] = taf
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(out)
	case Raw:
		if metar.RawText != "" {
			fmt.Println("METAR:", metar.RawText)
		}
		if taf.RawText != "" {
			fmt.Println("TAF:", taf.RawText)
		}
	default: // human
		if metar.StationID != "" {
			printHumanMETAR(metar)
		}
		if taf.StationID != "" {
			printHumanTAF(taf)
		}

		if metar.StationID != "" {
			fmt.Println(strings.Repeat("─", 80))
			PrintFlightCategoryTable(metar)
		}
	}
}

func printHumanMETAR(r aviation.Response) {
	fmt.Printf("🛩️  %s METAR (Observed: %s) | FltCat: %s\n", r.StationID, r.Time, r.FltCat)
	fmt.Printf("🌡️   Temp/Dew: %s / %s\n", r.Temperature, r.Dewpoint)
	fmt.Printf("💨   Wind: %s%s %s (gust %s)\n", r.WindDir, r.WindSpeed, r.WindDir, r.WindGust)
	fmt.Printf("👁️   Vis: %s | Alt: %s\n", r.Visibility, r.Barometer)
	fmt.Printf("☁️   Clouds: %s\n", cloudSummary(r.Clouds))
	if len(r.Weather) > 0 {
		fmt.Printf("🌩️   Weather: %s\n", strings.Join(r.Weather, ", "))
	}
	fmt.Println()
}

func printHumanTAF(r aviation.Response) {
	fmt.Printf("📅  %s TAF (Issued: %s | Valid: %s to %s)\n", r.StationID, r.Time, r.ValidFrom, r.ValidTo)
	for _, fp := range r.Forecast {
		prefix := "   "
		if fp.Type == "tempo" {
			prefix = "   TEMPO "
		}
		fmt.Printf("%s%s-%s: Wind %s%s | Vis %s\n", prefix, fp.StartTime, fp.EndTime, fp.WindDir, fp.WindSpeed, fp.Visibility)
		fmt.Printf("      Clouds: %s\n", cloudSummary(fp.Clouds))
		if len(fp.Weather) > 0 {
			fmt.Printf("      Weather: %s\n", strings.Join(fp.Weather, ", "))
		}
	}
	fmt.Println()
}

func cloudSummary(clouds []aviation.Cloud) string {
	if len(clouds) == 0 {
		return "CLR"
	}
	var s []string
	for _, c := range clouds {
		s = append(s, fmt.Sprintf("%s%03d", c.Type, c.Base/100))
	}
	return strings.Join(s, " ")
}
