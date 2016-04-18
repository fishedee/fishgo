package sdk

import (
//"fmt"
//"testing"
)

type MailTestCase struct {
	host     string
	from     string
	password string
	to       []string
	subject  string
	body     string
}

func sendSingleCase(singleTestCase MailTestCase) error {
	mailSdk := SmtpSdk{
		Host: singleTestCase.host,
	}
	from := singleTestCase.from
	to := singleTestCase.to
	subject := singleTestCase.subject
	body := singleTestCase.body
	mail := SmtpSdkMail{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}
	mailAuth := SmtpSdkMailAuth{
		UserName: from,
		Password: singleTestCase.password,
	}
	return mailSdk.Send(mailAuth, mail)
}

/*
func TestSmtpSdkNormal(t *testing.T) {
	testCase := []MailTestCase{
		//不同文本
		{"smtp.qq.com:25", "jdlau@qq.com", "a19890305272", []string{"jdlau@qq.com"}, "这是测试邮件", "纯文本"},
		{"smtp.qq.com:25", "jdlau@qq.com", "a19890305272", []string{"jdlau@qq.com"}, "这是测试邮件2", "<a href=\"www.baidu.com\"><p>html结构化文本</p></a>"},
	}

	for _, singleTestCase := range testCase {
		err := sendSingleCase(singleTestCase)
		if err != nil {
			t.Error("error is not nil" + err.Error())
		}
	}
}


func TestSmtpSdkError(t *testing.T) {
	testCase := []MailTestCase{
		//something wrong
		{"smtp.qq.com:25", "jdlau@qq.com", "a19890305272", nil, "这是测试邮件3", "纯文本"},
		{"smtp.qq.com:25", "jdlau@qq.com", "wrongpass", []string{"jdlau@qq.com"}, "这是测试邮件4", "纯文本"},
	}

	for _, singleTestCase := range testCase {
		fmt.Printf("%#v\n", singleTestCase)
		err := sendSingleCase(singleTestCase)
		if err == nil {
			t.Error("error is not nil" + err.Error())
		}
	}
}
*/
