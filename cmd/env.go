package cmd

import (
	"speedy-cli/internal/envcheck"

	"github.com/spf13/cobra"
)

var envFile string

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Validate required environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		required := []string{"DB_HOST", "API_KEY", "DB_USER", "DB_PASS"}
		result, checks := envcheck.Check(envFile, required)
		if jsonOut {
			return printJSON(map[string]any{"result": result, "checks": checks})
		}
		for _, c := range checks {
			if c.Present {
				success("✔ %s present", c.Key)
			} else {
				fail("✖ %s missing", c.Key)
			}
		}
		if result.Suggestion != "" {
			warn("💡 %s", result.Suggestion)
		}
		return nil
	},
}

func init() {
	envCmd.Flags().StringVar(&envFile, "file", ".env", "env file path")
	rootCmd.AddCommand(envCmd)
}
