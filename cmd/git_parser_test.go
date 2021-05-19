package cmd

import "testing"

func TestGitStatusListParse(t *testing.T) {
	testCase := []struct {
		status string
		expectFirstItem string
		expectLength int
	}{
		{
			status: `M  api/test1.go
 M api/test2.go
M  batch/test3.go
M  Makefile
`,
			expectFirstItem: "api/test1.go",
			expectLength: 4,
		},
		{
			status: ` M api/test2.go
M  batch/test3.go
M  Makefile
`,
			expectFirstItem: "api/test2.go",
			expectLength: 3,
		},
	}

	for _, tc := range testCase {
		result ,err:= GitStatusParse(tc.status)
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
	testCase := []string {
	"",
	"err",
	}

	for _, tc := range testCase {
		_ ,err:= GitStatusParse(tc)
		if err == nil {
			t.Fatal("err should not be nil")
		}
	}
}
