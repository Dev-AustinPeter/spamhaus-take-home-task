package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Dev-AustinPeter/spamhaus-take-home-task/constants"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/types"
)

var (
	Mutex    sync.RWMutex
	URLStore sync.Map

	semaphore = make(chan struct{}, constants.MAX_DOWNLOADS)

	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}
)

func ParseJson(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

func LoadData(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("[INFO] No existing data file found, starting fresh.")
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var tempStore map[string]*types.URLData
	if err := decoder.Decode(&tempStore); err != nil {
		log.Println("[ERROR] Error loading data file:", err)
		return
	}
	for k, v := range tempStore {
		URLStore.Store(k, v)
	}
}

func SaveData(filePath string) {
	tempStore := make(map[string]*types.URLData)
	URLStore.Range(func(key, value interface{}) bool {
		tempStore[key.(string)] = value.(*types.URLData)
		return true
	})
	data, err := json.MarshalIndent(tempStore, "", "  ")
	if err != nil {
		log.Println("[ERROR] Error marshaling data:", err)
		return
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Println("[ERROR] Error writing data to file:", err)
	}
}

func StartBatchSave(filepath string) {
	ticker := time.NewTicker(time.Duration(constants.BATCH_SAVE_INTERVAL) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		SaveData(filepath)
		log.Println("[INFO] Data batch saved.")
	}
}

func FilterByURL(urls []*types.URLData, query string) []*types.URLData {
	var filtered []*types.URLData
	for _, url := range urls {
		if url.URL == query {
			filtered = append(filtered, url)
		}
	}
	if filtered == nil {
		return []*types.URLData{}
	}
	return filtered
}

func FetchURL(url string) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	start := time.Now()
	resp, err := httpClient.Get(url)
	if err != nil {
		if data, exists := URLStore.Load(url); exists {
			data.(*types.URLData).FailureCount++
		}
		log.Printf("[ERROR] Failed to fetch URL: %s, Error: %v\n", url, err)

		return
	}

	defer resp.Body.Close()
	elapsed := time.Since(start).Seconds()

	if data, exists := URLStore.Load(url); exists {
		data.(*types.URLData).FetchTime = elapsed
		data.(*types.URLData).SuccessCount++
		data.(*types.URLData).LastFetched = time.Now().Format(time.RFC3339)

		log.Printf("[INFO] Successfully fetched URL: %s, Fetch Time: %.2f seconds, Success Count: %d, Failure Count: %d\n", url, elapsed, data.(*types.URLData).SuccessCount, data.(*types.URLData).FailureCount)
	}
}
