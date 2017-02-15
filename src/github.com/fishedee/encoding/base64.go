package encoding

import (
	"encoding/base64"
	"strings"
)

func EncodeBase64(in []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(in), nil
}

func DecodeBase64(in string) ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// url安全的base64编码
func EncodeBase64Safe(in []byte) (string, error) {
	encData := base64.StdEncoding.EncodeToString(in)
	encData = strings.Replace(encData, `+`, `-`, -1)
	encData = strings.Replace(encData, `/`, `_`, -1)
	return encData, nil
}
