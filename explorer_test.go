package configo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildSourceKeys(t *testing.T) {
	tests := []struct {
		name    string
		tagLine string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "nothing",
			tagLine: "",
			want:    map[string]string{},
		},
		{
			name:    "one",
			tagLine: "source=key",
			want: map[string]string{
				"source": "key",
			},
		},
		{
			name:    "many",
			tagLine: "json=key,env=KEY,default=value",
			want: map[string]string{
				"json":    "key",
				"env":     "KEY",
				"default": "value",
			},
		},
		{
			name:    "invalid",
			tagLine: "env=KEY,invalid",
			wantErr: true,
		},
		{
			name:    "multiple =",
			tagLine: "source=key=val",
			want: map[string]string{
				"source": "key=val",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := buildSourceKeys(tt.tagLine)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_explorer_exploreWithParents(t *testing.T) {
	tests := []struct {
		name    string
		s       any
		parents map[string][]string
		want    []exploredField
		wantErr bool
	}{
		{
			name:    "error: not a struct pointer",
			s:       struct{}{},
			wantErr: true,
		},
		{
			name: "empty struct",
			s:    &struct{}{},
			want: make([]exploredField, 0),
		},
		{
			name: "without tag",
			s: &struct {
				Test string
			}{},
			want: make([]exploredField, 0),
		},
		{
			name: "private field",
			s: &struct {
				test string `configo:"map=test"`
			}{},
			want: make([]exploredField, 0),
		},
		{
			name: "with parents",
			s: &struct {
				Test string `configo:"map=test"`
			}{},
			parents: map[string][]string{
				"map": {"parent1", "parent2"},
			},
			want: []exploredField{
				{
					name: "Test",
					sourceKeys: map[string]string{
						"map": "test",
					},
					parentSourceKeys: map[string][]string{
						"map": {"parent1", "parent2"},
					},
				},
			},
		},
		{
			name: "with tag",
			s: &struct {
				Test string `configo:"map=test"`
			}{},
			want: []exploredField{
				{
					name: "Test",
					sourceKeys: map[string]string{
						"map": "test",
					},
					parentSourceKeys: map[string][]string{},
				},
			},
		},
		{
			name: "with maleformed tag",
			s: &struct {
				Test string `configo:"invalid"`
			}{},
			wantErr: true,
		},
		{
			name: "multiple fields",
			s: &struct {
				Test  string `configo:"map=test"`
				Test2 string `configo:"map=test"`
			}{},
			want: []exploredField{
				{
					name: "Test",
					sourceKeys: map[string]string{
						"map": "test",
					},
					parentSourceKeys: map[string][]string{},
				},
				{
					name: "Test2",
					sourceKeys: map[string]string{
						"map": "test",
					},
					parentSourceKeys: map[string][]string{},
				},
			},
		},
		{
			name: "with wrong tag",
			s: &struct {
				Test string `wrong:"map=test"`
			}{},
			want: []exploredField{},
		},
		{
			name: "with multiple sources",
			s: &struct {
				Test string `configo:"map=test,env=TEST"`
			}{},
			want: []exploredField{
				{
					name: "Test",
					sourceKeys: map[string]string{
						"map": "test",
						"env": "TEST",
					},
					parentSourceKeys: map[string][]string{},
				},
			},
		},
		{
			name: "with nested struct",
			s: &struct {
				Nested struct {
					Test string `configo:"map=test,env=TEST"`
				} `configo:"map=parent"`
			}{},
			want: []exploredField{
				{
					name: "Test",
					sourceKeys: map[string]string{
						"map": "test",
						"env": "TEST",
					},
					parentSourceKeys: map[string][]string{
						"map": {"parent"},
					},
				},
			},
		},
		{
			name: "with nested error",
			s: &struct {
				Nested struct {
					Test string `configo:"invalid"`
				} `configo:"map=parent"`
			}{},
			wantErr: true,
		},
		{
			name: "with deep nested struct",
			s: &struct {
				Nested struct {
					Nested struct {
						Test string `configo:"map=test,env=TEST"`
					} `configo:"map=parent1"`
				} `configo:"map=parent2"`
			}{},
			want: []exploredField{
				{
					name: "Test",
					sourceKeys: map[string]string{
						"map": "test",
						"env": "TEST",
					},
					parentSourceKeys: map[string][]string{
						"map": {"parent2", "parent1"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := defaultExplorer.exploreWithParents(tt.s, tt.parents)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)

			// ignore value
			for i, f := range got {
				f.value = reflect.Value{}
				got[i] = f
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
