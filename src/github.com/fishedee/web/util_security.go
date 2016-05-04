package util

import (
	"errors"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"runtime"
)

type SecurityManagerConfig struct {
	IpWhite []string
}

type SecurityManager struct {
}

func NewSecurityManager(config SecurityManagerConfig) (*SecurityManager, error) {
	var netConfig string
	if runtime.GOOS == "darwin" {
		netConfig = "en0"
	} else {
		netConfig = "eth0"
	}
	ip, err := NewIfconfig().GetIP(netConfig)
	if err != nil {
		return nil, err
	}

	ipStr := ip.IP.String()
	if len(config.IpWhite) != 0 && ArrayIn(config.IpWhite, ipStr) == -1 {
		return nil, errors.New("当前IP: " + ipStr + "不在IP白名单中: " + Implode(config.IpWhite, ","))
	}

	return &SecurityManager{}, nil
}

func NewSecurityManagerFromConfig(configName string) (*SecurityManager, error) {
	ipwhite := globalBasic.Config.String(configName + "ipwhite")
	ipwhiteList := Explode(ipwhite, ",")
	return NewSecurityManager(SecurityManagerConfig{IpWhite: ipwhiteList})
}
