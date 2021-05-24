package cmd

import (
	"testing"
)

func TestGitCmdExecutor_ExecuteCmdGitStatus(t *testing.T) {
	cmd, _ := newGitCmdExecutor([]string{"s"}, []string{}, []string{}, "", false, true)
	_, err := cmd.ExecuteCmd(&GitStatusRunner{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGitCmdExecutor_ExecuteCmdGitCheckout(t *testing.T) {
	testCase := []struct {
		name                string
		targetIsNeed        bool
		combinableOptions   []string
		target              []string
		uncombinableOptions []string
		targetRegexp        string
		expectCmd           string
	}{
		{
			name:                "git checkout api/hoge/fuga.go",
			targetIsNeed:        true,
			combinableOptions:   []string{},
			uncombinableOptions: []string{},
			target:              []string{"cmd/hoge.go", "api/hoge/fuga.go", "api/hoge/fuga.json"},
			targetRegexp:        "api/hoge/.+go",
			expectCmd:           "/usr/local/bin/git checkout api/hoge/fuga.go",
		},
		{
			name:                "git checkout api/hoge/fuga.go api/hoge/hoge.go api/hoge_generated.go",
			targetIsNeed:        true,
			combinableOptions:   []string{},
			uncombinableOptions: []string{},
			target:              []string{"cmd/hoge.go", "api/hoge/fuga.go", "api/hoge/fuga.json", "api/hoge/hoge.go", "api/hoge_generated.go"},
			targetRegexp:        "api/hoge/.+go,api/.+_generated\\.go",
			expectCmd:           "/usr/local/bin/git checkout api/hoge/fuga.go api/hoge/hoge.go api/hoge_generated.go",
		},
		{
			name:                "git checkout --theirs api/hoge/fuga.go",
			targetIsNeed:        true,
			combinableOptions:   []string{},
			uncombinableOptions: []string{"--theirs"},
			target:              []string{"cmd/hoge.go", "api/hoge/fuga.go", "api/hoge/fuga.json", "api/hoge/hoge.go", "api/hoge_generated.go"},
			targetRegexp:        "api/hoge/.+go",
			expectCmd:           "/usr/local/bin/git checkout --theirs api/hoge/fuga.go api/hoge/hoge.go",
		},
		{
			name:                "git checkout --ours api/hoge/fuga.go",
			targetIsNeed:        true,
			combinableOptions:   []string{},
			uncombinableOptions: []string{"--ours"},
			target:              []string{"cmd/hoge.go", "api/hoge/fuga.go", "api/hoge/fuga.json", "api/hoge/hoge.go", "api/hoge_generated.go"},
			targetRegexp:        "api/hoge/.+go",
			expectCmd:           "/usr/local/bin/git checkout --ours api/hoge/fuga.go api/hoge/hoge.go",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			cmd, _ := newGitCmdExecutor(tc.combinableOptions, tc.target, tc.uncombinableOptions, tc.targetRegexp, tc.targetIsNeed, false)
			result, err := cmd.ExecuteCmd(&GitCheckoutRunner{})
			if err != nil {
				t.Fatal(err)
			}
			if result.executedCmd != tc.expectCmd {
				t.Fatalf("executed command should be [%s], but [%s]", tc.expectCmd, result.executedCmd)
			}
		})
	}
}

func TestGitCmdExecutor_optionToString(t *testing.T) {
	testCase := []struct {
		options []string
		expect  string
	}{
		{
			options: []string{
				"f",
				"b",
			},
			expect: "-fb",
		},
		{
			options: []string{
				"a",
			},
			expect: "-a",
		},
		{
			options: []string{},
			expect:  "",
		},
		{
			options: []string{
				"f",
				"b",
				"tc",
			},
			expect: "-fbtc",
		},
	}

	for _, tc := range testCase {
		o := optionsToString(tc.options)
		if o != tc.expect {
			t.Fatalf("result should be %s, but %s", tc.expect, o)
		}
	}
}
