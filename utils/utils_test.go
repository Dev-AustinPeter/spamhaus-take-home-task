package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Dev-AustinPeter/spamhaus-take-home-task/types"
	"github.com/stretchr/testify/assert"
)

func TestParseJson(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
		want    types.RequestUrlPayload
	}{
		{
			name:    "Valid JSON",
			body:    `{"url": "http://example.com"}`,
			wantErr: false,
			want:    types.RequestUrlPayload{URL: "http://example.com"},
		},
		{
			name:    "Invalid JSON",
			body:    `{"url": "http://example.com",}`,
			wantErr: true,
		},
		{
			name:    "Empty Body",
			body:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://example.com", io.NopCloser(bytes.NewReader([]byte(tt.body))))
			var payload types.RequestUrlPayload

			err := ParseJson(req, &payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJson() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && payload != tt.want {
				t.Errorf("ParseJson() got = %+v, want %+v", payload, tt.want)
			}
		})
	}
}

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		input      any
		expectBody string
		expectErr  bool
	}{
		{
			name:       "Valid JSON response",
			status:     http.StatusOK,
			input:      map[string]string{"message": "success"},
			expectBody: `{"message":"success"}`,
			expectErr:  false,
		},
		{
			name:       "Invalid JSON response",
			status:     http.StatusInternalServerError,
			input:      make(chan int), // Invalid JSON type
			expectBody: "",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			err := WriteJson(recorder, tt.status, tt.input)

			if (err != nil) != tt.expectErr {
				t.Errorf("WriteJson() error = %v, wantErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				if recorder.Code != tt.status {
					t.Errorf("WriteJson() status code = %d, want %d", recorder.Code, tt.status)
				}

				body := recorder.Body.String()
				if strings.TrimSpace(body) != strings.TrimSpace(tt.expectBody) {
					t.Errorf("WriteJson() body = %q, want %q", body, tt.expectBody)
				}
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		err        error
		expectBody string
	}{
		{
			name:       "Standard error",
			status:     http.StatusBadRequest,
			err:        errors.New("bad request error"),
			expectBody: `{"error":"bad request error"}`,
		},
		{
			name:       "Internal server error",
			status:     http.StatusInternalServerError,
			err:        errors.New("internal server error"),
			expectBody: `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			WriteError(recorder, tt.status, tt.err)

			if recorder.Code != tt.status {
				t.Errorf("WriteError() status code = %d, want %d", recorder.Code, tt.status)
			}

			body := recorder.Body.String()
			if strings.TrimSpace(body) != strings.TrimSpace(tt.expectBody) {
				t.Errorf("WriteError() body = %q, want %q", body, tt.expectBody)
			}
		})
	}
}

func TestLoadData(t *testing.T) {
	// Create a temporary file to simulate an existing data file
	tempFile, err := os.CreateTemp("", "test_data.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// mock data
	data := map[string]*types.URLData{
		"key1": {URL: "http://example.com"},
		"key2": {URL: "http://example.org"},
	}
	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)

	_, err = tempFile.Write(jsonData)
	assert.NoError(t, err)
	tempFile.Close()

	LoadData(tempFile.Name())

	// Validate stored data
	storedData, exists := URLStore.Load("key1")
	assert.True(t, exists)
	assert.Equal(t, "http://example.com", storedData.(*types.URLData).URL)

	storedData, exists = URLStore.Load("key2")
	assert.True(t, exists)
	assert.Equal(t, "http://example.org", storedData.(*types.URLData).URL)
}

func TestSaveData(t *testing.T) {
	// Create a temporary file to save test data
	tempFile, err := os.CreateTemp("", "test_save_data.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Mock URLStore with test data
	URLStore.Store("key1", &types.URLData{URL: "http://example.com"})
	URLStore.Store("key2", &types.URLData{URL: "http://example.org"})

	SaveData(tempFile.Name())

	savedData, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)

	var loadedData map[string]*types.URLData
	err = json.Unmarshal(savedData, &loadedData)
	assert.NoError(t, err)

	assert.Contains(t, loadedData, "key1")
	assert.Equal(t, "http://example.com", loadedData["key1"].URL)

	assert.Contains(t, loadedData, "key2")
	assert.Equal(t, "http://example.org", loadedData["key2"].URL)
}

func TestFilterByURL(t *testing.T) {
	// mock data
	urls := []*types.URLData{
		{URL: "http://example.com"},
		{URL: "http://example.org"},
	}

	tests := []struct {
		name     string
		query    string
		expected []*types.URLData
	}{
		{
			name:  "Filter existing URL",
			query: "http://example.com",
			expected: []*types.URLData{
				{URL: "http://example.com"},
			},
		},
		{
			name:  "Filter another existing URL",
			query: "http://example.org",
			expected: []*types.URLData{
				{URL: "http://example.org"},
			},
		},
		{
			name:     "Filter non-existing URL",
			query:    "http://notfound.com",
			expected: []*types.URLData{}, // Expecting an empty slice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterByURL(urls, tt.query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFetchURL(t *testing.T) {

	URLStore.Store("http://example.com", &types.URLData{URL: "http://example.com"})

	FetchURL("http://example.com")

	storedData, _ := URLStore.Load("http://example.com")
	urlData := storedData.(*types.URLData)
	assert.Equal(t, 1, urlData.SuccessCount)

}
