package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type GitCmdResult struct {
	result      []string
	executedCmd []string
	success     bool
}

type GitRunner interface {
	Run(*GitCmdExecutor) (*GitCmdResult, error)
}

type GitStatusRunner struct {
}

func (g *GitStatusRunner) Run(gitCmd *GitCmdExecutor) (*GitCmdResult, error) {
	cmd := gitCmd.commandBuilderInBulk("status")
	var r []byte
	r, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	res, err := GitStatusParse(string(r))
	if err != nil {
		return nil, err
	}
	return &GitCmdResult{result: res, success: true}, nil
}

type GitCheckoutRunner struct {
}

func (g *GitCheckoutRunner) Run(gitCmd *GitCmdExecutor) (*GitCmdResult, error) {
	cmd := gitCmd.commandBuilderInBulk("checkout")

	target := strings.Join(gitCmd.executePath, " ")
	fmt.Printf("checkout targets: %s \n", target)
	if !gitCmd.dryRun {
		_, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		fmt.Println("checkout completed")
	} else {
		fmt.Println("To execute the checkout, rerun without the dryrun option")
	}

	executedCmd := []string{cmd.String()}
	return &GitCmdResult{
		result:      nil,
		success:     true,
		executedCmd: executedCmd,
	}, nil
}

type GitBranchRunner struct {
}

func (g *GitBranchRunner) Run(gitCmd *GitCmdExecutor) (*GitCmdResult, error) {
	if !gitCmd.targetIsNeed {
		cmd := gitCmd.commandBuilderInBulk("branch")
		cmdResult, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		result, err := GitBranchParse(string(cmdResult))
		if err != nil {
			return nil, err
		}

		executedCmd := []string{cmd.String()}
		return &GitCmdResult{
			result:      result,
			executedCmd: executedCmd,
			success:     false,
		}, nil
	}

	executedCmds := gitCmd.commandsBuilder("branch")
	//TODO Execute command execution function asynchronously

	var resultExecutedCmds []string
	for _, e := range executedCmds {
		resultExecutedCmds = append(resultExecutedCmds, e.String())
	}

	return &GitCmdResult{
		result:      nil,
		executedCmd: resultExecutedCmds,
		success:     false,
	}, nil
}

type GitCmdExecutor struct {
	targetIsNeed        bool
	target              []string
	targetRegexp        []*regexp.Regexp
	executePath         []string
	combinableOptions   []string
	uncombinableOptions []string
	dryRun              bool
}

func (g *GitCmdExecutor) commandsBuilder(subCmd string) []*exec.Cmd {
	var commands []*exec.Cmd
	for _, e := range g.executePath {
		option := []string{optionsToString(g.combinableOptions)}
		option = append(option, append(g.uncombinableOptions, e)...)
		option = removeEmpty(option)
		if option != nil {
			commands = append(commands, exec.Command("git", subCmd, strings.Join(option, " ")))
			continue
		}
		commands = append(commands, exec.Command("git", subCmd))
	}
	return commands
}

func (g *GitCmdExecutor) commandBuilderInBulk(subCmd string) *exec.Cmd {
	var optionBase []string
	option := append(optionBase, optionsToString(g.combinableOptions))
	option = append(option, append(g.uncombinableOptions, g.executePath...)...)
	option = removeEmpty(option)

	if option != nil {
		return exec.Command("git", subCmd, strings.Join(option, " "))
	}
	return exec.Command("git", subCmd)
}

func removeEmpty(options []string) []string {
	var result []string
	for _, v := range options {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func (g *GitCmdExecutor) ExecuteCmd(runner GitRunner) (*GitCmdResult, error) {
	var executePath []string
	if g.targetIsNeed {
		for _, v := range g.target {
			for _, r := range g.targetRegexp {
				if r.MatchString(v) {
					executePath = append(executePath, v)
				}
			}
		}

		if len(executePath) == 0 {
			return nil, errors.New("there is no target path")
		}
		g.executePath = executePath
	}
	result, err := runner.Run(g)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func newGitCmdExecutor(combinableOptions, target, uncombinableOptions []string, targetRegexp string, targetIsNeed, dryRun bool) (*GitCmdExecutor, error) {
	var targets []string
	if strings.Contains(targetRegexp, ",") {
		targets = strings.Split(targetRegexp, ",")
	} else {
		if targetRegexp != "" {
			targets = append(targets, targetRegexp)
		}
	}

	var targetRegexps []*regexp.Regexp
	for _, v := range targets {
		targetRegexp, err := regexp.Compile(v)
		if err != nil {
			return nil, err
		}
		targetRegexps = append(targetRegexps, targetRegexp)
	}

	for _, v := range combinableOptions {
		if len(v) != 1 {
			return nil, errors.New("combinableOption should be one character")
		}
	}

	return &GitCmdExecutor{
		target:              target,
		targetRegexp:        targetRegexps,
		targetIsNeed:        targetIsNeed,
		combinableOptions:   combinableOptions,
		uncombinableOptions: uncombinableOptions,
		dryRun:              dryRun,
	}, nil
}

func optionsToString(options []string) string {
	if len(options) == 0 {
		return ""
	}
	var option = make([]byte, 0, len(options))
	option = append(option, '-')

	for _, o := range options {
		option = append(option, o...)
	}
	return string(option)
}
