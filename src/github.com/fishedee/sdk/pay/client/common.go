package client

import (
	"fmt"
	"sort"
	"strings"
	"crypto/md5"
	"errors"
)

func WechatGenSign(key string,m map[string]string) (string,error) {
	var signData []string
	for k, v := range m {
		if v != "" && k != "sign" && k != "key"{
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}
	fmt.Printf("%+v",signData)

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + key
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", errors.New("WechatGenSign md5.Write: "+err.Error())
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", errors.New("WechatGenSign md5.Sum: "+err.Error())
	}
	return strings.ToUpper(fmt.Sprintf("%x", signByte)),nil
}

func TruncatedText(data string,length int) string{
	data = FilterTheSpecialSymbol(data)
	if len([]rune(data)) > length {
		return string([]rune(data)[:length-1])
	}
	return data
}

//过滤特殊符号
func FilterTheSpecialSymbol(data string) string {
	// 定义转换规则
	specialSymbol := func(r rune) rune {
		if r == '`' || r == '[' || r == '~' || r == '!' || r == '@' || r == '#' || r == '$' ||
			r == '^' || r == '&' || r == '*' || r == '~' || r == '(' || r == ')' || r == '=' ||
			r == '~' || r == '|' || r == '{' || r == '}' || r == '~' || r == ':' || r == ';' ||
			r == '\'' || r == ',' || r == '\\' || r == '[' || r == ']' || r == '.' || r == '<' ||
			r == '>' || r == '/' || r == '?' || r == '~' || r == '！' || r == '@' || r == '#' ||
			r == '￥' || r == '…' || r == '&' || r == '*' || r == '（' || r == '）' || r == '—' ||
			r == '|' || r == '{' || r == '}' || r == '【' || r == '】' || r == '‘' || r == '；' ||
			r == '：' || r == '”' || r == '“' || r == '\'' || r == '"' || r == '。' || r == '，' ||
			r == '、' || r == '？' || r == '%' || r == '+' || r == '_' || r == ']' || r == '"' || r == '&' {
			return ' '
		}
		return r
	}
	data = strings.Map(specialSymbol, data)
	return strings.Replace(data, "\n", " ", -1)
}
