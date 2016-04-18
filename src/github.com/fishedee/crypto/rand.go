package crypto

import (
	"crypto/rand"
	"encoding/hex"
)

func CryptoRand(size int) string {
	result := make([]byte, size/2+1)
	rand.Read(result)
	resultStr := hex.EncodeToString(result)
	return resultStr[0:size]
}
