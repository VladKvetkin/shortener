package shortener

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateID(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want *regexp.Regexp
	}{
		{
			name: "valid short URL",
			url:  "https://practicum.yandex.ru/",
			want: regexp.MustCompile(`^.{8}$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := CreateID(tt.url)
			assert.NoError(t, err)
			assert.Regexp(t, tt.want, id)
		})
	}
}
