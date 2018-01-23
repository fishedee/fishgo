package sdk

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/fishedee/encoding"
	"github.com/fishedee/language"
	"gopkg.in/gomail.v2"
	"io"
	"strconv"
)

type SmtpSdk struct {
	Host      string
	SSL       bool
	TLSConfig *tls.Config
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

func (this *SmtpSdk) setAttach(m *gomail.Message, singleAttach SmtpSdkMailAttach) {
	base64name, _ := encoding.EncodeBase64([]byte(singleAttach.Name))
	base64nameString := "=?utf-8?B?" + string(base64name) + "?="
	m.Attach(base64nameString, gomail.SetCopyFunc(func(dest io.Writer) error {
		src := bytes.NewReader(singleAttach.Data)
		_, err := io.Copy(dest, src)
		return err
	}))
}
func (this *SmtpSdk) Send(mailAuth SmtpSdkMailAuth, mail SmtpSdkMail) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To...)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Body)
	for _, singleAttach := range mail.Attach {
		this.setAttach(m, singleAttach)
	}

	addressInfo := language.Explode(this.Host, ":")
	if len(addressInfo) != 2 {
		return errors.New("invalid host")
	}
	host := addressInfo[0]
	port, err := strconv.Atoi(addressInfo[1])
	if err != nil {
		return err
	}
	d := &gomail.Dialer{
		Host:     host,
		Port:     port,
		Username: mailAuth.UserName,
		Password: mailAuth.Password,
		SSL:      this.SSL,
	}
	d.TLSConfig = this.TLSConfig
	return d.DialAndSend(m)
}
