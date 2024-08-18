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

	requestData, err := parseRequestData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, validationErr := isValidUrl(requestData.URL)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(validationErr.Error()))
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

func parseRequestData(r *http.Request) (*UrlRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &requestData, nil
}
