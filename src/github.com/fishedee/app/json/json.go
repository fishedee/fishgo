package json

import (
	"bytes"
	"encoding/json"
)

var (
	jsonQuickTag *quickTag
)

func JsonMarshal(data interface{}) ([]byte, error) {
	quickTagInstance := jsonQuickTag.getTagTypeInstance(data)

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	err := encoder.Encode(quickTagInstance)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func JsonUnmarshal(in []byte, data interface{}) error {
	quickTagInstance := jsonQuickTag.getTagTypeInstance(data)

	err := json.Unmarshal(in, quickTagInstance)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	jsonQuickTag = newQuickTag("json")
}
