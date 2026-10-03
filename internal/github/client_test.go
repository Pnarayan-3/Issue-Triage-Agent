package github

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	handler func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.handler(req)
}

func mockResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestGetIssue_Success(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}

			if !strings.Contains(req.URL.Path, "/issues/123") {
				t.Errorf("unexpected URL: %s", req.URL.Path)
			}

			return mockResponse(
				http.StatusOK,
				`{
					"number": 123,
					"title": "Login API is failing",
					"body": "Users receive HTTP 500"
				}`,
			), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	issue, err := client.GetIssue("123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if issue.ID != 123 {
		t.Errorf("expected issue ID 123, got %d", issue.ID)
	}

	if issue.Title != "Login API is failing" {
		t.Errorf("unexpected title: %s", issue.Title)
	}

	if issue.Body != "Users receive HTTP 500" {
		t.Errorf("unexpected body: %s", issue.Body)
	}
}

func TestGetIssue_APIError(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {
			return mockResponse(
				http.StatusNotFound,
				`{"message":"Not Found"}`,
			), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	_, err := client.GetIssue("999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "status 404") {
		t.Errorf("expected 404 error, got %v", err)
	}
}

func TestAddComment_Success(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}

			if !strings.Contains(req.URL.Path, "/issues/123/comments") {
				t.Errorf("unexpected URL: %s", req.URL.Path)
			}

			return mockResponse(http.StatusCreated, `{}`), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	err := client.AddComment("123", "Automated triage result")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAddLabels_Success(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}

			if !strings.Contains(req.URL.Path, "/issues/123/labels") {
				t.Errorf("unexpected URL: %s", req.URL.Path)
			}

			return mockResponse(http.StatusOK, `[]`), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	err := client.AddLabels(
		"123",
		[]string{"bug", "backend", "review:auto"},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateLabel_Success(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", req.Method)
			}

			if !strings.Contains(req.URL.Path, "/labels") {
				t.Errorf("unexpected URL: %s", req.URL.Path)
			}

			return mockResponse(http.StatusCreated, `{}`), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	err := client.CreateLabel("backend")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetComments_Success(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", req.Method)
			}

			return mockResponse(
				http.StatusOK,
				`[
					{
						"id": 101,
						"body": "First comment"
					},
					{
						"id": 102,
						"body": "Second comment"
					}
				]`,
			), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	comments, err := client.GetComments("123")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}

	if comments[0].ID != 101 {
		t.Errorf("expected first comment ID 101, got %d", comments[0].ID)
	}
}

func TestGetTriageComment_Found(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			return mockResponse(
				http.StatusOK,
				`[
					{
						"id": 101,
						"body": "Normal comment"
					},
					{
						"id": 102,
						"body": "<!-- issue-triage-agent -->\nAutomated triage"
					}
				]`,
			), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	comment, err := client.GetTriageComment(
		"123",
		"<!-- issue-triage-agent -->",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if comment == nil {
		t.Fatal("expected triage comment, got nil")
	}

	if comment.ID != 102 {
		t.Errorf("expected comment ID 102, got %d", comment.ID)
	}
}

func TestGetTriageComment_NotFound(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			return mockResponse(
				http.StatusOK,
				`[
					{
						"id": 101,
						"body": "Normal comment"
					}
				]`,
			), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	comment, err := client.GetTriageComment(
		"123",
		"<!-- issue-triage-agent -->",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if comment != nil {
		t.Fatal("expected nil triage comment")
	}
}

func TestUpdateComment_Success(t *testing.T) {

	oldTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	http.DefaultClient.Transport = &mockRoundTripper{
		handler: func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", req.Method)
			}

			if !strings.Contains(req.URL.Path, "/issues/comments/102") {
				t.Errorf("unexpected URL: %s", req.URL.Path)
			}

			return mockResponse(http.StatusOK, `{}`), nil
		},
	}

	client := &Client{
		Token:      "test-token",
		Repository: "test-owner/test-repo",
	}

	err := client.UpdateComment(
		102,
		"Updated triage result",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}