package xraysub

import (
	"github.com/spf13/cobra"
	"log"
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

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version prints the build information for xraysub",
		Long:  `Version prints the build information for xraysub`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println(buildVer)
		},
	}
)

func Execute() error {
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}
