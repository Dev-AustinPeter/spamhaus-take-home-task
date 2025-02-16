package constants

const (
	DATA_FILE           = "data.json"
	RATE_LIMIT          = 5   // Maximum requests per IP per minute
	MAX_DOWNLOADS       = 3   // Max concurrent downloads
	FETCH_INTERVAL      = 60  // Seconds between background fetch runs
	BATCH_SAVE_INTERVAL = 300 // Save data every 5 minutes
)
