package sdk

import (
	//"fmt"
	"crypto/tls"
	"testing"
)

func GGTestSmtpSdkNormal(t *testing.T) {
	smtp := &SmtpSdk{
		Host: "smtp.163.com:25",
	}
	err := smtp.Send(SmtpSdkMailAuth{
		UserName: "15018749403@163.com",
		Password: "password",
	}, SmtpSdkMail{
		From:    "15018749403@163.com",
		To:      []string{"306766045@qq.com"},
		Subject: "测试标题_normal",
		Body:    "测试body_noraml<a href=\"http://www.baidu.com\"><p>html结构化文本</p></a>",
	})
	if err != nil {
		panic(err)
	}
}

func GGTestSmtpSdkSSL(t *testing.T) {
	smtp := &SmtpSdk{
		Host:      "smtp.163.com:465",
		SSL:       true,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	err := smtp.Send(SmtpSdkMailAuth{
		UserName: "15018749403@163.com",
		Password: "password",
	}, SmtpSdkMail{
		From:    "15018749403@163.com",
		To:      []string{"306766045@qq.com"},
		Subject: "测试标题_ssl",
		Body:    "测试body_ssl",
	})
	if err != nil {
		panic(err)
	}
}

func GGTestSmtpSdkAttach(t *testing.T) {
	smtp := &SmtpSdk{
		Host:      "smtp.163.com:465",
		SSL:       true,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	err := smtp.Send(SmtpSdkMailAuth{
		UserName: "15018749403@163.com",
		Password: "password",
	}, SmtpSdkMail{
		From:    "15018749403@163.com",
		To:      []string{"306766045@qq.com"},
		Subject: "测试标题_attach",
		Body:    "测试body_attach",
		Attach: []SmtpSdkMailAttach{
			SmtpSdkMailAttach{
				Name: "文件1.txt",
				Data: []byte("文件内容1"),
			},
			SmtpSdkMailAttach{
				Name: "文件2.txt",
				Data: []byte("文件内容2"),
			},
		},
	})
	if err != nil {
		panic(err)
	}
}
