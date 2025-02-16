# spamhaus-take-home-task

## Overview
This project is a Golang-based HTTP daemon that provides an API for submitting and retrieving URLs. It also includes a background process that fetches the most requested URLs periodically, while ensuring concurrency limits on downloads.

## Features
- **Submit URLs:** `POST /url` to submit URLs for processing.
- **Retrieve URL:** `GET /url` to list stored URL.
- **Retrieve latest 50 URLs:** `GET /urls` (sorted by request count or timestamp).
- **Background Fetching:** The top 10 most requested URLs are fetched every 60 seconds.
- **Concurrency Control:** No more than 3 URLs are downloaded in parallel.
- **Graceful Shutdown:** Ensures data persistence on shutdown.

## Project Structure
```
/spamhaus-take-home-task
│── /cmd/main.go        # Main entry point
│── /handlers           # HTTP route handlers
│── /service            # background jobs
│── /utils              # Utilities
│── /middleware         # Middleware (rate limiting)
│── /types              # Data models
│── /constants          # Constant values
│── /config             # Environment values
│── go.mod              # Go module dependencies
│── Dockerfile          # Containerization support
│── data.json           # Persistent storage for URLs
```

## Installation
### Prerequisites
- **Go 1.23+**
- **Docker (optional, for containerization)**

### Setup
Clone the repository:
```sh
git clone https://github.com/Dev-AustinPeter/spamhaus-take-home-task.git
cd spamhaus-take-home-task
```
Install dependencies:
```sh
go mod tidy
```
Run the application:
```sh
go run cmd/api/main.go
```

## API Endpoints
### **Submit a URL**
- **Endpoint:** `POST /url`
- **Request Body:**
```json
{
  "url": "http://example.com"
}
```
- **Response:** `202 Accepted`

### **Retrieve stored URL**
- **Endpoint:** `GET /url`
- **Query Params:**
  - `url=http://example.com` → Query param required.
- **Response:** JSON list of stored URL with submission counts.

### **Retrieve latest 50 URLs**
- **Endpoint:** `GET /urls`
- **Query Params:**
  - `sort=smallest` → Sort by submission count (default: sorted by timestamp)
- **Response:** JSON list of URLs sorted accordingly.

## Background Process
- Runs every **60 seconds**.
- Fetches the **top 10 most submitted URLs**.
- Limits **concurrent downloads to 3**.
- Logs download time, success and failures.

## Running with Docker
To run the application inside a Docker container:
```sh
docker build -t spamhaus-take-home-task .
docker run -p 8080:8080 spamhaus-take-home-task
```

## Graceful Shutdown
The server listens for termination signals (`SIGINT`, `SIGTERM`) and ensures data is saved before exiting.

## Running Tests
```sh
make test
or
go test -v ./...
```

## API Docs
[openapi](https://github.com/Dev-AustinPeter/spamhaus-take-home-task/blob/main/docs/openapi.yaml)