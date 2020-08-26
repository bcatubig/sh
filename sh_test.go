package sh

import (
	"os"
	"testing"
)

func TestCommand_Run(t *testing.T) {
	tests := []struct {
		name       string
		bin        string
		opts       []func(*Command)
		expectedRC int
		wantErr    bool
	}{
		{
			name:    "ls",
			bin:     "ls",
			wantErr: false,
		},
		{
			name: "ls with args",
			bin:  "ls",
			opts: []func(*Command){
				Args("-lha"),
			},
			wantErr: false,
		},
		{
			name: "echo env",
			bin:  "sh",
			opts: []func(*Command){
				Writers(os.Stdout),
				Environment(map[string]string{
					"NAME": "joe",
				}),
				Args(
					"-c",
					"echo \"hello, ${NAME}\"",
				),
			},
			wantErr: false,
		},
		{
			name: "exit with expected non-zero exit code",
			bin:  "sh",
			opts: []func(*Command){
				Args(
					"-c",
					"exit 2",
				),
				ExpectedReturnCode(2),
			},
			expectedRC: 2,
			wantErr:    false,
		},
		{
			name: "simulate exit error",
			bin:  "sh",
			opts: []func(*Command){
				Args(
					"-c",
					"exit 1",
				),
			},
			expectedRC: 1,
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(st *testing.T) {
			c := NewCommand(tc.bin, tc.opts...)
			output, err := c.Run()

			if tc.wantErr {
				if err == nil {
					st.Fatal("wanted error, got nil")
				}
				st.Log(err)
				return
			}

			if err != nil {
				st.Fatal(err)
			}

			if output.ReturnCode != tc.expectedRC {
				st.Fatalf("wrong error code: got %d, want %d", output.ReturnCode, tc.expectedRC)
			}

			st.Log(output)
		})
	}
}
