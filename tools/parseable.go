package tools

import (
	"encoding/json"
	"log"
	"net/http"
)

// These variables must be set by main.go before calling RegisterParseableTools
var (
	ParseableBaseURL string
	ParseableUser    string
	ParseablePass    string
)

func addBasicAuth(req *http.Request) {
	req.SetBasicAuth(ParseableUser, ParseablePass)
}

func listParseableStreams() ([]string, error) {
	url := ParseableBaseURL + "/api/v1/logstream"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	addBasicAuth(httpReq)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
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

func getParseableSchema(stream string) (map[string]string, error) {
	url := ParseableBaseURL + "/api/v1/logstream/" + stream + "/schema"
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	addBasicAuth(httpReq)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	var result struct {
		Fields []struct {
			Name     string          `json:"name"`
			DataType json.RawMessage `json:"data_type"`
		} `json:"fields"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	schema := make(map[string]string)
	for _, field := range result.Fields {
		var dtStr string
		if err := json.Unmarshal(field.DataType, &dtStr); err == nil {
			schema[field.Name] = dtStr
			continue
		}
		schema[field.Name] = string(field.DataType)
	}
	return schema, nil
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
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, nil, err
	}
	return stats, nil, nil
}
