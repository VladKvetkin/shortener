package shortener

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateId(t *testing.T) {
	tests := []struct {
		name string
		want *regexp.Regexp
	}{
		{
			name: "valid short URL",
			want: regexp.MustCompile(`^.{8}$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Regexp(t, tt.want, CreateId())
		})
	}
}
