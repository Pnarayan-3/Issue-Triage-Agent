package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Pnarayan-3/Issue-Triage-Agent/config"
)

type Client struct {
	APIKey            string
	Model             string
	MaxRetries        int
	RetryDelaySeconds int
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

func NewClient(cfg *config.Config) *Client {
	return &Client{
		APIKey:            os.Getenv("GEMINI_API_KEY"),
		Model:             cfg.GeminiModel,
		MaxRetries:        cfg.MaxRetries,
		RetryDelaySeconds: cfg.RetryDelaySeconds,
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

	url := "https://generativelanguage.googleapis.com/v1beta/models/" +
		c.Model +
		":generateContent?key=" +
		c.APIKey

	var responseBody []byte

	for attempt := 1; attempt <= c.MaxRetries; attempt++ {

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

		responseBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		// Retry temporary Gemini errors.
		if resp.StatusCode == http.StatusServiceUnavailable ||
			resp.StatusCode == http.StatusTooManyRequests {

			if attempt == c.MaxRetries {
				return nil, fmt.Errorf(
					"Gemini API failed after %d attempts: status %d: %s",
					c.MaxRetries,
					resp.StatusCode,
					string(responseBody),
				)
			}

			waitSeconds := c.RetryDelaySeconds * attempt

			fmt.Printf(
				"⚠️ Gemini returned status %d. Retrying in %d seconds... (attempt %d/%d)\n",
				resp.StatusCode,
				waitSeconds,
				attempt,
				c.MaxRetries,
			)

			time.Sleep(time.Duration(waitSeconds) * time.Second)

			continue
		}

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

	jsonText, err := extractJSON(generatedText)
	if err != nil {
		return nil, err
	}

	var result TriageResult

	err = json.Unmarshal([]byte(jsonText), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response JSON: %w", err)
	}

	return &result, nil
}

func extractJSON(text string) (string, error) {
	text = strings.TrimSpace(text)

	// Remove markdown code fences if Gemini adds them.
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")

	text = strings.TrimSpace(text)

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")

	if start == -1 || end == -1 || start >= end {
		return "", fmt.Errorf("no valid JSON object found in AI response")
	}

	return text[start : end+1], nil
}