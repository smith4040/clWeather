package format

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rodaine/table"
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
	wind := r.WindDir + r.WindSpeed
	if r.WindGust != "" {
		wind += "G" + r.WindGust
	}
	fmt.Printf("💨   Wind: %s\n", wind)
	fmt.Printf("👁️   Vis: %s | Alt: %s\n", r.Visibility, r.Barometer)
	fmt.Printf("☁️   Clouds: %s\n", cloudSummary(r.Clouds))
	fmt.Println()
}

func printHumanTAF(r aviation.Response) {
	fmt.Printf("📅  %s TAF (Issued: %s | Valid: %s to %s)\n", r.StationID, r.Time, r.ValidFrom, r.ValidTo)
	for _, fp := range r.Forecast {
		if fp.Start == "" {
			continue
		}
		prefix := ""
		if fp.Probability != "" {
			prefix = "PROB" + fp.Probability + " "
		} else if fp.Type != "" {
			prefix = fp.Type + " "
		}
		end := ""
		if fp.End != "" {
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
		s = append(s, fmt.Sprintf("%s%03d", c.Type, base))
	}
	return strings.Join(s, " ")
}

func PrintFlightCategoryTable(metar aviation.Response) {
	activeCol := -1
	switch metar.FltCat {
	case "LIFR":
		activeCol = 5
	case "IFR":
		activeCol = 6
	case "MVFR":
		activeCol = 7
	case "VFR":
		activeCol = 8
	}

	tbl := table.New("PIREP", "", "", "", "", "FltCat", "", "SIGMET")

	row1 := []interface{}{"Turb", "", "", "Ice", "", "", "", ""}
	tbl.AddRow(row1...)

	row2 := []interface{}{"MOD", "SEV", "LLWS", "MOD", "SEV", "LIFR", "IFR", "VFR"}
	for i := range row2 {
		if i == activeCol {
			row2[i] = fmt.Sprintf("\033[1;31m%s\033[0m", row2[i])
		}
	}
	tbl.AddRow(row2...)

	tbl.WithWriter(os.Stdout).Print()
	fmt.Println()
}
