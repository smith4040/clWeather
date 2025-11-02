package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/smith4040/clWeather/internal/aviation"
	"github.com/smith4040/clWeather/internal/format"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) error {
	station := "KJFK" // Default to KJFK for demo
	if len(args) > 0 {
		station = strings.ToUpper(args[0])
	}

	dataType, _ := cmd.Flags().GetString("type")
	output, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var metarRaw, tafRaw []byte
	var metarErr, tafErr error
	var metar aviation.Response
	var taf aviation.Response

	switch dataType {
	case "metar":
		metarRaw, metarErr = aviation.FetchMETAR(ctx, station)
		if metarErr != nil {
			return fmt.Errorf("failed to fetch METAR: %w", metarErr)
		}
	case "taf":
		tafRaw, tafErr = aviation.FetchTAF(ctx, station)
		if tafErr != nil {
			return fmt.Errorf("failed to fetch TAF: %w", tafErr)
		}
	default: // both
		metarRaw, metarErr = aviation.FetchMETAR(ctx, station)
		tafRaw, tafErr = aviation.FetchTAF(ctx, station)
		if metarErr != nil && tafErr != nil {
			return fmt.Errorf("both METAR and TAF failed: %v; %v", metarErr, tafErr)
		}
	}

	if verbose {
		if dataType == "metar" || dataType == "both" {
			if metarErr == nil {
				fmt.Println("METAR:", string(metarRaw))
			}
		}
		if dataType == "taf" || dataType == "both" {
			if tafErr == nil {
				fmt.Println("TAF:", string(tafRaw))
			}
		}
		return nil
	}

	if dataType == "metar" || dataType == "both" {
		if metarErr == nil {
			metar, metarErr = aviation.ParseMETAR(metarRaw)
			if metarErr != nil {
				return fmt.Errorf("failed to parse METAR: %w", metarErr)
			}
		}
	}
	if dataType == "taf" || dataType == "both" {
		if tafErr == nil {
			taf, tafErr = aviation.ParseTAF(tafRaw)
			if tafErr != nil {
				return fmt.Errorf("failed to parse TAF: %w", tafErr)
			}
		}
	}

	format.PrintAviation(metar, taf, format.OutputFormat(output))
	return nil
}
