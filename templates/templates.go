package templates

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

var Walk = filepath.Walk
var ReadFile = os.ReadFile

// ParseDir walks a directory, parsing template files within and minifying their content
func ParseDir(root string, funcs template.FuncMap) (*template.Template, error) {
	tmpl := template.New("").Funcs(funcs)

	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)

	// get and minify templates
	err := Walk("./templates", func(path string, _ fs.FileInfo, _ error) error {
		if filepath.Ext(path) == ".html" {
			b, err := ReadFile(path)
			if err != nil {
				return err
			}

			mb, err := minifier.Bytes("text/html", b)
			if err != nil {
				return err
			}

			_, err = tmpl.Parse(string(mb))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
