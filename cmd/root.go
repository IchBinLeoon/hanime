package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "hanime",
	Short: "Command-line tool to download videos from hanime.tv",
	Long:  "Command-line tool to download videos from hanime.tv\n\nComplete documentation is available at https://github.com/IchBinLeoon/hanime",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
