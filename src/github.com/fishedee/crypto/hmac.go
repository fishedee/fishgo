package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

func CryptoHMACMd5(data, key []byte) string {
	mac := hmac.New(md5.New, key)
	mac.Write(data)
	etag := mac.Sum(nil)
	return hex.EncodeToString(etag)
}
