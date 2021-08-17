package cmd

import (
	"encoding/base64"
	"github.com/Bpazy/xraysub/protocol"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

type genCmdConfig struct {
	url string
}

var genCmdCfg = &genCmdConfig{}

type Link struct {
	ssCfg *protocol.ShadowsocksConfig
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

		for _, cfg := range links {
			log.Printf("Shadowsocks cfg: %+v", cfg)
		}
	},
}

func parseLinks(uris []string) []*Link {
	links := make([]*Link, len(uris))
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
				ssCfg: cfg,
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
	cobra.CheckErr(genCmd.MarkFlagRequired(cUrl))
}
