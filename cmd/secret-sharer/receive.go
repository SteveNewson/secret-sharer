package main

import (
	"github.com/SteveNewson/secret-sharer/internal"
	"github.com/spf13/cobra"
)

var (
	numberOfSenders int
)

var receiveCommand = &cobra.Command{
	Use:   "receive",
	Short: "Request and receive a secret",
	Long: `Generates a temporary transport key for exchanging a secret in the open with another party

Receive a single shared secret
- secret-sharer receive

Receive 3 shared secrets
- secret-sharer receive --senders 3`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return CheckRequiredFlags(cmd.Flags())
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Receive(numberOfSenders)
	},
}

func init() {
	receiveCommand.Flags().IntVarP(&numberOfSenders, "senders", "s", 1, "The number of senders")

	rootCmd.AddCommand(receiveCommand)
}
