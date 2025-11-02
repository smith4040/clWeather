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

	var data []byte
	var err error
	switch dataType {
	case "metar":
		data, err = aviation.FetchMETAR(ctx, station)
	case "taf":
		data, err = aviation.FetchTAF(ctx, station)
	default:
		// Fetch both
		metar, _ := aviation.FetchMETAR(ctx, station)
		taf, _ := aviation.FetchTAF(ctx, station)
		data = append(metar, taf...) // Simple concat for verbose; parse separately below
		if err != nil {
			return fmt.Errorf("failed to fetch data: %w", err)
		}
	}

	if verbose {
		fmt.Println(string(data))
		return nil
	}

	var metarData, tafData aviation.Response
	if dataType == "both" || dataType == "metar" {
		metarData, err = aviation.ParseMETAR(data) // Reuse data for metar if both
		if err != nil {
			return fmt.Errorf("failed to parse METAR: %w", err)
		}
	}
	if dataType == "both" || dataType == "taf" {
		tafData, err = aviation.ParseTAF(data) // Reuse for taf if both
		if err != nil {
			return fmt.Errorf("failed to parse TAF: %w", err)
		}
	}

	format.PrintAviation(metarData, tafData, format.OutputFormat(output))
	return nil
}
