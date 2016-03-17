package crypto

import (
	"crypto/sha1"
	"encoding/hex"
)

func CryptoSha1(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	etag := hash.Sum(nil)
	etagString := hex.EncodeToString(etag)
	return etagString
}
