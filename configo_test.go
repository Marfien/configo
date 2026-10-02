package configo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	t.Setenv("CONFIGO_TEST", "value")

	assert.Equal(t, "env", Environment.Tag())

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr error
	}{
		{
			name: "set",
			key:  "CONFIGO_TEST",
		},
		{
			name:    "unset",
			key:     "UKNOWN",
			wantErr: ErrNotFoundInSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := Environment.GetString(tt.key, nil)

			if tt.wantErr != nil {
				assert.Error(t, gotErr)
				assert.ErrorIs(t, gotErr, tt.wantErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, "value", got)
		})
	}
}

// Default echoes the key back as the value, so the tag doubles as the default literal.
func TestDefault(t *testing.T) {
	assert.Equal(t, "default", Default.Tag())

	got, gotErr := Default.GetInt("42", nil)

	assert.NoError(t, gotErr)
	assert.Equal(t, 42, got)
}

func TestParse(t *testing.T) {
	t.Setenv("CONFIGO_TEST_HOST", "example.com")

	cfg := &struct {
		Host    string `configo:"env=CONFIGO_TEST_HOST,default=localhost"`
		Port    int    `configo:"env=CONFIGO_TEST_PORT,default=8080"`
		Verbose bool   `configo:"default=false"`
	}{}

	err := Parse(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "example.com", cfg.Host) // env wins over default
	assert.Equal(t, 8080, cfg.Port)          // env unset, default applies
	assert.False(t, cfg.Verbose)
}

func TestParse_missingValue(t *testing.T) {
	cfg := &struct {
		Host string `configo:"env=CONFIGO_TEST_MISSING"`
	}{}

	assert.ErrorContains(t, Parse(cfg), "Host")
}
