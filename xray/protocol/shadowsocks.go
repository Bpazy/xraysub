package protocol

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

type ShadowsocksConfig struct {
	Method   string
	Password string
	Hostname string
	Port     int
	Comment  string
}

var ssUriRe = regexp.MustCompile(`(?m)ss://(.+)@(.+):(\d+)(#(.+))?`)
var ssCfgRe = regexp.MustCompile(`(?m)(.+):(.+)`)

func ParseShadowsocksUri(uri string) (*ShadowsocksConfig, error) {
	ssUri := ssUriRe.FindStringSubmatch(uri)
	cfgBytes, err := base64.RawURLEncoding.DecodeString(ssUri[1])
	if err != nil {
		return nil, err
	}

	cfgs := ssCfgRe.FindStringSubmatch(string(cfgBytes))

	// comment
	comment, err := getComment(ssUri[5])
	if err != nil {
		return nil, err
	}

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
		Comment:  comment,
	}, nil
}

func getComment(uri string) (string, error) {
	if comment, err := url.QueryUnescape(uri); err != nil {
		return "", fmt.Errorf("failed to unescape comment: %w", err)
	} else {
		return comment, nil
	}
}
