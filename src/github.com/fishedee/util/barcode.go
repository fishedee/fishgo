// 二维码生成库
package util

import (
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// 二维码结构
type QrCode struct {
	content string // 内容
	width   int    // 宽
	height  int    // 高
}

// 初始化二维码结构
func NewQrCode(content string, width, height int) *QrCode {
	return &QrCode{
		content: content,
		width:   width,
		height:  height,
	}
}

/*
 * GetQRCode：生成图片二维码
 * fileName： 文件名称
 * img：图片格式为：png、jpg/jpeg、gif
 */
func (this *QrCode) GetQRCode() (barcode.Barcode, error) {
	var result barcode.Barcode

	// 将链接编码为二维码
	code, err := qr.Encode(this.content, qr.L, qr.Unicode)
	if err != nil {
		return result, err
	}

	// 重写大小
	code, err = barcode.Scale(code, this.width, this.height)
	if err != nil {
		return result, err
	}

	return code, err
}

/*
 * GetQrCodeFile：生成图片二维码文件
 */
func (this *QrCode) GetQrCodeFile() (string, error) {
	qrCode, err := this.GetQRCode()
	if err != nil {
		return "", err
	}

	fileName, err := CreateTempFile("qrcode", ".png")
	if err != nil {
		return fileName, err
	}
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fileName, err
	}
	defer file.Close()

	err = png.Encode(file, qrCode)
	if err != nil {
		panic(err)
	}
	return fileName, err
}
