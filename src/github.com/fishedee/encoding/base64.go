package encoding

import (
	"encoding/base64"
)

func EncodeBase64(in string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(in)), nil
}

func DecodeBase64(in string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}
	return (string)(result), nil
}
