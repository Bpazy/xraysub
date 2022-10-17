package xray

import (
	"fmt"
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

func (s ShadowsocksServer) GetAddress() string {
	return s.Address
}

func (s ShadowsocksServer) GetPort() int {
	return s.Port
}

type TrojanServer struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func (s TrojanServer) GetAddress() string {
	return s.Address
}

func (s TrojanServer) GetPort() int {
	return s.Port
}

type User struct {
	Id       string `json:"id"`
	AlterId  int    `json:"alterId"`
	Email    string `json:"email"`
	Security string `json:"security"`
}

type Vnext struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Users   []User `json:"users"`
}

func (v Vnext) GetAddress() string {
	return v.Address
}

func (v Vnext) GetPort() int {
	return v.Port
}

type OutboundSettings struct {
	Servers []*ShadowsocksServer `json:"servers"`
	Vnext   []*Vnext             `json:"vnext"`
}

func (s OutboundSettings) GetAddressPort() AddressPort {
	if len(s.Servers) != 0 {
		return s.Servers[0]
	}
	return s.Vnext[0]
}

type AddressPort interface {
	GetAddress() string
	GetPort() int
}

type StreamSettings struct {
	Network string `json:"network"`
}

type Mux struct {
	Enabled     bool `json:"enabled"`
	Concurrency int  `json:"concurrency"`
}

type OutBound struct {
	Tag      string `json:"tag"`
	Protocol string `json:"protocol"`

	Latency        *time.Duration    `json:"-"` // server's latency
	Inbound        *Inbound          `json:"-"` // bound inbound for detecting latency
	Comment        string            `json:"-"`
	Settings       *OutboundSettings `json:"settings"`
	StreamSettings *StreamSettings   `json:"streamSettings"`
	Mux            *Mux              `json:"mux"`
}

func (o OutBound) PrettyComment() string {
	ap := o.Settings.GetAddressPort()
	addr := fmt.Sprintf("%s:%d", ap.GetAddress(), ap.GetPort())
	if o.Comment != "" {
		return fmt.Sprintf("%s(%s)", o.Comment, addr)
	}
	return addr
}

type Config struct {
	Policy    *Policy     `json:"policy"`
	Log       *Log        `json:"log"`
	Inbounds  []*Inbound  `json:"inbounds"`
	Outbounds []*OutBound `json:"outbounds"`
	Routing   *Routing    `json:"routing"`
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
