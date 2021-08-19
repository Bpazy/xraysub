package cmd

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Bpazy/xraysub/protocol"
	"github.com/Bpazy/xraysub/vmess"
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

		xraycfg := &vmess.XrayConfig{
			Policy: &vmess.Policy{
				System: vmess.System{
					StatsOutboundUplink:   true,
					StatsOutboundDownlink: true,
				},
			},
			Log: &vmess.Log{
				Access:   "",
				Error:    "",
				Loglevel: "warning",
			},
			Inbounds:  getInbounds(),
			Outbounds: getOutBounds(links),
			Routing: &vmess.Routing{
				DomainStrategy: "IPIfNonMatch",
				DomainMatcher:  "linear",
				Rules: []*vmess.Rule{
					{
						Type:        "field",
						OutboundTag: "proxy",
						Port:        "0-65535",
					},
				},
			},
		}
		cfg, err := json.Marshal(xraycfg)
		cobra.CheckErr(err)
		err = ioutil.WriteFile(genCmdCfg.GetOutputFile(), cfg, 0644)
		cobra.CheckErr(err)
	},
}

func getOutBounds(links []*Link) []*vmess.ShadowsocksOutbound {
	var outbounds []*vmess.ShadowsocksOutbound
	for _, link := range links {
		outbounds = append(outbounds, &vmess.ShadowsocksOutbound{
			BaseOutbound: vmess.BaseOutbound{
				Tag:      "proxy", // 应该测速后选择最合适的设置 tag 为 proxy
				Protocol: "shadowsocks",
			},
			Settings: &vmess.OutboundSettings{
				Servers: []*vmess.ShadowsocksServer{
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
			StreamSettings: &vmess.StreamSettings{
				Network: "tcp",
			},
			Mux: &vmess.Mux{
				Enabled:     false,
				Concurrency: -1,
			},
		})
	}
	return outbounds
}

func getInbounds() []*vmess.Inbound {
	return []*vmess.Inbound{
		{
			Tag:      "socks",
			Port:     10808,
			Listen:   "0.0.0.0",
			Protocol: "socks",
			Sniffing: &vmess.Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: &vmess.InboundSettings{
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
			Sniffing: &vmess.Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: &vmess.InboundSettings{
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
