package types

import "time"

type RequestUrlPayload struct {
	URL string `json:"url"`
}

type URLData struct {
	URL          string    `json:"url"`
	Count        int       `json:"count"`
	LastFetched  string    `json:"last_fetched"`
	FetchTime    float64   `json:"fetch_time"`
	SuccessCount int       `json:"success_count"`
	FailureCount int       `json:"failure_count"`
	CreatedAt    time.Time `json:"created_at"`
}
