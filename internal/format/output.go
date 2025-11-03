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

func printHumanMETAR(r aviation.Response) {
	fmt.Printf("METAR %s (Observed: %s) | FltCat: %s\n", r.StationID, r.Time, r.FltCat)
	fmt.Printf("   Temp/Dew: %s / %s\n", r.Temperature, r.Dewpoint)

	wind := r.WindDir + r.WindSpeed
	if r.WindGust != "" {
		wind += "G" + r.WindGust
	}
	fmt.Printf("   Wind: %s\n", wind)

	fmt.Printf("   Vis: %s | Alt: %s\n", r.Visibility, r.Barometer)
	fmt.Printf("   Clouds: %s\n", cloudSummary(r.Clouds))
	fmt.Println()
}

func printHumanTAF(r aviation.Response) {
	fmt.Printf("TAF %s (Issued: %s | Valid: %s to %s)\n", r.StationID, r.Time, r.ValidFrom, r.ValidTo)
	for _, fp := range r.Forecast {
		if fp.Start == "N/A" {
			continue
		}

		prefix := ""
		if fp.Probability != "" {
			prefix = "PROB" + fp.Probability + " "
		} else if fp.Type != "" {
			prefix = fp.Type + " "
		}

		end := ""
		if fp.End != "N/A" {
			end = " → " + fp.End
		}

		wind := ""
		if fp.WindDir != "" && fp.WindSpeed != "" {
			wind = fp.WindDir + fp.WindSpeed
			if fp.WindGust != "" {
				wind += "G" + fp.WindGust
			}
		} else if fp.WindDir == "VRB" && fp.WindSpeed != "" {
			wind = "VRB" + fp.WindSpeed
		}

		fmt.Printf("  %s%s%s: Wind %s | Vis %s\n", prefix, fp.Start, end, wind, fp.Visibility)

		if len(fp.Weather) > 0 {
			fmt.Printf("     Weather: %s\n", strings.Join(fp.Weather, ", "))
		}
		if len(fp.Clouds) > 0 {
			fmt.Printf("     Clouds: %s\n", cloudSummary(fp.Clouds))
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
		base := c.Base / 100
		if base == 0 {
			base = 1 // avoid 000
		}
		s = append(s, fmt.Sprintf("%s%03d", c.Type, base))
	}
	return strings.Join(s, " ")
}

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
