package cmd

import (
	"fmt"

	"speedy-cli/internal/stocks"

	"github.com/spf13/cobra"
)

var (
	symbol string
	topN   int
)

var stocksCmd = &cobra.Command{
	Use:   "stocks",
	Short: "Show NEPSE stock prices and updates",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, movers := stocks.Get(symbol, topN)
		if jsonOut {
			return printJSON(map[string]any{"result": result, "stocks": movers})
		}

		info("🏦 NEPSE Movers:")
		for i, m := range movers {
			arrow := "↑"
			line := success
			if m.ChangePct < 0 {
				arrow = "↓"
				line = fail
			}
			line("%d. %-6s %s %.2f%%  NPR %.2f %s", i+1, m.Symbol, arrow, abs(m.ChangePct), m.Price, m.Spark)
		}
		if result.Suggestion != "" {
			warn("💡 %s", result.Suggestion)
		}
		if verbose {
			fmt.Printf("Source: %s\n", result.Message)
		}
		return nil
	},
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func init() {
	stocksCmd.Flags().StringVar(&symbol, "symbol", "", "specific NEPSE symbol")
	stocksCmd.Flags().IntVar(&topN, "top", 5, "number of movers to show")
	rootCmd.AddCommand(stocksCmd)
}
