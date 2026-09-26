package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	APIKey string
}

func NewClient() *Client {
	return &Client{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func (c *Client) Analyze(title string, body string) (*TriageResult, error) {

	if c.APIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	prompt := fmt.Sprintf(`
%s

GitHub Issue:

Title:
%s

Body:
%s
`, SystemPrompt, title, body)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + c.APIKey

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Gemini API returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	fmt.Println("Gemini raw response:")
	fmt.Println(string(responseBody))

	return nil, nil
}