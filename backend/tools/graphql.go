package tools

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func GQLQuery(url string, query string, variables map[string]string, bucket interface{}) error {
	data, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, bucket)
	if err != nil {
		return err
	}

	return nil
}
