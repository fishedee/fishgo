// 获取html的快照
// 先模拟出一个浏览器，然后调用casperjs命令，生成页面的快照，保存到文件中
// casperjs是基于phantomjs建立的，所以必须先安装phantomjs
package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Browser struct {
	viewPort  *ViewPort
	userAgent string
	proxy     *Proxy
}

type ViewPort struct {
	width  int
	height int
}

type Proxy struct {
	host string
	port int
}

// 新建浏览器结构体
func NewBrowser(userAgent, host string, port, width, height int) *Browser {
	proxy := &Proxy{}
	if host != "" && port != 0 {
		proxy.host = host
		proxy.port = port
	}
	return &Browser{
		viewPort: &ViewPort{
			width:  width,
			height: height,
		},
		userAgent: userAgent,
		proxy:     proxy,
	}
}

// 获取指定页面的快照
// url: 文件路径
// path: 图片保存文件名
func (this *Browser) Snapshot(url, path string) error {
	// 拼接casperjs命令
	cmd := exec.Command(
		"/usr/local/casperjs/1.1.3/package/bin/casperjs",
		"--ssl-protocol=any",
		"--ignore-ssl-errors=true",
		"/var/www/fishgo/src/github.com/fishedee/util/scripts/snapshot.js",
		"--width="+strconv.Itoa(this.viewPort.width),
		"--height="+strconv.Itoa(this.viewPort.height),
		"--useragent="+this.userAgent,
		"--url="+url,
		"--path="+path,
		"--load-images=true",
	)

	// 设置运行环境
	env := this.getEnv()
	cmd.Env = env

	// 运行命令
	msg, err := cmd.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to run cmd: %v\n, msg is %v\n, env is %v\n", err, string(msg), env))
	}

	return nil
}

// 设置环境变量
func (this *Browser) getEnv() []string {
	// 待添加path部分
	addPATH := ":/usr/local/casperjs/1.1.3/package/bin:/home/jd/.npm/casperjs/1.1.3/package/bin:/usr/local/phantomjs/bin/"

	// 获取当前设置
	env := os.Environ()

	// 替换path部分设置
	newEnv := []string{}
	for _, single := range env {
		singleValue := single
		if strings.Contains(single, "PATH") {
			singleValue = single + addPATH
		}
		newEnv = append(newEnv, singleValue)
	}

	// 添加语言设置
	newEnv = append(
		newEnv,
		"LC_CTYPE=zh_CN.UTF-8",
	)

	return newEnv
}
