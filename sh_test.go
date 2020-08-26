package sh

import "testing"

func TestCommand_Run(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "basic",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(st *testing.T) {
			st.Log("test")
		})
	}
}
