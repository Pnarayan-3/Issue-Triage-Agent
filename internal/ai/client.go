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

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
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

	// url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + c.APIKey
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.8-flash:generateContent?key=" + c.APIKey
	
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

	var geminiResponse GeminiResponse

	err = json.Unmarshal(responseBody, &geminiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	if len(geminiResponse.Candidates) == 0 {
		return nil, fmt.Errorf("Gemini returned no candidates")
	}

	if len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("Gemini returned no content")
	}

	generatedText := geminiResponse.Candidates[0].Content.Parts[0].Text

	fmt.Println("\nAI generated JSON:")
	fmt.Println(generatedText)

	var result TriageResult

	err = json.Unmarshal([]byte(generatedText), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse triage JSON: %w", err)
	}

	return &result, nil
}