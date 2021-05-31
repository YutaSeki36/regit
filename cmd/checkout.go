package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "checkout is used for file restore",
	Long:  `checkout is used for file restore`,
	Run:   checkout,
}

var theirs bool
var ours bool

func init() {
	checkoutCmd.PersistentFlags().StringP("target", "t", "", "Set the target file name to check out with a regular expression.")
	checkoutCmd.PersistentFlags().BoolVar(&theirs, "theirs", false, "")
	checkoutCmd.PersistentFlags().BoolVar(&ours, "ours", false, "")

	rootCmd.AddCommand(checkoutCmd)
}

func checkout(cmd *cobra.Command, args []string) {
	if target, err := cmd.PersistentFlags().GetString("target"); err == nil {
		if target == "" {
			fmt.Println("target should not be blank")
			os.Exit(2)
		}

		var gitStatusResult *GitCmdResult
		// git checkout -s
		{
			cmd, err := newGitCmdExecutor([]string{"s"}, []string{}, []string{}, "", false, false)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			gitStatusResult, err = cmd.ExecuteCmd(&GitStatusRunner{})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		// git checkout
		{
			var options []string
			if theirs && ours {
				fmt.Println("you cannot select both theirs and ours option")
				os.Exit(1)
			}
			if theirs {
				options = append(options, "--theirs")
			}
			if ours {
				options = append(options, "--ours")
			}

			cmd, err := newGitCmdExecutor([]string{}, gitStatusResult.result, options, target, true, dryRun)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			cmd.ExecuteCmd(&GitCheckoutRunner{})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}
