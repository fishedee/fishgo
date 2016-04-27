package encoding

import (
	"encoding/base64"
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
