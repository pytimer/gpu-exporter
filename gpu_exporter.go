package main

import (
	"github.com/pytimer/gpu_exporter/cmd"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "gpu_exporter",
		Short:   "GPU exporter collect node nvidia gpu metrics",
		Long:    "GPU exporter collect node all nvidia gpu metrics expose to different backends, such as Influxdb/Prometheus",
		Version: "1.0.0",
	}
	rootCmd.AddCommand(
		cmd.NewCollectorCommand(),
	)
	rootCmd.Execute()
}
