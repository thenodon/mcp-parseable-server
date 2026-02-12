package tools

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

// These variables must be set by main.go before calling RegisterParseableTools
var (
	ParseableBaseURL string
	ParseableUser    string
	ParseablePass    string
)

// package-level HTTP client; initialized in init() to respect UNSECURE env var
var HTTPClient *http.Client

func init() {
	// UNSECURE environment variable controls whether TLS verification is skipped.
	// Accepts the same values as strconv.ParseBool (true/1/t etc.).
	unsecureEnv := os.Getenv("UNSECURE")
	if ok, _ := strconv.ParseBool(unsecureEnv); ok {
		HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		slog.Info("UNSECURE=true: HTTP client will skip TLS verification")
	} else {
		HTTPClient = http.DefaultClient
	}
}

func addBasicAuth(req *http.Request) {
	req.SetBasicAuth(ParseableUser, ParseablePass)
}

func doParseableQuery(query string, streamName string, startTime string, endTime string) ([]map[string]interface{}, error) {
	payload := map[string]string{
		"query":      query,
		"streamName": streamName,
		"startTime":  startTime,
		"endTime":    endTime,
	}
	jsonPayload, _ := json.Marshal(payload)
	url := ParseableBaseURL + parseableSQLPath
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	addBasicAuth(httpReq)
	resp, err := HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	body, _ := io.ReadAll(resp.Body)

	// Try to unmarshal as array of rows
	var arrResult []map[string]interface{}
	if err := json.Unmarshal(body, &arrResult); err != nil {
		return nil, err
	}
	return arrResult, nil
}

func listParseableStreams() ([]string, error) {
	url := ParseableBaseURL + "/api/v1/logstream"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	addBasicAuth(httpReq)
	resp, err := HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	var apiResult []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResult); err != nil {
		return nil, err
	}
	streams := make([]string, 0, len(apiResult))
	for _, s := range apiResult {
		streams = append(streams, s.Name)
	}
	return streams, nil
}

func getParseableSchema(stream string) (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/logstream/" + stream + "/schema"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	addBasicAuth(httpReq)
	resp, err := HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func getParseableStats(streamName string) (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/logstream/" + streamName + "/stats"
	stats, m, err := doSimpleGet(url)
	if err != nil {
		return m, err
	}
	return stats, nil
}

func getParseableInfo(streamName string) (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/logstream/" + streamName + "/info"
	info, m, err := doSimpleGet(url)
	if err != nil {
		return m, err
	}
	return info, nil
}

func getParseableAbout() (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/about"
	about, m, err := doSimpleGet(url)
	if err != nil {
		return m, err
	}
	return about, nil
}

func getParseableRoles() (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/roles"
	roles, m, err := doSimpleGet(url)
	if err != nil {
		return m, err
	}
	return roles, nil
}

func doSimpleGet(url string) (map[string]interface{}, map[string]interface{}, error) {
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	addBasicAuth(httpReq)
	resp, err := HTTPClient.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()
	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, nil, err
	}
	return stats, nil, nil
}
