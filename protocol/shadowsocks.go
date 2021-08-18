package protocol

import (
	"encoding/base64"
	"regexp"
	"strconv"
)

var ssUriRe = regexp.MustCompile(`(?m)ss://(.+)@(.+):(\d+)`)
var ssCfgRe = regexp.MustCompile(`(?m)(.+):(.+)`)

type ShadowsocksConfig struct {
	Method   string
	Password string
	Hostname string
	Port     int
}

func ParseShadowsocksUri(uri string) (*ShadowsocksConfig, error) {
	ssUri := ssUriRe.FindStringSubmatch(uri)
	cfgBytes, err := base64.RawURLEncoding.DecodeString(ssUri[1])
	if err != nil {
		return nil, err
	}

	cfgs := ssCfgRe.FindStringSubmatch(string(cfgBytes))

	port := ssUri[3]
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	return &ShadowsocksConfig{
		Method:   cfgs[1],
		Password: cfgs[2],
		Hostname: ssUri[2],
		Port:     p,
	}, nil
}
