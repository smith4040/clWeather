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
	station := "KJFK"
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

	switch dataType {
	case "metar":
		metarRaw, metarErr = aviation.FetchMETAR(ctx, station)
	case "taf":
		tafRaw, tafErr = aviation.FetchTAF(ctx, station)
	default:
		metarRaw, metarErr = aviation.FetchMETAR(ctx, station)
		tafRaw, tafErr = aviation.FetchTAF(ctx, station)
	}

	if verbose {
		if metarErr == nil {
			fmt.Println("=== METAR ===")
			fmt.Println(string(metarRaw))
		}
		if tafErr == nil {
			fmt.Println("=== TAF ===")
			fmt.Println(string(tafRaw))
		}
		return nil
	}

	var metar aviation.Response
	var taf aviation.Response

	if metarErr == nil {
		metar, metarErr = aviation.ParseMETAR(metarRaw)
		if metarErr != nil {
			return fmt.Errorf("parse METAR: %w", metarErr)
		}
	}
	if tafErr == nil {
		taf, tafErr = aviation.ParseTAF(tafRaw)
		if tafErr != nil {
			return fmt.Errorf("parse TAF: %w", tafErr)
		}
	}

	format.PrintAviation(metar, taf, format.OutputFormat(output))
	return nil
}
