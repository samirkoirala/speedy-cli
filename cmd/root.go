package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	jsonOut bool
)

var rootCmd = &cobra.Command{
	Use:   "speedy-cli",
	Short: "Fast and fun terminal diagnostics tool",
	Long:  "speedy-cli checks speed, ports, env vars, and NEPSE stocks with colorful output.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show additional debug details")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output in JSON format")
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func info(msg string, args ...any) {
	color.New(color.FgCyan, color.Bold).Printf(msg+"\n", args...)
}

func success(msg string, args ...any) {
	color.New(color.FgGreen).Printf(msg+"\n", args...)
}

func warn(msg string, args ...any) {
	color.New(color.FgYellow).Printf(msg+"\n", args...)
}

func fail(msg string, args ...any) {
	color.New(color.FgRed).Printf(msg+"\n", args...)
}
