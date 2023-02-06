package assets_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/carterjs/webtools/assets"
)

func TestServer_ServeHTTP(t *testing.T) {
	var tests = map[string]struct {
		name           string
		files          map[string][]byte
		requestPath    string
		expectedStatus int
		expectedBody   []byte
	}{
		"CSS minification": {
			files: map[string][]byte{
				"/style.css": []byte(`
					/*
						This should be shortened!
					*/

					body {
						background-color: red;


						
					}
				`),
			},
			requestPath:    "/assets/style.css",
			expectedStatus: 200,
			expectedBody:   []byte("body{background-color:red}"),
		},
		"JS minification": {
			files: map[string][]byte{
				"/js/script.js": []byte(`
					// some unnecessary comment


					console.log("done");
				`),
			},
			requestPath:    "/assets/js/script.js",
			expectedStatus: 200,
			expectedBody:   []byte(`console.log("done")`),
		},
		"text without minification": {
			files: map[string][]byte{
				"/text.txt": []byte(`Hello, world!`),
			},
			requestPath:    "/assets/text.txt",
			expectedStatus: 200,
			expectedBody:   []byte(`Hello, world!`),
		},
		"file without extension": {
			files: map[string][]byte{
				"/test": []byte("test"),
			},
			requestPath:    "/assets/test",
			expectedStatus: 200,
			expectedBody:   []byte("test"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assets.GetFileServer = getFileServerGetter(tc.files)

			server := httptest.NewServer(assets.NewServer("/assets", "some/fake/path"))

			resp, err := http.Get(server.URL + tc.requestPath)
			if err != nil {
				t.Fatalf("error in request: %v", err)
			}

			if resp.StatusCode != tc.expectedStatus {
				t.Fatalf("expected %d, found %d", tc.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("error reading body: %v", err)
			}

			if !reflect.DeepEqual(body, tc.expectedBody) {
				t.Fatalf("expected %s, found %s", tc.expectedBody, body)
			}
		})
	}
}

func TestServer_GetVersionedPath(t *testing.T) {
	var tests = map[string]struct {
		name                  string
		serverVersion         int64
		files                 map[string][]byte
		inputPath             string
		expectedVersionedPath string
		requestPath           string
		expectedStatus        int
		expectedBody          []byte
	}{
		"versioned path returns file": {
			files: map[string][]byte{
				"/style.css": []byte(`
					body {
						background-color: red;
					}
				`),
			},
			serverVersion:         0,
			inputPath:             "/assets/style.css",
			expectedVersionedPath: "/assets/0_style.css",
			requestPath:           "/assets/0_style.css",
			expectedStatus:        200,
			expectedBody:          []byte("body{background-color:red}"),
		},
		"original path still works": {
			files: map[string][]byte{
				"/style.css": []byte(`
					body {
						background-color: red;
					}
				`),
			},
			serverVersion:         123,
			inputPath:             "/assets/style.css",
			expectedVersionedPath: "/assets/123_style.css",
			requestPath:           "/assets/style.css",
			expectedStatus:        200,
			expectedBody:          []byte(`body{background-color:red}`),
		},
		"file without extension": {
			files: map[string][]byte{
				"/test": []byte("test"),
			},
			serverVersion:         0,
			inputPath:             "/assets/test",
			expectedVersionedPath: "/assets/0_test",
			requestPath:           "/assets/0_test",
			expectedStatus:        200,
			expectedBody:          []byte("test"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assets.GetFileServer = getFileServerGetter(tc.files)
			assets.GetVersion = func() int64 {
				return tc.serverVersion
			}

			assetsServer := assets.NewServer("/assets", "some/fake/path")
			server := httptest.NewServer(assetsServer)

			versionedPath := assetsServer.GetVersionedPath(tc.inputPath)
			if versionedPath != tc.expectedVersionedPath {
				t.Fatalf("expected %s, found %s", tc.expectedVersionedPath, versionedPath)
			}

			resp, err := http.Get(server.URL + tc.requestPath)
			if err != nil {
				t.Fatalf("error in request: %v", err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("error reading body: %v", err)
			}

			if !reflect.DeepEqual(body, tc.expectedBody) {
				t.Fatalf("expected %s, found %s", tc.expectedBody, body)
			}
		})
	}
}

func getFileServerGetter(files map[string][]byte) func(path string) http.Handler {
	return func(path string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Looking for", r.URL.Path)
			if b, ok := files[r.URL.Path]; ok {
				w.WriteHeader(http.StatusOK)
				w.Write(b)
			} else {
				// file not found
				w.WriteHeader(http.StatusNotFound)
			}
		})
	}
}
