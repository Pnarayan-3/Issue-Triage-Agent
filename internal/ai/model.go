package ai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ListModels() error {
	url := "https://generativelanguage.googleapis.com/v1beta/models?key=" + c.APIKey

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"Gemini API returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	models, ok := result["models"].([]interface{})

	if !ok {
		fmt.Println(string(body))
		return nil
	}

	fmt.Println("Available Gemini models:")

	for _, model := range models {
		m, ok := model.(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := m["name"].(string)
		displayName, _ := m["displayName"].(string)

		fmt.Printf("- %s (%s)\n", name, displayName)
	}

	return nil
}