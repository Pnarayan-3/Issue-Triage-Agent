package ai

import "testing"

func TestExtractJSON_ValidJSON(t *testing.T) {
	input := `{"type":"BUG","priority":"HIGH"}`

	result, err := extractJSON(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != input {
		t.Errorf("expected %s, got %s", input, result)
	}
}

func TestExtractJSON_MarkdownCodeFence(t *testing.T) {
	input := "```json\n{\"type\":\"BUG\",\"priority\":\"HIGH\"}\n```"

	expected := `{"type":"BUG","priority":"HIGH"}`

	result, err := extractJSON(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestExtractJSON_ExtraText(t *testing.T) {
	input := `Here is the triage result:
{"type":"BUG","priority":"HIGH"}
Please review the result.`

	expected := `{"type":"BUG","priority":"HIGH"}`

	result, err := extractJSON(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestExtractJSON_Whitespace(t *testing.T) {
	input := "  \n  {\"type\":\"BUG\"}  \n  "
	expected := `{"type":"BUG"}`

	result, err := extractJSON(input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestExtractJSON_NoJSON(t *testing.T) {
	input := "The issue appears to be a backend bug."

	_, err := extractJSON(input)

	if err == nil {
		t.Fatal("expected an error when no JSON object is present")
	}
}

func TestExtractJSON_IncompleteJSON(t *testing.T) {
	input := `{"type":"BUG","priority":"HIGH"`

	_, err := extractJSON(input)

	if err == nil {
		t.Fatal("expected an error for incomplete JSON")
	}
}

func TestExtractJSON_EmptyInput(t *testing.T) {
	input := ""

	_, err := extractJSON(input)

	if err == nil {
		t.Fatal("expected an error for empty input")
	}
}