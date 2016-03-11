package sdk

import (
	"bytes"
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

type SmtpSdkMailAuth struct {
	UserName string
	Password string
}

func (this *SmtpSdk) auth(username, password string) smtp.Auth {
	identity := ""
	host := strings.Split(this.Host, ":")[0]
	return smtp.PlainAuth(identity, username, password, host)
}

func (this *SmtpSdk) Send(mailAuth SmtpSdkMailAuth, mail SmtpSdkMail) error {
	auth := this.auth(mailAuth.UserName, mailAuth.Password)

	msg := bytes.Buffer{}

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
	_, err := msg.WriteString("To: " + toString + "\r\nFrom: " + mail.From + "<" + mail.From + ">\r\nSubject: " + mail.Subject + "\r\n" + "Content-Type: multipart/mixed; boundary=" + assignBoundary + "\r\n\r\n")
	if err != nil {
		return err
	}

	//body
	_, err = msg.WriteString(middleBoundary + "\r\n" + "Content-Type:text/html; charset=UTF-8\r\n\r\n")
	if err != nil {
		return err
	}
	_, err = msg.WriteString(mail.Body + "\r\n\r\n")
	if err != nil {
		return err
	}

	//attach
	for _, value := range mail.Attach {
		_, err = msg.WriteString(middleBoundary + "\r\n" + "Content-Type: application/octet-stream; charset=UTF-8\r\nContent-Disposition: attachment;filename=" + value.Name + "\r\nContent-Description: File Transfer" + "\r\nContent-Transfer-Encoding: base64\r\n\r\n")
		if err != nil {
			return err
		}
		afterData := this.base64Encode(value.Data)
		_, err = msg.Write(afterData)
		if err != nil {
			return err
		}
		_, err = msg.WriteString("\r\n\r\n")
		if err != nil {
			return err
		}
	}

	//结束标志
	_, err = msg.WriteString(endBoundary + "\r\n")
	if err != nil {
		return err
	}

	//send message
	return smtp.SendMail(this.Host, auth, mail.From, mail.To, msg.Bytes())
}

func (this *SmtpSdk) base64Encode(data []byte) []byte {
	singleDataLen := len(data)
	afterDataLen := base64.StdEncoding.EncodedLen(singleDataLen)
	afterData := make([]byte, afterDataLen)
	base64.StdEncoding.Encode(afterData, data)
	return afterData
}
