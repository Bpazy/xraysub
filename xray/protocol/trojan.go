package protocol

import (
	"regexp"
	"strconv"
)

type TrojanConfig struct {
	Password string
	Host     string
	Port     int
}

var trojanUriRe = regexp.MustCompile(`(?m)trojan://(.+)@(.+):(\d+)\??`)

func ParseTrojanUri(uri string) (*TrojanConfig, error) {
	trojanUri := trojanUriRe.FindStringSubmatch(uri)
	p, err := strconv.Atoi(trojanUri[3])
	if err != nil {
		return nil, err
	}
	return &TrojanConfig{
		Password: trojanUri[1],
		Host:     trojanUri[2],
		Port:     p,
	}, nil
}
