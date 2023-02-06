package templates_test

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"
	"text/template"

	"github.com/carterjs/webtools/templates"
)

func TestParseDir(t *testing.T) {
	var tests = map[string]struct {
		files          map[string][]byte
		expectedResult *template.Template
		expectedErr    bool
	}{
		"valid temmplates": {
			files: map[string][]byte{
				"home.html": []byte(`
					{{ define "home" }}
						
						<h1>Home</h1>
		
					{{ end }}
				`),
			},
			expectedResult: template.Must(template.New("").Parse(`{{ define "home" }}<h1>Home</h1>{{ end }}`)),
		},
		"invalid template": {
			files: map[string][]byte{
				"template.html": []byte(`{{ fail }}`),
			},
			expectedErr: true,
		},
		"file walking error": {
			expectedErr: true,
		},
		"file reading error": {
			files: map[string][]byte{
				"fail.html": {},
			},
			expectedResult: template.New(""),
			expectedErr:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defineFiles(tc.files)

			tmpl, err := templates.ParseDir("", template.FuncMap{})
			if err != nil {
				if !tc.expectedErr {
					t.Fatalf("unexpected error: %v", err)
				}

				return
			} else if tc.expectedErr {
				t.Fatal("expected error")
			}

			if !reflect.DeepEqual(tc.expectedResult.Tree, tmpl.Tree) {
				t.Fatalf("expected %v, found %v", tc.expectedResult.Tree, tmpl.Tree)
			}

		})
	}
}

func defineFiles(files map[string][]byte) {
	templates.ReadFile = func(name string) ([]byte, error) {
		if name == "fail.html" {
			return nil, errors.New("fail")
		}

		if f, ok := files[name]; ok {
			return f, nil
		} else {
			return nil, errors.New("file doesn't exist")
		}
	}

	templates.Walk = func(root string, fn filepath.WalkFunc) error {
		if files == nil {
			return errors.New("fail")
		}

		for path := range files {
			err := fn(path, nil, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
