package cmd

import (
	"fmt"
	"time"

	"speedy-cli/internal/animation"

	"github.com/spf13/cobra"
)

var funCmd = &cobra.Command{
	Use:   "fun-mode",
	Short: "Show fun mini terminal animations",
	Run: func(cmd *cobra.Command, args []string) {
		info("🎉 Fun mode activated")
		animation.RocketLaunch()
		for i := 0; i <= 100; i += 10 {
			fmt.Printf("\rDownload %3d%% %s", i, animation.Progress(i))
			time.Sleep(90 * time.Millisecond)
		}
		fmt.Println()
		for i := 0; i <= 100; i += 20 {
			fmt.Printf("\rUpload   %3d%% %s", i, animation.Progress(i))
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println("\n🚀 Done. Your terminal just got faster vibes.")
	},
}

func init() {
	rootCmd.AddCommand(funCmd)
}
