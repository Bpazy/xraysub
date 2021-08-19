package cmd

import (
	"github.com/Bpazy/xraysub/gen"
	"github.com/spf13/cobra"
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
	genCmd.Flags().StringVarP(&gen.Cfg.OutputFile, "output-file", "o", "", "output configuration to file")
	cobra.CheckErr(genCmd.MarkFlagRequired(cUrl))
}
