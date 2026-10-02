package configo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setOrAppendToString(t *testing.T) {
	tests := []struct {
		name  string
		m     map[string]string
		want  map[string]string
		key   string
		value string
	}{
		{
			name:  "new key",
			m:     make(map[string]string),
			key:   "key",
			value: "value",
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name: "append to existing",
			m: map[string]string{
				"key": "first",
			},
			key:   "key",
			value: "second",
			want: map[string]string{
				"key": "first,second",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setOrAppendToString(tt.m, tt.key, tt.value)
			assert.Equal(t, tt.m, tt.want)
		})
	}
}

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want map[string]string
	}{
		{
			name: "no args",
			args: make([]string, 0),
			want: make(map[string]string),
		},
		{
			name: "key val with equal sign",
			args: []string{
				"--key=value",
			},
			want: map[string]string{
				"--key": "value",
			},
		},
		{
			name: "value contains equal sign",
			args: []string{
				"--key=k=v",
			},
			want: map[string]string{
				"--key": "k=v",
			},
		},
		{
			name: "key val with space",
			args: []string{"--key", "value"},
			want: map[string]string{
				"--key": "value",
			},
		},
		{
			name: "flag",
			args: []string{"--flag"},
			want: map[string]string{
				"--flag": "true",
			},
		},
		{
			name: "flag and kv",
			args: []string{"--flag", "--key", "value"},
			want: map[string]string{
				"--flag": "true",
				"--key":  "value",
			},
		},
		{
			name: "trailing flag",
			args: []string{"--key", "value", "--flag"},
			want: map[string]string{
				"--key":  "value",
				"--flag": "true",
			},
		},
		{
			name: "repeated key is joined",
			args: []string{"--key", "a", "--key=b"},
			want: map[string]string{
				"--key": "a,b",
			},
		},
		{
			name: "negative value is treated as the next key",
			args: []string{"--key", "-1"},
			want: map[string]string{
				"--key": "true",
				"-1":    "true",
			},
		},
		{
			name: "positional arg becomes a flag",
			args: []string{"value"},
			want: map[string]string{
				"value": "true",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseArgs(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}
