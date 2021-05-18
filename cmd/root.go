package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "regit",
	Short: "Regit is for git commands that support regular expressions",
	Long: `Write later`,
	Run: func(cmd *cobra.Command, args []string) {
		// ルートコマンドにアクセスした時はhelpなどを出す？
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}