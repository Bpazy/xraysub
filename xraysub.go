package xraysub

import (
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
)

var (
	// buildVer represents 'xraysub' build version
	buildVer string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "xraysub",
		Short: "",
		Long:  ``,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version prints the build information for xraysub",
		Long:  `Version prints the build information for xraysub`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof(buildVer)
		},
	}
)

func Execute() error {
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}
