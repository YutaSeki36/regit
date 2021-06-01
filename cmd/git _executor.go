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

type asyncRunnerResult struct {
	executedCmdString string
	error             error
}

type GitBranchRunner struct {
}

const concurrencyLimitNum = 10

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

	target := strings.Join(gitCmd.executePath, " ")
	fmt.Printf("branch delete targets: %s \n", target)

	executedCmds := gitCmd.commandsBuilder("branch")
	if !gitCmd.dryRun {
		concurrencyLimitSignal := make(chan struct{}, concurrencyLimitNum)
		responseCh := make(chan *asyncRunnerResult, len(executedCmds))
		defer close(concurrencyLimitSignal)
		defer close(responseCh)

		for _, executedCmd := range executedCmds {
			go gitBranchAsyncRunner(executedCmd, responseCh, concurrencyLimitSignal)
		}

		for i := 0; i < len(executedCmds); i++ {
			chi := <-responseCh
			if chi.error != nil {
				fmt.Printf("Failed: %s. Reason: %s", chi.executedCmdString, chi.error.Error())
				fmt.Println()
			}
		}

		fmt.Println("Finished")
	} else {
		fmt.Println("To execute the del_branch, rerun without the dryrun option")
	}

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

func gitBranchAsyncRunner(cmd *exec.Cmd, result chan<- *asyncRunnerResult, sig chan struct{}) {
	sig <- struct{}{}
	cmdResult, err := cmd.Output()
	if err != nil {
		result <- &asyncRunnerResult{
			error:             err,
			executedCmdString: cmd.String(),
		}
		<-sig
		return
	}
	if err := CheckGitBranchDeleteResult(string(cmdResult)); err != nil {
		result <- &asyncRunnerResult{
			error:             err,
			executedCmdString: cmd.String(),
		}
		<-sig
		return
	}

	result <- &asyncRunnerResult{
		error:             nil,
		executedCmdString: cmd.String(),
	}
	<-sig
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
		option := []string{subCmd}
		option = append(option, optionsToString(g.combinableOptions))
		option = append(option, append(g.uncombinableOptions, e)...)
		option = removeEmpty(option)
		if len(option) != 0 {
			commands = append(commands, exec.Command("git", option...))
			continue
		}
		commands = append(commands, exec.Command("git", subCmd))
	}
	return commands
}

func (g *GitCmdExecutor) commandBuilderInBulk(subCmd string) *exec.Cmd {
	option := []string{subCmd}
	option = append(option, optionsToString(g.combinableOptions))
	option = append(option, append(g.uncombinableOptions, g.executePath...)...)
	option = removeEmpty(option)

	if option != nil {
		return exec.Command("git", option...)
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
		ep := map[string]struct{}{}
		concurrencyLimitSignal := make(chan struct{}, concurrencyLimitNum)
		responseCh := make(chan string, len(g.target)*len(g.targetRegexp))
		defer close(concurrencyLimitSignal)
		defer close(responseCh)
		for _, v := range g.target {
			for _, r := range g.targetRegexp {
				go func(r *regexp.Regexp, target string, response chan<- string, sig chan struct{}) {
					sig <- struct{}{}
					if r.MatchString(target) {
						responseCh <- target
					} else {
						responseCh <- ""
					}
					<-sig
				}(r, v, responseCh, concurrencyLimitSignal)
			}
		}

		for i := 0; i < len(g.target)*len(g.targetRegexp); i++ {
			res := <-responseCh
			if res != "" {
				ep[res] = struct{}{}
			}
		}

		for k := range ep {
			executePath = append(executePath, k)
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
		if len(v) != 1 && v != "" {
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
	result := string(option)
	if result == "-" {
		return ""
	}
	return result
}
