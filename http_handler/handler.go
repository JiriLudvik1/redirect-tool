package http_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"redirect-tool/redis_service"
)

type Handler struct {
	RedisService *redis_service.RedisService
}

type UrlRequest struct {
	URL string `json:"url"`
}

func NewHandler(redisService *redis_service.RedisService) *Handler {
	return &Handler{
		RedisService: redisService,
	}
}

func (h *Handler) ShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)

	var requestData UrlRequest
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redirectHash, err := h.RedisService.CreateRedirectEntry(requestData.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(redirectHash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlHash := r.URL.Path[1:]
	originalUrl, err := h.RedisService.GetOriginalUrl(urlHash)
	if err != nil || originalUrl == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusFound)
}
