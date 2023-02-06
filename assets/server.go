package assets

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

var maxAge = fmt.Sprint(60 * 60 * 24 * 30) // 30 days in seconds

var GetVersion = func() int64 {
	return time.Now().Unix()
}

var GetFileServer = func(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}

type Server struct {
	minifier   *minify.M
	prefix     string
	fileServer http.Handler
	version    int64
	aliases    map[string]string
}

func NewServer(prefix, path string) *Server {
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minifier.AddFunc("image/svg+xml", svg.Minify)
	minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	return &Server{
		minifier:   minifier,
		prefix:     prefix,
		fileServer: http.StripPrefix(prefix, GetFileServer(path)),
		version:    GetVersion(),
		aliases:    make(map[string]string),
	}
}

// GetVersionedPath returns a path to the filename that includes a version prefix
func (server *Server) GetVersionedPath(path string) string {
	dir := filepath.Dir(path)
	file := filepath.Base(path)

	smartCachePath := filepath.Join(dir, fmt.Sprintf("%d_%s", server.version, file))
	server.aliases[smartCachePath] = path

	return smartCachePath
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if actualPath, ok := server.aliases[r.URL.Path]; ok {
		r.URL.Path = actualPath
	}

	w.Header().Set("Cache-Control", "public, max-age="+maxAge)

	// this may not be the best way to do this
	if ext := filepath.Ext(r.URL.Path); ext == ".css" || ext == ".js" || ext == ".svg" {
		mw := server.minifier.ResponseWriter(w, r)
		defer mw.Close()

		w = mw
	}

	server.fileServer.ServeHTTP(w, r)
}
