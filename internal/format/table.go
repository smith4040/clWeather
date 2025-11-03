package format

import (
	"fmt"
	"os"

	"github.com/rodaine/table"
	"github.com/smith4040/clWeather/internal/aviation"
)

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

	tbl := table.New("PIREP", "", "", "", "", "FltCat", "", "SIGMET").
		WithHeaderFormatter(func(format string, vals ...interface{}) string {
			return fmt.Sprintf("\033[1m%s\033[0m", fmt.Sprintf(format, vals...))
		})

	// Row 1
	row1 := []interface{}{"Turb", "", "", "Ice", "", "", "", ""}
	tbl.AddRow(row1...)

	// Row 2 – highlight active
	row2 := []interface{}{"MOD", "SEV", "LLWS", "MOD", "SEV", "LIFR", "IFR", "VFR"}
	for i, v := range row2 {
		s := fmt.Sprint(v)
		if i == activeCol {
			s = fmt.Sprintf("\033[1;31m%s\033[0m", s) // bold red
		}
		row2[i] = s
	}
	tbl.AddRow(row2...)

	// Print
	tbl.WithWriter(os.Stdout).Print()
	fmt.Println()
}
