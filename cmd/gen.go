package cmd

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Bpazy/xraysub/protocol"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

type genCmdConfig struct {
	url        string
	outputFile string
}

var genCmdCfg = &genCmdConfig{}

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
	Tag      string          `json:"tag"`
	Port     int             `json:"port"`
	Listen   string          `json:"listen"`
	Protocol string          `json:"protocol"`
	Sniffing Sniffing        `json:"sniffing"`
	Settings InboundSettings `json:"settings"`
}

type BaseOutbound struct {
	Tag      string `json:"tag"`
	Protocol string `json:"protocol"`
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
	Servers []ShadowsocksServer `json:"servers"`
}

type StreamSettings struct {
	Network string `json:"network"`
}

type Mux struct {
	Enabled     bool `json:"enabled"`
	Concurrency int  `json:"concurrency"`
}

type ShadowsocksOutbound struct {
	BaseOutbound
	Settings       *OutboundSettings `json:"settings"`
	StreamSettings *StreamSettings   `json:"streamSettings"`
	Mux            *Mux              `json:"mux"`
}

type XrayConfig struct {
	Policy    *Policy                `json:"policy"`
	Log       *Log                   `json:"log"`
	Inbounds  []*Inbound             `json:"inbounds"`
	Outbounds []*ShadowsocksOutbound `json:"outbounds"`
}

func (c genCmdConfig) GetOutputFile() string {
	if c.outputFile != "" {
		return c.outputFile
	}
	// default output file
	return "./xray-config.json"
}

type Link struct {
	SsCfg *protocol.ShadowsocksConfig
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate xray configuration file from url",
	Long:  "generate xray configuration file from url",
	Run: func(cmd *cobra.Command, args []string) {
		c := resty.New()
		res, err := c.R().Get(genCmdCfg.url)
		cobra.CheckErr(err)

		dst, err := base64.StdEncoding.DecodeString(res.String())
		cobra.CheckErr(err)

		uris := strings.Split(strings.TrimSpace(string(dst)), "\n")
		links := parseLinks(uris)

		// pretty print
		for _, cfg := range links {
			j, _ := json.Marshal(cfg)
			log.Printf("Shadowsocks cfg: %s", string(j))
		}

		xraycfg := &XrayConfig{
			Policy: &Policy{
				System: System{
					StatsOutboundUplink:   true,
					StatsOutboundDownlink: true,
				},
			},
			Log: &Log{
				Access:   "",
				Error:    "",
				Loglevel: "warning",
			},
			Inbounds:  getInbounds(),
			Outbounds: getOutBounds(links),
		}
		cfg, err := json.Marshal(xraycfg)
		cobra.CheckErr(err)
		err = ioutil.WriteFile(genCmdCfg.GetOutputFile(), cfg, 0644)
		cobra.CheckErr(err)
	},
}

func getOutBounds(links []*Link) []*ShadowsocksOutbound {
	var outbounds []*ShadowsocksOutbound
	for _, link := range links {
		outbounds = append(outbounds, &ShadowsocksOutbound{
			BaseOutbound: BaseOutbound{
				Tag:      "proxy",
				Protocol: "shadowsocks",
			},
			Settings: &OutboundSettings{
				Servers: []ShadowsocksServer{
					{
						Address:  link.SsCfg.Hostname,
						Method:   link.SsCfg.Method,
						Ota:      false,
						Password: link.SsCfg.Password,
						Port:     link.SsCfg.Port,
						Level:    1,
					},
				},
			},
			StreamSettings: &StreamSettings{
				Network: "tcp",
			},
			Mux: &Mux{
				Enabled:     false,
				Concurrency: -1,
			},
		})
	}
	return outbounds
}

func getInbounds() []*Inbound {
	return []*Inbound{
		{
			Tag:      "socks",
			Port:     10808,
			Listen:   "0.0.0.0",
			Protocol: "socks",
			Sniffing: Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: InboundSettings{
				Auth:             "noauth",
				Udp:              true,
				AllowTransparent: false,
			},
		},
		{
			Tag:      "http",
			Port:     10809,
			Listen:   "0.0.0.0",
			Protocol: "http",
			Sniffing: Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: InboundSettings{
				Udp:              false,
				AllowTransparent: false,
			},
		},
	}
}

func parseLinks(uris []string) []*Link {
	var links []*Link
	for _, uri := range uris {
		p, err := protocol.GetProtocol(uri)
		if err != nil {
			log.Warn("unrecognized protocol: " + uri)
			continue
		}

		switch p {
		case protocol.Shadowsocks:
			cfg, err := protocol.ParseShadowsocksUri(uri)
			if err != nil {
				log.Warn("illegal shadowsocks uri schema: " + uri)
				continue
			}
			links = append(links, &Link{
				SsCfg: cfg,
			})
		case protocol.Vmess:

		}
	}
	return links
}

func init() {
	rootCmd.AddCommand(genCmd)

	const cUrl = "url"
	genCmd.Flags().StringVarP(&genCmdCfg.url, cUrl, "u", "", "subscription address(URL)")
	genCmd.Flags().StringVarP(&genCmdCfg.outputFile, "output-file", "o", "", "output configuration to file")
	cobra.CheckErr(genCmd.MarkFlagRequired(cUrl))
}
