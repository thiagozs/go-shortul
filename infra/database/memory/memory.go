package memory

import (
	"encoding/json"
	"fmt"
	"sync"
)

// URLStats holds the statistics for a shortened URL
type URLStats struct {
	Count           int      `json:"count"`
	LastIPs         []string `json:"last_ips"`
	Referrers       []string `json:"referrers"`
	LastGeoLocation string   `json:"last_geo_location"`
}

// URLStore to hold the shortened URLs, their original URLs, and statistics
type URLStore struct {
	sync.RWMutex
	urls  map[string]string
	stats map[string]*URLStats
}

// NewURLStore creates a new URLStore
func NewURLStore() *URLStore {
	return &URLStore{
		urls:  make(map[string]string),
		stats: make(map[string]*URLStats),
	}
}

// Save stores a shortened URL with its original URL and initializes statistics
func (s *URLStore) Save(shortURL, originalURL string) error {
	s.Lock()
	defer s.Unlock()
	s.urls[shortURL] = originalURL
	s.stats[shortURL] = &URLStats{LastIPs: []string{}, Referrers: []string{}}
	return nil
}

// Get retrieves the original URL from a shortened URL
func (s *URLStore) Get(shortURL string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	originalURL, found := s.urls[shortURL]
	return originalURL, found
}

// GetStats retrieves the statistics for a given shortened URL
func (s *URLStore) GetStats(shortURL string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	stats, found := s.stats[shortURL]

	rr, err := json.Marshal(stats)
	if err != nil {
		return "", false
	}

	return string(rr), found
}

// UpdateURL updates the original URL for a given short URL
func (s *URLStore) UpdateURL(shortURL, newOriginalURL string) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.urls[shortURL]; !exists {
		return fmt.Errorf("short URL not found")
	}

	s.urls[shortURL] = newOriginalURL
	return nil
}

// UpdateStats updates the statistics for a given shortened URL
func (s *URLStore) UpdateStats(shortURL, ip, referrer, geoLocation string) error {
	s.Lock()
	defer s.Unlock()
	stats, found := s.stats[shortURL]
	if !found {
		return fmt.Errorf("short URL not found")
	}

	// Update count
	stats.Count++

	// Update last IPs
	if len(stats.LastIPs) >= 5 {
		stats.LastIPs = stats.LastIPs[1:] // Keep only the last 5 IPs
	}
	stats.LastIPs = append(stats.LastIPs, ip)

	// Update referrers
	if referrer != "" {
		if len(stats.Referrers) >= 5 {
			stats.Referrers = stats.Referrers[1:] // Keep only the last 5 referrers
		}
		stats.Referrers = append(stats.Referrers, referrer)
	}

	// Update last geo location
	stats.LastGeoLocation = geoLocation

	return nil
}

// Flush removes all key-value pairs from the URLStore and returns a backup JSON
func (s *URLStore) Flush() (map[string]string, error) {
	s.Lock()
	defer s.Unlock()
	backup := s.urls // Create a backup before flushing
	s.urls = make(map[string]string)
	s.stats = make(map[string]*URLStats)
	return backup, nil
}

// Backup returns a copy of all stored URLs as a JSON string
func (s *URLStore) Backup() ([]byte, error) {
	s.RLock()
	defer s.RUnlock()
	return json.Marshal(s.urls)
}

// Import loads URLs from a JSON string into the URLStore
func (s *URLStore) Import(data []byte) error {
	var importedURLs map[string]string
	if err := json.Unmarshal(data, &importedURLs); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	for shortURL, originalURL := range importedURLs {
		s.urls[shortURL] = originalURL
		s.stats[shortURL] = &URLStats{LastIPs: []string{}, Referrers: []string{}} // Initialize stats
	}
	return nil
}
