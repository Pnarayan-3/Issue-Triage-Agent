package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ConfidenceThreshold float64
}

func Load() (*Config, error) {
	threshold := 0.80

	value := os.Getenv("TRIAGE_CONFIDENCE_THRESHOLD")

	if value != "" {
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid TRIAGE_CONFIDENCE_THRESHOLD: %w",
				err,
			)
		}

		threshold = parsed
	}

	if threshold < 0 || threshold > 1 {
		return nil, fmt.Errorf(
			"TRIAGE_CONFIDENCE_THRESHOLD must be between 0 and 1",
		)
	}

	return &Config{
		ConfidenceThreshold: threshold,
	}, nil
}