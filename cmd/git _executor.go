package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type GitCmdResult struct {
	result  []string
	success bool
}

type GitRunner interface {
	Run(*GitCmdExecutor) (*GitCmdResult, error)
}

type GitStatusRunner struct {
}

func (g *GitStatusRunner) Run(gitCmd *GitCmdExecutor) (*GitCmdResult, error) {
	cmd := commandBuilder(gitCmd, "status", gitCmd.executePath)
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

func commandBuilder(gitCmd *GitCmdExecutor, subCmd string, target []string) *exec.Cmd {
	t := strings.Join(target, " ")
	if t == "" {
		return exec.Command("git", subCmd, optionsToString(gitCmd.combinableOptions))
	}

	return exec.Command("git", subCmd, optionsToString(gitCmd.combinableOptions), strings.Join(target, " "))
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
			fmt.Println("there is no target path")
			return nil, nil
		}
		g.executePath = executePath
	}
	result, err := runner.Run(g)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func newGitCmdExecutor(combinableOptions, target, uncombinableOptions []string, targetRegexp string, targetIsNeed bool) (*GitCmdExecutor, error) {
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
