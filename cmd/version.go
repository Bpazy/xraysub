package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the build information for xraysub",
	Long:  `print the build information for xraysub`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(buildVer)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
