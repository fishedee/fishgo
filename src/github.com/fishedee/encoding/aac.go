// aac音频编码
package encoding

import (
	"errors"

	fdkaac "github.com/winlinvip/go-fdkaac/fdkaac"
)

// 编码器
type AacEncoder struct {
	aacEncoder *fdkaac.AacEncoder
}

// 新建编码器
// channels: 声道数，如：2
// sampleRate: 采样率，如：44100
// bitrateBps: 码率,比特率，如：48000
func NewAacEncoder(channels, sampleRate, bitrateBps int) (*AacEncoder, error) {
	aacEnc := &AacEncoder{}

	// 新建编码器
	encoder := fdkaac.NewAacEncoder()

	// 初始化
	if err := encoder.InitLc(channels, sampleRate, bitrateBps); err != nil {
		return aacEnc, err
	}

	aacEnc.aacEncoder = encoder

	return aacEnc, nil
}

// 编码
func (this *AacEncoder) Encode(data []byte) ([]byte, error) {
	result := []byte{}
	if this.aacEncoder == nil {
		return result, errors.New("请先初始化encoder！")
	}

	var err error
	if result, err = this.aacEncoder.Encode(data); err != nil {
		return result, err
	}

	return result, nil
}

// 关闭编码器
func (this *AacEncoder) Close() {
	if this.aacEncoder == nil {
		return
	}
	this.aacEncoder.Close()
}

// 冲刷编码器
func (this *AacEncoder) Flush() ([]byte, error) {
	result := []byte{}
	if this.aacEncoder == nil {
		return result, errors.New("请先初始化encoder！")
	}

	return this.aacEncoder.Flush()
}

// 获取编码器的声道数
func (this *AacEncoder) ChannelNum() (int, error) {
	result := 0
	if this.aacEncoder == nil {
		return result, errors.New("请先初始化encoder！")
	}

	return this.aacEncoder.Channels(), nil
}

// 获取编码器的帧大小
func (this *AacEncoder) FrameSize() (int, error) {
	result := 0
	if this.aacEncoder == nil {
		return result, errors.New("请先初始化encoder！")
	}

	return this.aacEncoder.FrameSize(), nil
}

// 获取编码器的每个aac帧的字节数
func (this *AacEncoder) NbBytesPerFrame() (int, error) {
	result := 0
	if this.aacEncoder == nil {
		return result, errors.New("请先初始化encoder！")
	}

	return this.aacEncoder.NbBytesPerFrame(), nil
}

// 解码器
type AacDecoder struct {
	aacDecoder *fdkaac.AacDecoder
}

// 新建解码器
// asc: SequenceHeader
func NewAacDecoder(asc []byte) (*AacDecoder, error) {
	aacDec := &AacDecoder{}

	// 新建解码器
	decoder := fdkaac.NewAacDecoder()

	// 初始化
	if len(asc) != 0 {
		if err := decoder.InitRaw(asc); err != nil {
			return aacDec, err
		}
	} else {
		if err := decoder.InitAdts(); err != nil {
			return aacDec, err
		}
	}

	aacDec.aacDecoder = decoder

	return aacDec, nil
}

// 解码
func (this *AacDecoder) Decode(data []byte) ([]byte, error) {
	result := []byte{}
	if this.aacDecoder == nil {
		return result, errors.New("请先初始化decoder！")
	}

	// 解码
	var err error
	if result, err = this.aacDecoder.Decode(data); err != nil {
		return result, err
	}
	return result, nil
}

// 关闭解码器
func (this *AacDecoder) Close() {
	if this.aacDecoder == nil {
		return
	}
	this.aacDecoder.Close()
}
