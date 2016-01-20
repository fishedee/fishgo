package util

import (
	"errors"
	"net"
)

type Ifconfig struct {
}

func NewIfconfig() *Ifconfig {
	return &Ifconfig{}
}

func (this *Ifconfig) GetIP(name string) (*net.IPNet, error) {
	in, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	addrs, err := in.Addrs()
	if err != nil {
		return nil, err
	}
	for _, singleAddr := range addrs {
		addrIpNet, ok := singleAddr.(*net.IPNet)
		if ok {
			return addrIpNet, nil
		}
	}
	return nil, errors.New(name + " has nothing ip")
}
