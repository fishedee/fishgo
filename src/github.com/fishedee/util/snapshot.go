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
		"casperjs",
		"--ssl-protocol=any",
		"--ignore-ssl-errors=true",
		"/home/jd/Project/fishgo/src/github.com/fishedee/util/scripts/snapshot.js",
		"--width="+strconv.Itoa(this.viewPort.width),
		"--height="+strconv.Itoa(this.viewPort.height),
		"--useragent="+this.userAgent,
		"--url="+url,
		"--path="+path,
		"--load-images=true",
	)

	// 设置运行环境
	env := os.Environ()
	env = append(
		env,
		"LC_CTYPE=zh_CN.UTF-8",
		"PATH=/usr/local/node/bin:/usr/local/phantomjs/bin",
	)
	cmd.Env = env

	// 运行命令
	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to run cmd: %v\n", err))
	}

	return nil
}
