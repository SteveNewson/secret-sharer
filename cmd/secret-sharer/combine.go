package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/SteveNewson/secret-sharer/internal"
	"github.com/spf13/cobra"
	"os"
)

var combineCommand = &cobra.Command{
	Use:   "combine",
	Short: "Combine a secret that was split using Shamir's secret sharing",
	Long: `Combines a secret from its component parts.

This 
- echo -e "9d32f670afce9fbeedde38e5\n82e06c97ccfbc70cd224f6a0\n3718d0fedf92f7f2e12e1343" \
  | secret-sharer combine`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return CheckRequiredFlags(cmd.Flags())
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		components, err := convertHexComponents(readHexComponentsFromStdin())
		if err != nil {
			return err
		}

		return internal.Combine(components)
	},
}

func init() {
	rootCmd.AddCommand(combineCommand)
}

func readHexComponentsFromStdin() []string {
	var hexComponents []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		hexComponents = append(hexComponents, scanner.Text())
	}
	return hexComponents
}

func convertHexComponents(hexComponents []string) ([][]byte, error) {
	var components [][]byte
	for _, hexPart := range hexComponents {
		b, err := hex.DecodeString(hexPart)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %q: %w", hexPart, err)
		}

		components = append(components, b)
	}

	return components, nil
}
