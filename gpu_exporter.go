package main

import (
	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "gpu_exporter",
		Short:   "GPU exporter collect node nvidia gpu metrics",
		Long:    "GPU exporter collect node all nvidia gpu metrics expose to Prometheus",
		Version: version,
	}
	rootCmd.AddCommand(
		newCollectorCommand(),
	)
	rootCmd.Execute()
}
