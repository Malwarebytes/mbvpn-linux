package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type MsgType int

const (
	MsgOutput  MsgType = iota // General output
	MsgSuccess                // Success messages
	MsgError                  // Error messages
	MsgWarning                // Warnings
)

// PrintMsg prints a message to the user with appropriate formatting.
func PrintMsg(message string, msgType MsgType) {
	switch msgType {
	case MsgSuccess:
		fmt.Fprintln(os.Stdout, "✔ "+message)
	case MsgError:
		color.New(color.FgRed).Fprintln(os.Stderr, "✖ "+message)
	case MsgWarning:
		color.New(color.FgYellow).Fprintln(os.Stdout, "⚠ "+message)
	default:
		fmt.Fprintln(os.Stdout, message)
	}
}

// PrintTable prints tabular data to the user.
func PrintTable(headers []string, rows [][]string) {
	// Print headers
	color.New(color.FgCyan).Fprintln(os.Stdout, strings.Join(headers, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(os.Stdout, strings.Join(row, "\t"))
	}
}

// LogError logs an error message for debugging purposes.
func LogError(err error) {
	color.New(color.FgRed).Fprintln(os.Stderr, "Error:", err)
}
