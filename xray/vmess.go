package xray

import (
	"time"
)

type System struct {
	StatsOutboundUplink   bool `json:"statsOutboundUplink"`
	StatsOutboundDownlink bool `json:"statsOutboundDownlink"`
}

type Policy struct {
	System System `json:"system"`
}

type Log struct {
	Access   string `json:"access"`
	Error    string `json:"error"`
	Loglevel string `json:"loglevel"`
}

type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

type InboundSettings struct {
	Auth             string `json:"auth,omitempty"`
	Udp              bool   `json:"udp"`
	AllowTransparent bool   `json:"allowTransparent"`
}

type Inbound struct {
	Tag      string           `json:"tag"`
	Port     int              `json:"port"`
	Listen   string           `json:"listen"`
	Protocol string           `json:"protocol"`
	Sniffing *Sniffing        `json:"sniffing"`
	Settings *InboundSettings `json:"settings"`
}

type ShadowsocksServer struct {
	Address  string `json:"address"`
	Method   string `json:"method"`
	Ota      bool   `json:"ota"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Level    int    `json:"level"`
}

type OutboundSettings struct {
	Servers []*ShadowsocksServer `json:"servers"`
}

type StreamSettings struct {
	Network string `json:"network"`
}

type Mux struct {
	Enabled     bool `json:"enabled"`
	Concurrency int  `json:"concurrency"`
}

type BaseOutbound struct {
	Tag      string `json:"tag"`
	Protocol string `json:"protocol"`

	Latency *time.Duration `json:"-"` // server's latency
	Inbound *Inbound       `json:"-"` // bound inbound for detecting latency
}

type ShadowsocksOutbound struct {
	BaseOutbound
	Settings       *OutboundSettings `json:"settings"`
	StreamSettings *StreamSettings   `json:"streamSettings"`
	Mux            *Mux              `json:"mux"`
}

type Config struct {
	Policy    *Policy                `json:"policy"`
	Log       *Log                   `json:"log"`
	Inbounds  []*Inbound             `json:"inbounds"`
	Outbounds []*ShadowsocksOutbound `json:"outbounds"`
	Routing   *Routing               `json:"routing"`
}

type Rule struct {
	Type        string   `json:"type"`
	InboundTag  []string `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag"`
	Port        string   `json:"port,omitempty"`
}

type Routing struct {
	DomainStrategy string  `json:"domainStrategy"`
	DomainMatcher  string  `json:"domainMatcher"`
	Rules          []*Rule `json:"rules"`
}
