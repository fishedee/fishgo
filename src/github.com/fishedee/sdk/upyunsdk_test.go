package sdk

import (
	. "github.com/fishedee/assert"
	"testing"
)

func TestUpyunSdkPut(t *testing.T) {
	sdk := &UpyunSdk{
		Bucket:   "image-fish",
		Operator: "fishedee",
		Password: "fish123456",
	}
	err := sdk.PutString("/test1", []byte("Hello Data1"))
	AssertEqual(t, err, nil)
}
