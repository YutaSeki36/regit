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
		if len(v) < 3 {
			return nil, errors.New("status format is invalid")
		}
		result = append(result, v[3:])
	}

	return result, nil
}
