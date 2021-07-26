package main

import "fmt"

var (
	// info = teal
	warn = yellow
	fata = red
)

var (
	red    = Color("\033[1;31m%s\033[0m")
	green  = Color("\033[1;32m%s\033[0m")
	yellow = Color("\033[1;33m%s\033[0m")
	teal   = Color("\033[1;36m%s\033[0m")
)

// Color function returns a colored string
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString, fmt.Sprint(args...))
	}
	return sprint
}
