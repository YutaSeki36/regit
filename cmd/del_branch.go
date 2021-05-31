package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// delBranchCmd represents the delBranch command
var delBranchCmd = &cobra.Command{
	Use:   "del_branch",
	Short: "del_branch is used to delete a branch.",
	Long:  `del_branch is used to delete a branch.`,
	Run:   delBranch,
}

func init() {
	delBranchCmd.PersistentFlags().StringP("target", "t", "", "Set the branch name to be deleted with a regular expression.")
	rootCmd.AddCommand(delBranchCmd)
}

func delBranch(cmd *cobra.Command, args []string) {
	if target, err := cmd.PersistentFlags().GetString("target"); err == nil {
		if target == "" {
			fmt.Println("target should not be blank")
			os.Exit(2)
		}

		var gitBranchResult *GitCmdResult
		// git branch
		{
			cmd, err := newGitCmdExecutor([]string{""}, []string{}, []string{}, "", false, false)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			gitBranchResult, err = cmd.ExecuteCmd(&GitBranchRunner{})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		// git branch -d
		{
			cmd, err := newGitCmdExecutor([]string{"d"}, gitBranchResult.result, []string{}, target, true, dryRun)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			cmd.ExecuteCmd(&GitBranchRunner{})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}
