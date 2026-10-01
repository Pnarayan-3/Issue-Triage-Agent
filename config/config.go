package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ConfidenceThreshold float64
	GeminiModel         string
	MaxRetries          int
	RetryDelaySeconds   int
}

func Load() (*Config, error) {
	confidenceThreshold := 0.80
	geminiModel := "YOUR_CURRENT_WORKING_MODEL"
	maxRetries := 3
	retryDelaySeconds := 1

	if value := os.Getenv("TRIAGE_CONFIDENCE_THRESHOLD"); value != "" {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid TRIAGE_CONFIDENCE_THRESHOLD: %w",
				err,
			)
		}

		confidenceThreshold = parsed
	}

	if value := os.Getenv("GEMINI_MODEL"); value != "" {
		geminiModel = value
	}

	if value := os.Getenv("GEMINI_MAX_RETRIES"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid GEMINI_MAX_RETRIES: %w",
				err,
			)
		}

		maxRetries = parsed
	}

	if value := os.Getenv("GEMINI_RETRY_DELAY_SECONDS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid GEMINI_RETRY_DELAY_SECONDS: %w",
				err,
			)
		}

		retryDelaySeconds = parsed
	}

	if confidenceThreshold < 0 || confidenceThreshold > 1 {
		return nil, fmt.Errorf(
			"TRIAGE_CONFIDENCE_THRESHOLD must be between 0 and 1",
		)
	}

	if maxRetries < 1 {
		return nil, fmt.Errorf(
			"GEMINI_MAX_RETRIES must be at least 1",
		)
	}

	if retryDelaySeconds < 0 {
		return nil, fmt.Errorf(
			"GEMINI_RETRY_DELAY_SECONDS cannot be negative",
		)
	}

	return &Config{
		ConfidenceThreshold: confidenceThreshold,
		GeminiModel:         geminiModel,
		MaxRetries:          maxRetries,
		RetryDelaySeconds:   retryDelaySeconds,
	}, nil
}