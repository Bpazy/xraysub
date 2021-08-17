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
)

func GetProtocol(uri string) (Type, error) {
	split := strings.Split(uri, "://")
	if len(split) == 0 {
		return noneProtocolType()
	}
	switch split[0] {
	case "ss":
		return Shadowsocks, nil
	case "vmess":
		return Vmess, nil
	}
	return noneProtocolType()
}

var ErrWrongProtocol = errors.New("wrong protocol")

func noneProtocolType() (Type, error) {
	return None, ErrWrongProtocol
}
