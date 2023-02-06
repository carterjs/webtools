package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Query makes a GraphQL request and parses the response into a type
func Query[T any](url string, q string, variables map[string]any) (*T, error) {
	bodyBytes, err := json.Marshal(map[string]any{
		"query":     q,
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errBody string
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("error response from api: %v", resp.Status)
	}

	var responseBody struct {
		Data *T `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	return responseBody.Data, nil
}
