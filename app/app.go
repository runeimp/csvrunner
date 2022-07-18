package app

import (
	"fmt"
	"os"
)

const (
	Name    = "CSV Runner"
	Label   = Name + " v" + Version
	Version = "1.0.0"
)

var Debug bool

func PrintDebug(msg string, parts ...any) {
	if Debug {
		msg = "DEBUG: " + msg
		PrintError(msg, parts...)
	}
}

func PrintError(msg string, parts ...any) {
	fmt.Fprintf(os.Stderr, msg, parts...)
	fmt.Fprintln(os.Stderr)
}
