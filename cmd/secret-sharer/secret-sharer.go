package main

import (
	"github.com/SteveNewson/secret-sharer/internal"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "secret-sharer",
	Short:   "Tools for sharing secrets over potentially insecure channels",
	Version: internal.Version + " (" + internal.Build + ")",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
