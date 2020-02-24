package util_session

import (
	"github.com/yvasiyarov/php_session_decoder"
)

func EncodePhp(data map[interface{}]interface{}) ([]byte, error) {
	data2 := make(php_session_decoder.PhpSession)
	for key, value := range data {
		data2[key.(string)] = value
	}
	encoder := php_session_decoder.NewPhpEncoder(data2)
	result, err := encoder.Encode()
	return []byte(result), err
}

func DecodePhp(data []byte) (map[interface{}]interface{}, error) {
	decoder := php_session_decoder.NewPhpDecoder(string(data))
	data2, err := decoder.Decode()
	result := map[interface{}]interface{}{}
	for key, value := range data2 {
		result[key] = value
	}
	return result, err
}
