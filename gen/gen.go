package gen

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Bpazy/xraysub/xray"
	"github.com/Bpazy/xraysub/xray/protocol"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// CmdConfig config for command: xraysub gen
type CmdConfig struct {
	Url           string // subscription link
	OutputFile    string // xray-core's configuration path
	Ping          bool   // speed test to choose the fastest node
	XrayCorePath  string // xray-core path, for some case such as: speed test
	XraySocksPort int    // xray-core listen socks port
	XrayHttpPort  int    // xray-core listen http port
}

var Cfg = &CmdConfig{}

func NewGenCmdRun() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		c := resty.New()
		res, err := c.R().Get(Cfg.Url)
		cobra.CheckErr(err)
		dst, err := base64.StdEncoding.DecodeString(res.String())
		cobra.CheckErr(err)
		uris := strings.Split(strings.TrimSpace(string(dst)), "\n")
		links := parseLinks(uris)

		xrayCfg := getXrayConfig(links)
		if Cfg.Ping {
			_, err := ping(xrayCfg)
			cobra.CheckErr(err)
		}

		writeFile(xrayCfg, Cfg.OutputFile)
	}
}

// speed test, return the fastest node
func ping(cfg *xray.XrayConfig) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "xray.config.json")
	if err != nil {
		return "", err
	}
	// random port for speed test
	rand.Seed(time.Now().UnixNano())
	httpPort := rand.Intn(1000) + 40000
	socksPort := rand.Intn(1000) + 40000
	for _, inbound := range cfg.Inbounds {
		if inbound.Protocol == "socks" {
			inbound.Port = socksPort
		}
		if inbound.Protocol == "http" {
			inbound.Port = httpPort
		}
	}
	for _, outbound := range cfg.Outbounds {
		outbound.Tag = "proxy"
		j, err := json.Marshal(cfg)
		if err != nil {
			return "", err
		}
		_, err = f.Write(j)
		if err != nil {
			return "", err
		}

		cmd := exec.Command(Cfg.XrayCorePath, "-c", f.Name(), "-format=json")
		if err != nil {
			return "", fmt.Errorf("init xray-core error: %w", err)
		}
		stdoutBuf := new(bytes.Buffer)
		stderrBuf := new(bytes.Buffer)
		cmd.Stdout = stdoutBuf
		cmd.Stderr = stderrBuf
		err = cmd.Start()
		if err != nil {
			return "", fmt.Errorf("exec xray-core error: %w", err)
		}
		log.Infof("xray-core pid: %d", cmd.Process.Pid)

		client := resty.New()
		proxy := "http://127.0.0.1:" + strconv.Itoa(httpPort)
		client.SetProxy(proxy)
		client.SetTimeout(5 * time.Second)
		start := time.Now()
		_, err = client.R().Get("https://www.google.com/generate_204")
		if err != nil {
			log.Infof("start xray-core error: (stdout: %s), (stderr: %s)", stdoutBuf.String(), stderrBuf.String())
			log.Errorf("request failed by proxy %s: %+v", proxy, err)
			err = killProcess(cmd)
			if err != nil {
				return "", err
			}
			continue
		}
		log.Infof("%s spent %dms", outbound.Settings.Servers[0].Address, time.Since(start).Milliseconds())
		err = killProcess(cmd)
		if err != nil {
			return "", err
		}
	}

	// rollback port
	for _, inbound := range cfg.Inbounds {
		if inbound.Protocol == "socks" {
			inbound.Port = Cfg.XraySocksPort
		}
		if inbound.Protocol == "http" {
			inbound.Port = Cfg.XrayHttpPort
		}
	}
	return "", nil
}

func killProcess(cmd *exec.Cmd) error {
	err := cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("kill xray-core error: %w", err)
	}
	return nil
}

// write xray-core's configuration
func writeFile(cfg *xray.XrayConfig, path string) {
	j, err := json.Marshal(cfg)
	cobra.CheckErr(err)
	err = ioutil.WriteFile(path, j, 0644)
	cobra.CheckErr(err)
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

type Link struct {
	SsCfg *protocol.ShadowsocksConfig
}

// build xray-core config
func getXrayConfig(links []*Link) *xray.XrayConfig {
	return &xray.XrayConfig{
		Policy: &xray.Policy{
			System: xray.System{
				StatsOutboundUplink:   true,
				StatsOutboundDownlink: true,
			},
		},
		Log: &xray.Log{
			Access:   "",
			Error:    "",
			Loglevel: "warning",
		},
		Inbounds:  getInbounds(),
		Outbounds: getOutBounds(links),
		Routing: &xray.Routing{
			DomainStrategy: "IPIfNonMatch",
			DomainMatcher:  "linear",
			Rules: []*xray.Rule{
				{
					Type:        "field",
					OutboundTag: "proxy",
					Port:        "0-65535",
				},
			},
		},
	}
}

func getOutBounds(links []*Link) []*xray.ShadowsocksOutbound {
	var outbounds []*xray.ShadowsocksOutbound
	for _, link := range links {
		outbounds = append(outbounds, &xray.ShadowsocksOutbound{
			BaseOutbound: xray.BaseOutbound{
				Tag:      "", // 应该测速后选择最合适的设置 tag 为 proxy
				Protocol: "shadowsocks",
			},
			Settings: &xray.OutboundSettings{
				Servers: []*xray.ShadowsocksServer{
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
			StreamSettings: &xray.StreamSettings{
				Network: "tcp",
			},
			Mux: &xray.Mux{
				Enabled:     false,
				Concurrency: -1,
			},
		})
	}
	return outbounds
}

func getInbounds() []*xray.Inbound {
	return []*xray.Inbound{
		{
			Tag:      "socks",
			Port:     Cfg.XraySocksPort,
			Listen:   "0.0.0.0",
			Protocol: "socks",
			Sniffing: &xray.Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: &xray.InboundSettings{
				Auth:             "noauth",
				Udp:              true,
				AllowTransparent: false,
			},
		},
		{
			Tag:      "http",
			Port:     Cfg.XrayHttpPort,
			Listen:   "0.0.0.0",
			Protocol: "http",
			Sniffing: &xray.Sniffing{
				Enabled:      true,
				DestOverride: []string{"http", "tls"},
			},
			Settings: &xray.InboundSettings{
				Udp:              false,
				AllowTransparent: false,
			},
		},
	}
}
