package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

func TestCSVs(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput string
		wantError  error
	}{
		{
			name:       "empty",
			input:      "",
			wantOutput: "",
			wantError:  fmt.Errorf("missing header"),
		},
		{
			name:       "bad header",
			input:      "foo, bar, baz, qux, quax",
			wantOutput: "",
			wantError:  fmt.Errorf("invalid header"),
		},
		{
			name:       "empty body",
			input:      "type, client, tx, amount",
			wantOutput: "client, available, held, total, locked\n",
		},
		{
			name: "example from requirements doc",
			input: `type, client, tx, amount
deposit, 1, 1, 1.0
deposit, 2, 2, 2.0
deposit, 1, 3, 2.0
withdrawal, 1, 4, 1.5
withdrawal, 2, 5, 3.0 
`,
			wantOutput: `client, available, held, total, locked
1, 1.5000, 0.0000, 1.5000, false
2, -1.0000, 0.0000, -1.0000, false
`,
		},
		{
			name: "dispute",
			input: `type, client, tx, amount
deposit, 1, 1, 2.0
dispute, 1, 1			
`,
			wantOutput: `client, available, held, total, locked
1, 0, 2.0000, 2.0000, false
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := New()
			var buf bytes.Buffer
			err := e.Run(strings.NewReader(test.input), &buf)
			if test.wantError == nil {
				if err != nil {
					t.Fatalf("e.Run() returned error %q", err)
				}
			} else {
				errStr := fmt.Sprintf("%v", err)
				wantErrStr := fmt.Sprintf("%v", test.wantError)
				if !strings.HasPrefix(errStr, wantErrStr) {
					t.Fatalf("e.Run() returned error %q, want something starting with %q", errStr, wantErrStr)
				}
			}
			output := buf.String()
			if test.wantOutput != output {
				diff, err := fileDiff(test.wantOutput, output)
				if err != nil {
					t.Fatalf("diffing results: %v", err)
				}
				t.Errorf("e.Run() gave an unexpected result (-want, +got):\n%s", diff)
			}
		})
	}
}

func fileDiff(want, got string) (string, error) {
	wantFile, err := ioutil.TempFile("/tmp", "want")
	if err != nil {
		return "", fmt.Errorf("making first temp file: %v", err)
	}

	gotFile, err := ioutil.TempFile("/tmp", "got")
	if err != nil {
		return "", fmt.Errorf("making second temp file: %v", err)
	}

	if nw, err := wantFile.WriteString(want); err != nil {
		return "", fmt.Errorf("writing contents of wantFile (%d bytes written): %v", nw, err)
	}
	if nw, err := gotFile.WriteString(got); err != nil {
		return "", fmt.Errorf("writing contents of gotFile (%d bytes written): %v", nw, err)
	}

	cmd := exec.Command("diff", "-u", wantFile.Name(), gotFile.Name())
	// err is non-nil if there is a non-zero diff, so just ignore it.
	out, _ := cmd.CombinedOutput()
	return string(out), nil
}
