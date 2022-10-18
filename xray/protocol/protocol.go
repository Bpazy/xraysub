package protocol

import (
	"errors"
	"strings"
)

type Type string

const (
	None        = ""
	Vmess       = "vmess"
	Shadowsocks = "ss"
	Trojan      = "trojan"
)

func GetProtocol(uri string) (Type, error) {
	split := strings.Split(uri, "://")
	if len(split) == 0 {
		return noneProtocolType()
	}
	switch split[0] {
	case Shadowsocks:
		return Shadowsocks, nil
	case Vmess:
		return Vmess, nil
	case Trojan:
		return Trojan, nil
	}
	return noneProtocolType()
}

var ErrWrongProtocol = errors.New("wrong protocol")

func noneProtocolType() (Type, error) {
	return None, ErrWrongProtocol
}
