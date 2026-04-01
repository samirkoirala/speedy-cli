package cmd

import (
	"fmt"
	"strings"

	"speedy-cli/internal/animation"
	"speedy-cli/internal/speedtest"

	"github.com/spf13/cobra"
)

var (
	speedServer string
	speedGraph  bool
)

var speedCmd = &cobra.Command{
	Use:   "check-speed",
	Short: "Check internet speed (download/upload/latency)",
	RunE: func(cmd *cobra.Command, args []string) error {
		servers := speedtest.DefaultServers()
		if speedServer != "" {
			servers = strings.Split(speedServer, ",")
		}

		stop := animation.StartSpinner("🚀 Running speed tests")
		result, stats := speedtest.RunParallel(servers, verbose)
		stop()

		if jsonOut {
			return printJSON(map[string]any{"result": result, "stats": stats})
		}

		info("🌐 Speed Test Results:")
		scale := stats.DownloadMbps
		if stats.UploadMbps > scale {
			scale = stats.UploadMbps
		}
		if scale < 10 {
			scale = 10
		}
		success("Download: %.2f Mbps %s", stats.DownloadMbps, animation.Bar(stats.DownloadMbps, scale))
		success("Upload:   %.2f Mbps %s", stats.UploadMbps, animation.Bar(stats.UploadMbps, scale))
		warn("Ping:     %.2f ms ⚡", stats.PingMs)
		if speedGraph {
			fmt.Println(animation.DualGraph("Download", stats.DownloadMbps, "Upload", stats.UploadMbps))
		}
		if verbose {
			for _, s := range stats.ByServer {
				fmt.Printf("- %s: ↓ %.2f Mbps, ↑ %.2f Mbps, ping %.2f ms\n", s.Server, s.DownloadMbps, s.UploadMbps, s.PingMs)
			}
		}
		if result.Suggestion != "" {
			warn("💡 %s", result.Suggestion)
		}
		return nil
	},
}

func init() {
	speedCmd.Flags().StringVar(&speedServer, "server", "", "comma-separated download endpoints")
	speedCmd.Flags().BoolVar(&speedGraph, "graph", false, "show ASCII graph")
	rootCmd.AddCommand(speedCmd)
}
