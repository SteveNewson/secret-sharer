package main

import (
	"fmt"
	"github.com/SteveNewson/secret-sharer/internal"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var (
	numberOfComponents int
	threshold          int
)

var splitCommand = &cobra.Command{
	Use:   "split",
	Short: "Split a secret using Shamir's secret sharing",
	Long: `Splits a secret into multiple parts which can be recovered with a subset of the parts.

This will split a secret into 5 parts with 3 parts required to reconstitute it
- echo "mysecret" | secret-sharer split --number 5 --threshold 3`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return CheckRequiredFlags(cmd.Flags())
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		secret, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to read the secret from STDIN: %v", err)
			return err
		}

		return internal.Split(secret, numberOfComponents, threshold)
	},
}

func init() {
	splitCommand.Flags().IntVarP(&numberOfComponents, "number", "n", 0, "The number of components to create")
	splitCommand.Flags().IntVarP(&threshold, "threshold", "t", 0, "The threshold of components required to reconstitute secret")

	if err := splitCommand.MarkFlagRequired("number"); err != nil {
		os.Exit(1)
	}

	if err := splitCommand.MarkFlagRequired("threshold"); err != nil {
		os.Exit(1)
	}

	rootCmd.AddCommand(splitCommand)
}
