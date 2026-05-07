package storage

import "testing"

func TestValidateAPIBaseURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{name: "https ok", raw: "https://api.example.com/v1", wantErr: false},
		{name: "http ok", raw: "http://localhost:8080/api", wantErr: false},
		{name: "empty", raw: "", wantErr: true},
		{name: "ftp scheme", raw: "ftp://example.com/api", wantErr: true},
		{name: "userinfo", raw: "https://user:pass@example.com/api", wantErr: true},
		{name: "missing host", raw: "https:///path", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateAPIBaseURL(tt.raw)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected: %v", err)
			}
		})
	}
}
