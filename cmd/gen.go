package cmd

import (
	"github.com/Bpazy/xraysub/gen"
	"github.com/Bpazy/xraysub/util"
	"github.com/spf13/cobra"
	"runtime"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate xray configuration file from subscription url",
	Long:  "generate xray configuration file from subscription url",
	Run:   gen.NewGenCmdRun(),
}

func init() {
	rootCmd.AddCommand(genCmd)

	const cUrl = "url"
	genCmd.Flags().StringVarP(&gen.Cfg.Url, cUrl, "u", "", "subscription address(URL)")
	util.CheckErr(genCmd.MarkFlagRequired(cUrl))
	genCmd.Flags().StringVarP(&gen.Cfg.OutputFile, "output-file", "o", "./xray-config.json", "output configuration to file")
	genCmd.Flags().BoolVarP(&gen.Cfg.DetectLatency, "detect-latency", "", true, "detect server's latency to choose the fastest node")
	genCmd.Flags().StringVarP(&gen.Cfg.XrayCorePath, "xray", "", getDefaultXrayPath(), "xray-core path for detecting server's latency")
	genCmd.Flags().IntVarP(&gen.Cfg.XraySocksPort, "xray-socks-port", "", 1080, "xray-core listen socks port")
	genCmd.Flags().IntVarP(&gen.Cfg.XrayHttpPort, "xray-http-port", "", 1081, "xray-core listen http port")
	genCmd.Flags().IntVarP(&gen.Cfg.DetectThreadNumber, "detect-thread-number", "", 5, "detect server's latency threads number")
}

func getDefaultXrayPath() string {
	defaultXrayPath := "./xray"
	if runtime.GOOS == "windows" {
		defaultXrayPath = defaultXrayPath + ".exe"
	}
	return defaultXrayPath
}
