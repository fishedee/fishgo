package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func CryptoMd5(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	etag := hash.Sum(nil)
	etagString := hex.EncodeToString(etag)
	return etagString
}
