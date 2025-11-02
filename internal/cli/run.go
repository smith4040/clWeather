package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/smith4040/clWeather/internal/aviation"
	"github.com/smith4040/clWeather/internal/format"
	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) error {
	// ---- 1. Resolve station (default KJFK) ----
	station := "KJFK"
	if len(args) > 0 {
		station = strings.ToUpper(args[0])
	}

	// ---- 2. Flags ----
	dataType, _ := cmd.Flags().GetString("type") // metar | taf | both
	output, _ := cmd.Flags().GetString("output") // human | raw | json
	verbose, _ := cmd.Flags().GetBool("verbose")

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	// ---- 3. Fetch raw payloads ----
	var (
		metarRaw []byte
		tafRaw   []byte
		metarErr error
		tafErr   error
	)

	switch dataType {
	case "metar":
		metarRaw, metarErr = aviation.FetchMETAR(ctx, station)
	case "taf":
		tafRaw, tafErr = aviation.FetchTAF(ctx, station)
	default: // both
		metarRaw, metarErr = aviation.FetchMETAR(ctx, station)
		tafRaw, tafErr = aviation.FetchTAF(ctx, station)
	}

	// ---- 4. Verbose mode – dump raw JSON and exit ----
	if verbose {
		if dataType == "metar" || dataType == "both" {
			if metarErr == nil {
				fmt.Println("=== METAR raw JSON ===")
				os.Stdout.Write(metarRaw)
				fmt.Println()
			}
		}
		if dataType == "taf" || dataType == "both" {
			if tafErr == nil {
				fmt.Println("=== TAF raw JSON ===")
				os.Stdout.Write(tafRaw)
				fmt.Println()
			}
		}
		return nil
	}

	// ---- 5. Parse what we successfully fetched ----
	var metar aviation.Response
	var taf aviation.Response

	if dataType == "metar" || dataType == "both" {
		if metarErr != nil {
			fmt.Fprintf(os.Stderr, "METAR fetch failed: %v\n", metarErr)
		} else {
			metar, metarErr = aviation.ParseMETAR(metarRaw)
			if metarErr != nil {
				return fmt.Errorf("failed to parse METAR: %w", metarErr)
			}
		}
	}

	if dataType == "taf" || dataType == "both" {
		if tafErr != nil {
			fmt.Fprintf(os.Stderr, "TAF fetch failed: %v\n", tafErr)
		} else {
			taf, tafErr = aviation.ParseTAF(tafRaw)
			if tafErr != nil {
				return fmt.Errorf("failed to parse TAF: %w", tafErr)
			}
		}
	}

	// ---- 6. Print ----
	format.PrintAviation(metar, taf, format.OutputFormat(output))
	return nil
}
