package cmd

import (
	"github.com/Bpazy/xraysub/gen"
	"github.com/spf13/cobra"
	"runtime"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate xray configuration file from url",
	Long:  "generate xray configuration file from url",
	Run:   gen.NewGenCmdRun(),
}

func init() {
	rootCmd.AddCommand(genCmd)

	const cUrl = "url"
	genCmd.Flags().StringVarP(&gen.Cfg.Url, cUrl, "u", "", "subscription address(URL)")
	genCmd.Flags().StringVarP(&gen.Cfg.OutputFile, "output-file", "o", "./xray-config.json", "output configuration to file")
	genCmd.Flags().BoolVarP(&gen.Cfg.Ping, "ping", "", true, "speed test (default true) to choose the fastest node")
	defaultXrayPath := "./xray"
	if runtime.GOOS == "windows" {
		defaultXrayPath = defaultXrayPath + ".exe"
	}
	genCmd.Flags().StringVarP(&gen.Cfg.XrayCorePath, "xray", "", defaultXrayPath, "speed test (default true) to choose the fastest node")
	cobra.CheckErr(genCmd.MarkFlagRequired(cUrl))
}
