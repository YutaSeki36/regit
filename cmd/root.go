package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "regit",
	Short: "Regit is for git commands that support regular expressions",
	Long:  `Regit is for git commands that support regular expressions`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var dryRun bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dryRun", "d", false, "dryRun enable flag")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
