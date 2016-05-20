package crypto

import "crypto/rand"

var (
	randStr = []byte("0123456789abcdefghijklmnopqrstuvwxyz")
)

func generateRand(size int, targetLength byte) string {
	result := make([]byte, size)
	rand.Read(result)
	for singleIndex, singleByte := range result {
		result[singleIndex] = randStr[singleByte%targetLength]
	}
	return string(result)
}

func CryptoRand(size int) string {
	return generateRand(size, byte(len(randStr)))
}

func CryptoRandDigit(size int) string {
	return generateRand(size, 10)
}
