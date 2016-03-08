package sdk

import (
	"bytes"
	"encoding/base64"
	"net/smtp"
	"strings"
)

type MailSdk struct {
	Addr    string
	To      []string
	From    string
	Subject string
	Type    string
	Message []byte
}

func (this *MailSdk) Auth(password string) smtp.Auth {
	identity := ""
	host := strings.Split(this.Addr, ":")[0]
	return smtp.PlainAuth(identity, this.From, password, host)
}

func (this *MailSdk) SendMail(auth smtp.Auth) error {
	//content type
	code := "UTF-8"
	var content_type string
	if this.Type == "html" {
		content_type = "Content-Type: text/" + this.Type + "; charset=UTF-8"
	} else if this.Type == "excel" {
		content_type = "Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet; charset=" + code + "\r\nContent-Disposition: attachment;filename=" + this.Subject + ".xlsx" + "\r\nContent-Description: File Transfer" + "\r\nContent-Transfer-Encoding: base64"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	//mime
	toString := ""
	trimSign := ","
	for _, value := range this.To {
		toString += value + trimSign
	}
	toString = strings.Trim(toString, trimSign)
	msg := []byte("To: " + toString + "\r\nFrom: " + this.From + "<" + this.From + ">\r\nSubject: " + this.Subject + "\r\n" + content_type + "\r\n\r\n")

	//encode body
	msgs := [][]byte{}
	msgs = append(msgs, msg)
	msgLen := len(this.Message)
	finalLen := base64.StdEncoding.EncodedLen(msgLen)
	encodeMsg := make([]byte, finalLen)
	base64.StdEncoding.Encode(encodeMsg, this.Message)

	//final message
	msgs = append(msgs, encodeMsg)
	msg = bytes.Join(msgs, []byte(""))
	return smtp.SendMail(this.Addr, auth, this.From, this.To, msg)
}
