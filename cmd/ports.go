package cmd

import (
	"fmt"

	"speedy-cli/internal/portcheck"

	"github.com/spf13/cobra"
)

var port int

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "Check used ports and suggest fixes",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, usages := portcheck.Check(port)
		if jsonOut {
			return printJSON(map[string]any{"result": result, "ports": usages})
		}

		if len(usages) == 0 {
			success("✔ No active listening ports found")
			return nil
		}
		for _, u := range usages {
			if u.InUse {
				fail("✖ Port %d in use by process %d (%s)", u.Port, u.PID, u.Process)
				warn("💡 Suggestion: kill -9 %d", u.PID)
			} else {
				success("✔ Port %d is free", u.Port)
			}
		}
		if verbose {
			fmt.Printf("Status: %s | %s\n", result.Status, result.Message)
		}
		return nil
	},
}

func init() {
	portsCmd.Flags().IntVar(&port, "port", 0, "specific port to check")
	rootCmd.AddCommand(portsCmd)
}
