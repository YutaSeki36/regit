package cmd

import "testing"

func TestGitCmdExecutor_ExecuteCmdGitStatus(t *testing.T) {
	cmd, _ := newGitCmdExecutor([]string{"s"}, []string{}, []string{}, "", false)
	_, err := cmd.ExecuteCmd(&GitStatusRunner{})
	if err != nil {
		t.Fatal(err)
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
