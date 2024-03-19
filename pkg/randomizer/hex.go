package randomizer

import (
	"encoding/hex"
	"fmt"
)

func GenerateRandomHex(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be even")
	}
	bytes, err := GenerateRandomBytes(length / 2)
	if err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
