package api

import (
	"bytes"
	"fmt"
	"net/http"
)

type SnykOptions struct {
	GroupId   string
	ApiKey    string
	UserAgent string
}

func clientDo(so SnykOptions, method string, path string, body []byte) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, constructUrl(path), bytes.NewReader(body))

	generateHeaders(so, req)

	return client.Do(req)
}

func generateHeaders(so SnykOptions, req *http.Request) {
	authToken := fmt.Sprintf("token %s", so.ApiKey)
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", so.UserAgent)
}

func constructUrl(path string) string {
	snykEndpoint := "https://snyk.io/api/v1%s"
	return fmt.Sprintf(snykEndpoint, path)
}
