package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/thiagozs/go-shorturl/config"
	"github.com/thiagozs/go-shorturl/pkg/utils"
)

// Handler holds the logger and URLStore
type Handler struct {
	params *HandlerParams
}

func NewHandler(opts ...Options) (*Handler, error) {
	params, err := newHandlerParams(opts...)
	if err != nil {
		return nil, err
	}

	if params.Logger() == nil {
		return nil, fmt.Errorf("logger is required")
	} else if params.Store() == nil {
		return nil, fmt.Errorf("store is required")
	}

	return &Handler{params: params}, nil
}

// shortenHandler handles requests to shorten a URL
func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Generate a short URL and store it
	shortURL, err := utils.GenerateShortURL()
	if err != nil {
		h.params.Logger().Error("Failed to generate short URL", slog.String("error", err.Error()))
		http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
		return
	}

	h.params.Store().Save(shortURL, originalURL)

	// Respond with the short URL in JSON format
	response := map[string]string{"short_url": fmt.Sprintf("http://localhost:%s/%s", h.params.Port(), shortURL)}
	if !h.params.Local() {
		if h.params.HTTPS() {
			response = map[string]string{"short_url": fmt.Sprintf("https://%s/%s", h.params.Domain(), shortURL)}
		} else {
			response = map[string]string{"short_url": fmt.Sprintf("http://%s/%s", h.params.Domain(), shortURL)}
		}
	}

	h.params.Logger().Info("URL shortened", slog.String("original_url", originalURL), slog.String("short_url", shortURL))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// redirectHandler handles requests to redirect from a short URL to the original URL
func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:] // Get the short URL from the path
	originalURL, found := h.params.Store().Get(shortURL)

	if !found {
		h.params.Logger().Warn("Short URL not found", slog.String("short_url", shortURL))
		http.NotFound(w, r)
		return
	}

	// Get client IP and referrer for stats
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	referrer := r.Referer()
	geoLocation := utils.GetGeoLocation(ip) // Fetch geolocation using the helper function

	// Update statistics
	if err := h.params.Store().UpdateStats(shortURL, ip,
		referrer, geoLocation); err != nil {
		h.params.Logger().Error("Failed to update stats", slog.String("short_url", shortURL), slog.String("error", err.Error()))
		http.Error(w, "Failed to update stats", http.StatusInternalServerError)
		return
	}

	h.params.Logger().Info("Redirecting", slog.String("short_url", shortURL), slog.String("original_url", originalURL), slog.String("geo_location", geoLocation))
	http.Redirect(w, r, originalURL, http.StatusFound) // Redirect to the original URL
}

// statsHandler handles requests to retrieve statistics for a shortened URL
func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("short_url")
	statsStr, found := h.params.Store().GetStats(shortURL)

	if !found {
		h.params.Logger().Warn("Short URL not found for stats",
			slog.String("short_url", shortURL))
		http.NotFound(w, r)
		return
	}

	var stats map[string]interface{}
	err := json.Unmarshal([]byte(statsStr), &stats)
	if err != nil {
		h.params.Logger().Error("Failed to marshal stats", slog.String("short_url", shortURL), slog.String("error", err.Error()))
		http.Error(w, "Failed to marshal stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.params.Logger().Info("Retrieved stats for URL", slog.String("short_url", shortURL))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// updateHandler handles requests to update the original URL for a given short URL
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("short_url")
	newOriginalURL := r.URL.Query().Get("new_url")

	// Validate inputs
	if shortURL == "" || newOriginalURL == "" {
		http.Error(w, "Both short_url and new_url parameters are required", http.StatusBadRequest)
		return
	}

	// Update the URL in the store
	if err := h.params.Store().UpdateURL(shortURL, newOriginalURL); err != nil {
		h.params.Logger().Warn("Failed to update URL", slog.String("short_url", shortURL), slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.params.Logger().Info("Updated URL", slog.String("short_url", shortURL), slog.String("new_original_url", newOriginalURL))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"url": shortURL, "message": "URL updated successfully"})

}

// flushHandler handles requests to flush all key-value pairs from memory and returns a JSON backup
func (h *Handler) FlushHandler(w http.ResponseWriter, r *http.Request) {
	backup, err := h.params.Store().Flush()
	if err != nil {
		h.params.Logger().Error("Failed to flush URLs", slog.String("error", err.Error()))
		http.Error(w, "Failed to flush URLs", http.StatusInternalServerError)
		return
	}

	// Return the backup as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	h.params.Logger().Info("Flushed all URLs and returned backup")
	json.NewEncoder(w).Encode(backup)
}

// backupHandler returns the current URL mappings as a JSON object
func (h *Handler) BackupHandler(w http.ResponseWriter, r *http.Request) {
	backup, err := h.params.Store().Backup()
	if err != nil {
		h.params.Logger().Error("Failed to generate backup", slog.String("error", err.Error()))
		http.Error(w, "Failed to generate backup", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(backup)
	w.WriteHeader(http.StatusOK)
	h.params.Logger().Info("Returned current backup")
}

// importHandler handles requests to import URLs from a JSON object
func (h *Handler) ImportHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.params.Logger().Error("Failed to read request body", slog.String("error", err.Error()))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := h.params.Store().Import(body); err != nil {
		h.params.Logger().Error("Failed to import URLs", slog.String("error", err.Error()))
		http.Error(w, "Failed to import URLs", http.StatusBadRequest)
		return
	}

	h.params.Logger().Info("URLs imported successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"message": "URLs imported successfully"})
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) SetConfig(cfg *config.Config) {
	h.params.SetConfig(cfg)
	h.params.SetHost(cfg.GetHost())
	h.params.SetLocal(cfg.GetLocal())
	h.params.SetDomain(cfg.GetDomain())
	h.params.SetHTTPS(cfg.GetHTTPS())
	h.params.SetPort(cfg.GetPort())
}
