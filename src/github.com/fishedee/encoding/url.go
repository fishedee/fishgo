package encoding

import (
	"net/url"
)

func EncodeUrl(data string) (string, error) {
	return url.QueryEscape(data), nil
}

func DecodeUrl(data string) (string, error) {
	return url.QueryUnescape(data)
}
