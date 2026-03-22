package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHTML(t *testing.T) {
	cases := []struct {
		name           string
		status         int
		contentType    string
		body           string
		expectedError  string
		expectedResult string
	}{
		{
			name:           "success",
			status:         200,
			contentType:    "text/html",
			body:           "<html><body>Hello World</body></html>",
			expectedResult: "<html><body>Hello World</body></html>",
		},
		{
			name:           "success with charset",
			status:         200,
			contentType:    "text/html; charset=utf-8",
			body:           "<html><body>Hello World</body></html>",
			expectedResult: "<html><body>Hello World</body></html>",
		},
		{
			name:          "error 404",
			status:        404,
			contentType:   "text/html",
			body:          "Not Found",
			expectedError: "bad status code: 404",
		},
		{
			name:          "error 500",
			status:        500,
			contentType:   "text/html",
			body:          "Internal Server Error",
			expectedError: "bad status code: 500",
		},
		{
			name:          "non-html content type",
			status:        200,
			contentType:   "application/json",
			body:          `{"key": "value"}`,
			expectedError: "non-html response: application/json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()

			actual, err := getHTML(server.URL)

			if tc.expectedError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, but got nil", tc.expectedError)
				}
				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("expected error containing %q, but got %q", tc.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tc.expectedResult {
				t.Errorf("expected %q, got %q", tc.expectedResult, actual)
			}
		})
	}

	t.Run("invalid URL", func(t *testing.T) {
		_, err := getHTML(" :invalid")
		if err == nil {
			t.Fatal("expected error for invalid URL, but got nil")
		}
	})
}
