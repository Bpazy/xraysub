package cmd

import (
	"github.com/Bpazy/xraysub/xray"
	"github.com/spf13/cobra"
)

var xrayCmd = &cobra.Command{
	Use:   "xray",
	Short: "Xray-core related commands",
	Long:  "Xray-core related commands",
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download xray-core",
	Long:  "download xray-core",
	RunE:  xray.NewXrayDownloadCmdRun(),
}

func init() {
	xrayCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(xrayCmd)

	downloadCmd.Flags().StringVarP(&xray.Cfg.GhProxy, "gh-proxy", "", "", "github proxy: https://github.com/hunshcn/gh-proxy")
}
