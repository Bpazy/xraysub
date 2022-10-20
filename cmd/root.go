package cmd

import (
	"github.com/Bpazy/xraysub/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	// buildVer represents 'xraysub' build version
	buildVer string
	verbose  bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "xraysub",
		Short: "",
		Long: `不 畏 浮 云 遮 望 眼 · 金 睛 如 炬 耀 苍 穹
K E E P   R I D I N G   /   N E V E R   L O O K   B A C K`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func initConfig() {
	log.SetFormatter(&log.TextFormatter{
		DisableQuote: true,
	})

	file, err := getLogFile()
	util.CheckErr(err)
	log.SetOutput(file)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func getLogFile() (*os.File, error) {
	f, err := os.OpenFile("xraysub.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return f, err
}

func Execute() {
	util.CheckErr(rootCmd.Execute())
}
