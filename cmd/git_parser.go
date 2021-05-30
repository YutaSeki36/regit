package cmd

import (
	"errors"
	"strings"
)

func GitStatusParse(status string) ([]string, error) {
	if status == "" {
		return nil, errors.New("status should not be blank")
	}
	statusList := strings.Split(status, "\n")

	var result []string
	for _, v := range statusList {
		if len(v) == 0 {
			continue
		}
		if len(v) < 4 {
			return nil, errors.New("status format is invalid")
		}
		result = append(result, v[3:])
	}

	return result, nil
}

func GitBranchParse(branches string) ([]string, error) {
	if branches == "" {
		return nil, errors.New("branches should not be blank")
	}
	branchList := strings.Split(branches, "\n")

	var result []string
	for _, b := range branchList {
		if len(b) == 0 {
			continue
		}
		if len(b) < 3 {
			return nil, errors.New("executedCmdString format is invalid")
		}
		prefix, branchName := b[:2], b[2:]
		if strings.Contains(prefix, "*") {
			continue
		}

		result = append(result, branchName)
	}

	return result, nil
}

func CheckGitBranchDeleteResult(result string) error {
	if result == "" {
		return errors.New("result text should not be blank")
	}
	if len(result) < 4 {
		return errors.New("executedCmdString format is invalid")
	}
	resultPrefix := result[:5]
	if resultPrefix == "error" {
		return errors.New(result[7:])
	}

	return nil
}
