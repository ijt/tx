package engine

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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
1, 1.5, 0.0, 1.5, false
2, 2.0, 0.0, 2.0, false
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := New()
			var buf bytes.Buffer
			err := e.Run(strings.NewReader(test.input), &buf)
			output := buf.String()
			errStr := fmt.Sprintf("%v", err)
			wantErrStr := fmt.Sprintf("%v", test.wantError)
			if !strings.HasPrefix(errStr, wantErrStr) {
				t.Fatalf("e.Run() returned error %q, want something starting with %q", errStr, wantErrStr)
			}
			if diff := cmp.Diff(test.wantOutput, output); diff != "" {
				t.Errorf("e.Run() gave an unexpected result (-want, +got):\n%s", diff)
			}
		})
	}
}
