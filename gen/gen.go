package gen

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Bpazy/xraysub/util"
	"github.com/Bpazy/xraysub/xray"
	"github.com/Bpazy/xraysub/xray/protocol"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
		fmt.Printf("Requesting subscriptions from %s\n", Cfg.Url)
		c.SetTimeout(5 * time.Second)
		res, err := c.R().Get(Cfg.Url)
		util.CheckErr(err)
		dst, err := base64.StdEncoding.DecodeString(res.String())
		util.CheckErr(err)
		uris := strings.Split(strings.TrimSpace(string(dst)), "\n")
		links := parseLinks(uris)

		xrayCfg := getXrayConfig(links)
		if Cfg.Ping {
			err := ping(xrayCfg)
			util.CheckErr(err)
		}

		writeFile(xrayCfg, Cfg.OutputFile)
	}
}

var xrayCoreProcess *os.Process

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-c
		fmt.Println("Got signal:", s)
		killXrayCoreProcess()
	}()
}

func killXrayCoreProcess() {
	if xrayCoreProcess != nil {
		_ = xrayCoreProcess.Kill()
	}
}

// speed test, return the fastest node
func ping(xCfg *xray.Config) error {
	fmt.Println("Ping test processing")
	if len(xCfg.Outbounds) == 0 {
		return errors.New("outbounds empty")
	}
	// 根据 outbounds 生成 inbounds 和 routing rules
	var inbounds []*xray.Inbound
	var routingRules []*xray.Rule
	inboundPorts := randomInboundPorts(xCfg.Outbounds)
	for i, outbound := range xCfg.Outbounds {
		inbound := getInboundFromOutbound(i, inboundPorts[i])
		outbound.Inbound = inbound
		inbounds = append(inbounds, inbound)
		routingRules = append(routingRules, getRoutingRules(inbound, outbound))
	}
	oldInbounds := xCfg.Inbounds
	xCfg.Inbounds = inbounds
	xCfg.Routing.Rules = routingRules
	defer func() {
		xCfg.Inbounds = oldInbounds
	}()

	f, err := writeTempConfig(xCfg)
	if err != nil {
		return err
	}

	cmd := exec.Command(Cfg.XrayCorePath, "-c", f.Name(), "-format=json")
	if err != nil {
		return fmt.Errorf("init xray-core command error: %w", err)
	}

	xlf, err := appendXrayCoreLogFile()
	if err != nil {
		return fmt.Errorf("create xray-core.log error: %w", err)
	}
	defer xlf.Close()
	cmd.Stdout = xlf
	cmd.Stderr = xlf
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("exec xray-core error: %w", err)
	}
	log.Infof("xray-core PID: %d", cmd.Process.Pid)
	xrayCoreProcess = cmd.Process

	wg := new(sync.WaitGroup)
	outboundChan := make(chan *xray.ShadowsocksOutbound)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go pingWorker(outboundChan, wg)
	}
	for _, outbound := range xCfg.Outbounds {
		outboundChan <- outbound
	}
	close(outboundChan)
	wg.Wait()

	killXrayCoreProcess()

	// filter fasted outbound
	var fastedOutbound *xray.ShadowsocksOutbound
	for _, outbound := range xCfg.Outbounds {
		if outbound.PingDelay == nil {
			continue
		}
		if fastedOutbound == nil {
			fastedOutbound = outbound
		} else if fastedOutbound.PingDelay.Milliseconds() > outbound.PingDelay.Milliseconds() {
			fastedOutbound = outbound
		}
	}
	if fastedOutbound == nil {
		fmt.Println("All nodes ping test failed")
	} else {
		s := fastedOutbound.Settings.Servers[0]
		fmt.Printf("Got fastest node: %s:%d\n", s.Address, s.Port)
	}
	xCfg.Routing.Rules = []*xray.Rule{
		{
			Type:        "field",
			OutboundTag: fastedOutbound.Tag,
			Port:        "0-65535",
		},
	}

	return nil
}

func appendXrayCoreLogFile() (*os.File, error) {
	f, err := os.OpenFile("xray-core.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return f, err
}

func pingWorker(oc chan *xray.ShadowsocksOutbound, wg *sync.WaitGroup) {
	defer wg.Done()

	for outbound := range oc {
		client := resty.New()
		proxy := "http://127.0.0.1:" + strconv.Itoa(outbound.Inbound.Port)
		client.SetProxy(proxy)
		client.SetTimeout(5 * time.Second)
		start := time.Now()
		_, err := client.R().Get("https://www.google.com/generate_204")
		if err != nil {
			log.Errorf("request failed by proxy %s: %+v", proxy, err)
		} else {
			since := time.Since(start)
			outbound.PingDelay = &since
			s := outbound.Settings.Servers[0]
			log.Infof("%s:%d cost %dms", s.Address, s.Port, outbound.PingDelay.Milliseconds())
		}
	}
}

func randomInboundPorts(outbounds []*xray.ShadowsocksOutbound) []int {
	var ports []int
	offset := 0
	for range outbounds {
		for {
			p := 40000 + offset
			offset++
			listenable := portListenable(p)
			if listenable {
				ports = append(ports, p)
				break
			}
		}
	}
	return ports
}

func portListenable(p int) bool {
	listen, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(p))
	if err == nil {
		// 端口可用
		_ = listen.Close()
		return true
	}
	return false
}

func writeTempConfig(xCfg *xray.Config) (*os.File, error) {
	j, err := json.Marshal(xCfg)
	if err != nil {
		return nil, err
	}
	f, err := ioutil.TempFile(os.TempDir(), "xray.config.json")
	if err != nil {
		return nil, fmt.Errorf("create temp file 'xray.config.json' error: %w", err)
	}
	_, err = f.Write(j)
	if err != nil {
		return nil, fmt.Errorf("write to temp file 'xray.config.json' error: %w", err)
	}
	return f, nil
}

func getRoutingRules(inbound *xray.Inbound, outbound *xray.ShadowsocksOutbound) *xray.Rule {
	return &xray.Rule{
		Type:        "field",
		InboundTag:  []string{inbound.Tag},
		OutboundTag: outbound.Tag,
	}
}

func getInboundFromOutbound(i int, port int) *xray.Inbound {
	return &xray.Inbound{
		Tag:      "inbound" + strconv.Itoa(i),
		Port:     port,
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
	}
}

// write xray-core's configuration
func writeFile(cfg *xray.Config, path string) {
	j, err := json.Marshal(cfg)
	util.CheckErr(err)
	err = ioutil.WriteFile(path, j, 0644)
	fmt.Printf("The xray-core's configuration file is saved %s\n", path)
	util.CheckErr(err)
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
func getXrayConfig(links []*Link) *xray.Config {
	return &xray.Config{
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
		},
	}
}

func getOutBounds(links []*Link) []*xray.ShadowsocksOutbound {
	var outbounds []*xray.ShadowsocksOutbound
	for i, link := range links {
		outbounds = append(outbounds, &xray.ShadowsocksOutbound{
			BaseOutbound: xray.BaseOutbound{
				Tag:      "outbound" + strconv.Itoa(i), // 应该测速后选择最合适的设置 tag 为 proxy
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
