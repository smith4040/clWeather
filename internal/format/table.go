package format

import (
	"fmt"
	"os"

	"github.com/rodaine/table"
	"github.com/smith4040/clWeather/internal/aviation"
)

func PrintFlightCategoryTable(metar aviation.Response) {
	if metar.StationID == "" {
		return
	}

	// Map flight category to column index (LIFR=5, IFR=6, MVFR=7, VFR=8)
	// We'll highlight the active one
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

	tbl := table.New("PIREP", "", "", "", "", "FltCat", "", "SIGMET").
		WithHeaderFormatter(func(format string, vals ...interface{}) string {
			return fmt.Sprintf("\033[1m%s\033[0m", fmt.Sprintf(format, vals...))
		})

	// Row 1: Turb |        |        | Ice    |        |        |        |        |
	row1 := []interface{}{"Turb", "", "", "Ice", "", "", "", ""}
	tbl.AddRow(row1...)

	// Row 2: MOD  | SEV    | LLWS   | MOD    | SEV    | LIFR   | IFR    |        |
	row2 := []interface{}{"MOD", "SEV", "LLWS", "MOD", "SEV", "LIFR", "IFR", ""}
	for i, v := range row2 {
		s := fmt.Sprint(v)
		if i == activeCol {
			s = fmt.Sprintf("\033[1;31m%s\033[0m", s) // Bold red
		}
		row2[i] = s
	}
	tbl.AddRow(row2...)

	// Row 3: (empty for future use)
	tbl.AddRow("", "", "", "", "", "", "", "")

	// Print to stdout
	tbl.WithWriter(os.Stdout).Print()
	fmt.Println()
}
