package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(inValue string) string {
	data := []byte(inValue)
	hash := md5.Sum(data) // md5.Sum returns a [16]byte array

	// Convert the byte array to a hexadecimal string for display
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
