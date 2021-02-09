package main

import (
	"bufio"
	"fmt"
	"github.com/SteveNewson/secret-sharer/internal"
	"github.com/SteveNewson/secret-sharer/internal/ansi"
	"github.com/spf13/cobra"
	"os"
)

var (
	secret       string
	transportKey string
	contextInfo  string
)

var wrapCommand = &cobra.Command{
	Use:   "wrap",
	Short: "Encrypts a secret in a supplied transport key",
	Long: `Encrypts a secret in a supplied transport key.

Wrap a secret:
- secret-sharer wrap --secret "hello world" --transport-key "foobar=="

Wrap a secret provided via stdin appending a newline:
- echo "hello world" | secret-sharer wrap --transport-key "foobar=="

Wrap a secret provided via stdin without appending a newline:
- echo -n "hello world" | secret-sharer wrap --transport-key "foobar=="`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return CheckRequiredFlags(cmd.Flags())
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if secret == "" {
			fmt.Println("Enter the secret terminated with <" + ansi.Info("ctrl+]") + "> <" + ansi.Info("enter") + ">")
			reader := bufio.NewReader(os.Stdin)
			secret, _ = reader.ReadString('\x1D')
			fmt.Println("")
		}

		return internal.Wrap(transportKey, secret, contextInfo)
	},
}

func init() {
	wrapCommand.Flags().StringVarP(&secret, "secret", "s", "", "The secret to encrypt (uses stdin if not provided)")
	wrapCommand.Flags().StringVarP(&transportKey, "transport-key", "k", "", "The secret to encrypt (uses stdin if not provided)")
	wrapCommand.Flags().StringVarP(&contextInfo, "context-info", "c", "", "Additional context information to supply the encryption")

	if err := wrapCommand.MarkFlagRequired("transport-key"); err != nil {
		os.Exit(1)
	}

	rootCmd.AddCommand(wrapCommand)
}
