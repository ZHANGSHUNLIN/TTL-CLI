package opener

import (
	"strings"
	"testing"
)

func TestTarget(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		wantErr string
	}{
		{name: "markdown link", value: "[example](https://example.com)", want: "https://example.com"},
		{name: "markdown link trims target", value: "[example]( http://example.com )", want: "http://example.com"},
		{name: "plain URL", value: "https://example.com", want: "https://example.com"},
		{name: "plain value trims surrounding whitespace", value: "  https://example.com  ", want: "https://example.com"},
		{name: "empty", wantErr: "resource value is empty"},
		{name: "incomplete markdown", value: "[example](https://example.com", wantErr: "valid Markdown link"},
		{name: "missing destination", value: "[example]", wantErr: "valid Markdown link"},
		{name: "empty destination", value: "[example]()", wantErr: "empty Markdown link target"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Target(tt.value)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Target() error = %v", err)
				}
				if got != tt.want {
					t.Fatalf("Target() = %q, want %q", got, tt.want)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Target() error = %v, want substring %q", err, tt.wantErr)
			}
		})
	}
}
