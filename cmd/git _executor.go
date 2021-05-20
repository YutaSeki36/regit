package cmd

type GitCmdExecutor struct {
	subCmd  string
	options []string
}

func (g *GitCmdExecutor) ExecuteCmd() (string, error) {

	return "", nil
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
