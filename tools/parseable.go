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

func listParseableStreams() ([]map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/logstream"
	return doSimpleGetArray(url)
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
	stats, _, err := doSimpleGet(url)
	return stats, err
}

func getParseableInfo(streamName string) (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/logstream/" + streamName + "/info"
	info, _, err := doSimpleGet(url)
	return info, err
}

func getParseableAbout() (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/about"
	about, _, err := doSimpleGet(url)
	return about, err
}

func getParseableRoles() (map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/roles"
	roles, _, err := doSimpleGet(url)
	return roles, err
}

func getParseableUsers() ([]map[string]interface{}, error) {
	url := ParseableBaseURL + "/api/v1/users"
	return doSimpleGetArray(url)
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
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, err
	}
	return response, nil, nil
}

func doSimpleGetArray(url string) ([]map[string]interface{}, error) {
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
	var response []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
