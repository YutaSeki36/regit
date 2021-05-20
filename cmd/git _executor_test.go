package cmd

import "testing"

func TestGitCmdExceutor_ExecuteCmd(t *testing.T) {

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
