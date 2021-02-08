package internal

import (
	"fmt"
	"github.com/hashicorp/vault/shamir"
)

func Split(secret []byte, number int, threshold int) error {
	components, err := shamir.Split(secret, number, threshold)
	if err != nil {
		return fmt.Errorf("failed to split secret into compnents: %w", err)
	}

	for _, component := range components {
		fmt.Printf("%x\n", component)
	}

	return nil
}

func Combine(components [][]byte) error {
	secretBytes, err := shamir.Combine(components)
	if err != nil {
		return fmt.Errorf("failed to combine secret: %w", err)
	}

	fmt.Printf("%s", string(secretBytes))

	return nil
}
