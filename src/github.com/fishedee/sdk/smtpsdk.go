package sdk

import (
	"encoding/base64"
	"net/smtp"
	"strings"
)

type SmtpSdk struct {
	Host string
}

type SmtpSdkMail struct {
	From    string
	To      []string
	Subject string
	Body    string
	Attach  []SmtpSdkMailAttach
}

type SmtpSdkMailAttach struct {
	Name string
	Data []byte
}

func (this *SmtpSdk) Auth(username, password string) smtp.Auth {
	identity := ""
	host := strings.Split(this.Host, ":")[0]
	return smtp.PlainAuth(identity, username, password, host)
}

func (this *SmtpSdk) Send(auth smtp.Auth, mail SmtpSdkMail) error {
	//boundary
	boundary := "XX-JDED00099ASCI"
	assignBoundary := "\"" + boundary + "\""
	middleBoundary := "--" + boundary
	endBoundary := middleBoundary + "--"

	//header
	toString := ""
	trimSign := ","
	for _, value := range mail.To {
		toString += value + trimSign
	}
	toString = strings.Trim(toString, trimSign)
	msg := []byte("To: " + toString + "\r\nFrom: " + mail.From + "<" + mail.From + ">\r\nSubject: " + mail.Subject + "\r\n" + "Content-Type: multipart/mixed; boundary=" + assignBoundary + "\r\n\r\n")

	//body
	msg = append(msg, []byte(middleBoundary+"\r\n"+"Content-Type:text/plain; charset=UTF-8\r\n\r\n")...)
	msg = append(msg, []byte(mail.Body+"\r\n\r\n")...)

	//attach
	for _, value := range mail.Attach {
		msg = append(msg, []byte(middleBoundary+"\r\n"+"Content-Type: application/octet-stream; charset=UTF-8\r\nContent-Disposition: attachment;filename="+value.Name+"\r\nContent-Description: File Transfer"+"\r\nContent-Transfer-Encoding: base64\r\n\r\n")...)
		afterData := this.base64Encode(value.Data)
		msg = append(msg, afterData...)
		msg = append(msg, []byte("\r\n\r\n")...)
	}

	//结束标志
	msg = append(msg, []byte(endBoundary+"\r\n")...)

	//send message
	return smtp.SendMail(this.Host, auth, mail.From, mail.To, msg)
}

func (this *SmtpSdk) base64Encode(data []byte) []byte {
	singleDataLen := len(data)
	afterDataLen := base64.StdEncoding.EncodedLen(singleDataLen)
	afterData := make([]byte, afterDataLen)
	base64.StdEncoding.Encode(afterData, data)
	return afterData
}
