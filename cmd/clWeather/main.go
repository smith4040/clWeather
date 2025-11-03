package main

import (
	"os"

	"github.com/smith4040/clWeather/internal/cli"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "clweather [station]",
		Short: "Aviation METAR/TAF in your terminal",
		Long:  `A fast CLI for METAR and TAF data from aviationweather.gov.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  cli.Run,
	}

	root.PersistentFlags().StringP("type", "t", "both", "Data type: metar|taf|both")
	root.PersistentFlags().StringP("output", "o", "human", "Output format: human|raw|json")
	root.PersistentFlags().BoolP("verbose", "v", false, "Show full API response")

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
