package cmd

import "testing"

func TestGitStatusListParse(t *testing.T) {
	testCase := []struct {
		status          string
		expectFirstItem string
		expectLength    int
	}{
		{
			status: `M  api/test1.go
 M api/test2.go
M  batch/test3.go
M  Makefile
`,
			expectFirstItem: "api/test1.go",
			expectLength:    4,
		},
		{
			status: ` M api/test2.go
M  batch/test3.go
M  Makefile
`,
			expectFirstItem: "api/test2.go",
			expectLength:    3,
		},
	}

	for _, tc := range testCase {
		result, err := GitStatusParse(tc.status)
		if err != nil {
			t.Fatal()
		}
		if result[0] != tc.expectFirstItem {
			t.Fatalf("result[0] should be %s, but %s", tc.expectFirstItem, result[0])
		}
		if len(result) != tc.expectLength {
			t.Fatalf("result length should be %d, but %d", tc.expectLength, len(result))
		}
	}
}

func TestGitStatusListParseError(t *testing.T) {
	testCase := []string{
		"",
		"err",
	}

	for _, tc := range testCase {
		_, err := GitStatusParse(tc)
		if err == nil {
			t.Fatal("err should not be nil")
		}
	}
}

func TestGitBranchParse(t *testing.T) {
	testCase := []struct {
		branches     string
		expectLength int
	}{
		{
			branches: `
* main
  feature/uuusu/fix-bugs
  feature/yse/add-function
  debug/20210528/check
`,
			expectLength: 3,
		},
		{
			branches: `
* master
  feature
  future
  development
  prod
`,
			expectLength: 4,
		},
		{
			branches: `
  dev
  feature
  future
* QA
`,
			expectLength: 3,
		},
		{
			branches: `
* main
`,
			expectLength: 0,
		},
	}

	for _, tc := range testCase {
		result, err := GitBranchParse(tc.branches)
		if err != nil {
			t.Fatal()
		}
		if len(result) != tc.expectLength {
			t.Fatalf("result length should be %d, but %d", tc.expectLength, len(result))
		}
	}
}

func TestGitbranchParseError(t *testing.T) {
	testCase := []string{
		"",
		"er",
		"* ",
		"  ",
	}

	for _, tc := range testCase {
		_, err := GitBranchParse(tc)
		if err == nil {
			t.Fatal("err should not be nil")
		}
	}
}
