package crypto

import (
	"crypto/sha512"
	"encoding/hex"
)

func Sha3(inValue string) string {
	// Create a new SHA-512 hasher
	hasher := sha512.New()

	// Write the input data to the hasher
	hasher.Write([]byte(inValue))

	// Get the hash sum
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string for display
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
