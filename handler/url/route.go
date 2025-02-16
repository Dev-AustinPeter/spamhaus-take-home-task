package url

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Dev-AustinPeter/spamhaus-take-home-task/middleware"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/types"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router, middleware *middleware.RateLimiter) {
	log.Println("[INFO] Registering URL routes...")

	router.Handle("/url", middleware.Limit(http.HandlerFunc(h.handleSubmit))).Methods("POST")
	router.Handle("/url", middleware.Limit(http.HandlerFunc(h.handleGet))).Methods("GET")
	router.Handle("/urls", middleware.Limit(http.HandlerFunc(h.handleListAll))).Methods("GET")
}

func (h *Handler) handleSubmit(w http.ResponseWriter, r *http.Request) {
	//get JSON payload
	log.Println("[INFO] Received request to submit URL...")
	var payload types.RequestUrlPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, exists := utils.URLStore.Load(payload.URL); !exists {
		utils.URLStore.Store(payload.URL, &types.URLData{URL: payload.URL, Count: 1, CreatedAt: time.Now()})
	} else {
		if data, ok := utils.URLStore.Load(payload.URL); ok {
			data.(*types.URLData).Count++
		}
	}
	utils.WriteJson(w, http.StatusAccepted, payload)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("url")

	if query == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("URL is required"))
		return
	}

	var urls []*types.URLData

	utils.URLStore.Range(func(_, value interface{}) bool {
		urls = append(urls, value.(*types.URLData))
		return true
	})

	urls = utils.FilterByURL(urls, query)

	if len(urls) == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("URL not found"))
		return
	}

	utils.FetchURL(query)

	utils.WriteJson(w, http.StatusOK, urls[0])
}
func (h *Handler) handleListAll(w http.ResponseWriter, r *http.Request) {
	sortOrder := r.URL.Query().Get("sort")
	var urls []*types.URLData

	utils.URLStore.Range(func(_, value interface{}) bool {
		urls = append(urls, value.(*types.URLData))
		return true
	})

	if sortOrder == "smallest" {
		sort.Slice(urls, func(i, j int) bool { return urls[i].Count < urls[j].Count })
	} else {
		sort.Slice(urls, func(i, j int) bool {
			return urls[i].CreatedAt.Format(time.RFC3339) > urls[j].CreatedAt.Format(time.RFC3339)
		})
	}

	if len(urls) > 50 {
		urls = urls[:50]
	}

	utils.WriteJson(w, http.StatusOK, urls)
}
