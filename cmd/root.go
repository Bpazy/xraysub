package cmd

import (
	"github.com/Bpazy/xraysub/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	// buildVer represents 'xraysub' build version
	buildVer string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "xraysub",
		Short: "",
		Long: `不 畏 浮 云 遮 望 眼 · 金 睛 如 炬 耀 苍 穹
K E E P   R I D I N G   /   N E V E R   L O O K   B A C K`,
	}
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})

	file, err := getLogFile()
	util.CheckErr(err)
	logrus.SetOutput(file)
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
