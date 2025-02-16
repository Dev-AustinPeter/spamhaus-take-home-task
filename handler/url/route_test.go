package url

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Dev-AustinPeter/spamhaus-take-home-task/types"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/utils"
)

func TestHandleSubmit(t *testing.T) {
	handler := &Handler{}

	// Mock request payload
	payload := types.RequestUrlPayload{
		URL: "http://example.com",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Create a request
	req := httptest.NewRequest("POST", "/url", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Call the handler
	handler.handleSubmit(w, req)

	// Check response status
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		t.Errorf("Expected status %d but got %d", http.StatusAccepted, res.StatusCode)
	}

	// Verify that the URL was stored
	if data, exists := utils.URLStore.Load(payload.URL); exists {
		urlData := data.(*types.URLData)
		if urlData.Count != 1 {
			t.Errorf("Expected count to be 1 but got %d", urlData.Count)
		}
		if urlData.URL != payload.URL {
			t.Errorf("Expected URL to be %s but got %s", payload.URL, urlData.URL)
		}
	} else {
		t.Errorf("Expected URL to be stored but it was not found")
	}
}

func TestHandleGet_Success(t *testing.T) {
	handler := &Handler{}

	// Mock stored URL
	testURL := "http://example.com"
	utils.URLStore.Store(testURL, &types.URLData{
		URL:         testURL,
		Count:       2,
		LastFetched: "2024-01-01T00:00:00Z",
	})

	// Create a request with query param
	req := httptest.NewRequest("GET", "/url?url="+testURL, nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.handleGet(w, req)

	// Check response status
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, res.StatusCode)
	}
}

func TestHandleGet_MissingQuery(t *testing.T) {
	handler := &Handler{}

	// Create a request without query param
	req := httptest.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.handleGet(w, req)

	// Check response status
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d but got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestHandleGet_NotFound(t *testing.T) {
	handler := &Handler{}

	// Create a request with a non-existent URL
	req := httptest.NewRequest("GET", "/url?url=http://notfound.com", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.handleGet(w, req)

	// Check response status
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status %d but got %d", http.StatusNotFound, res.StatusCode)
	}
}

func TestHandleListAll_DefaultSort(t *testing.T) {
	handler := &Handler{}

	// Mock stored URLs
	utils.URLStore.Store("http://example1.com", &types.URLData{
		URL:       "http://example1.com",
		Count:     5,
		CreatedAt: time.Now().Add(-10 * time.Minute),
	})

	utils.URLStore.Store("http://example2.com", &types.URLData{
		URL:       "http://example2.com",
		Count:     3,
		CreatedAt: time.Now(),
	})

	// Create a request
	req := httptest.NewRequest("GET", "/urls", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.handleListAll(w, req)

	// Check response status
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, res.StatusCode)
	}

	// Decode response
	var urls []types.URLData
	json.NewDecoder(res.Body).Decode(&urls)

	// Verify the sorting order (newest first)
	if len(urls) < 2 || urls[0].URL != "http://example2.com" {
		t.Errorf("Expected latest URL to be first, but got %v", urls[0].URL)
	}
}

func TestHandleListAll_SmallestSort(t *testing.T) {
	handler := &Handler{}

	// Mock stored URLs
	utils.URLStore.Store("http://example1.com", &types.URLData{
		URL:       "http://example1.com",
		Count:     10,
		CreatedAt: time.Now().Add(-10 * time.Minute),
	})

	utils.URLStore.Store("http://example2.com", &types.URLData{
		URL:       "http://example2.com",
		Count:     1,
		CreatedAt: time.Now(),
	})

	// Create a request with sorting by smallest count
	req := httptest.NewRequest("GET", "/urls?sort=smallest", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.handleListAll(w, req)

	// Check response status
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, res.StatusCode)
	}

	// Decode response
	var urls []types.URLData
	json.NewDecoder(res.Body).Decode(&urls)

	// Verify sorting order (smallest count first)
	if len(urls) < 2 || urls[0].URL != "http://example2.com" {
		t.Errorf("Expected smallest count URL first, but got %v", urls[0].URL)
	}
}

func TestHandleListAll_LimitTo50(t *testing.T) {
	handler := &Handler{}

	// Mock stored URLs (more than 50)
	for i := 1; i <= 60; i++ {
		utils.URLStore.Store(
			"http://example"+fmt.Sprintf("%d", i)+".com",
			&types.URLData{
				URL:       "http://example" + fmt.Sprintf("%d", i) + ".com",
				Count:     i,
				CreatedAt: time.Now(),
			},
		)
	}

	// Create a request
	req := httptest.NewRequest("GET", "/urls", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.handleListAll(w, req)

	// Check response status
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, res.StatusCode)
	}

	// Decode response
	var urls []types.URLData
	json.NewDecoder(res.Body).Decode(&urls)

	// Verify limit is applied
	if len(urls) > 50 {
		t.Errorf("Expected at most 50 URLs but got %d", len(urls))
	}
}
