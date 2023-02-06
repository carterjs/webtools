package graphql_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/carterjs/webtools/graphql"
)

type response struct {
	Value string `json:"value"`
}

func TestQuery(t *testing.T) {
	var tests = map[string]struct {
		host             string
		variables        map[string]any
		response         any
		expectedErr      bool
		expectedResponse *response
	}{
		"no variables": {
			response:         &response{"Hello, world!"},
			expectedResponse: &response{"Hello, world!"},
		},
		"different response type": {
			response: struct {
				Value bool `json:"value"`
			}{},
			expectedErr:      true,
			expectedResponse: nil,
		},
		"bad host": {
			host:        "/notgood",
			expectedErr: true,
		},
		"server error": {
			response:    "fail",
			expectedErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := getServer(tc.response)

			if tc.host == "" {
				tc.host = s.URL
			}
			resp, err := graphql.Query[response](tc.host, "", nil)
			if err != nil && !tc.expectedErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && tc.expectedErr {
				t.Fatal("expected error")
			}

			if !reflect.DeepEqual(tc.expectedResponse, resp) {
				t.Fatalf("expected %#v, found %#v", tc.expectedResponse, resp)
			}
		})
	}
}

func getServer(value any) *httptest.Server {
	type response struct {
		Data any `json:"data"`
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value == "fail" {
			http.Error(w, "intentional failure", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(response{
			Data: value,
		})
	}))
}
