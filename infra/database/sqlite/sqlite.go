package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// URLStats holds the statistics for a shortened URL
type URLStats struct {
	Count           int      `json:"count"`
	LastIPs         []string `json:"last_ips"`
	Referrers       []string `json:"referrers"`
	LastGeoLocation string   `json:"last_geo_location"`
}

// URLStore to hold the SQLite DB connection and manage URLs and statistics
type URLStore struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewURLStore creates a new URLStore and initializes the database
func NewURLStore(dbFilePath string, logger *slog.Logger) (*URLStore, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	// Initialize the database tables
	if err := initializeDB(db); err != nil {
		return nil, err
	}

	return &URLStore{db: db, logger: logger}, nil
}

// initializeDB sets up the necessary tables
func initializeDB(db *sql.DB) error {
	createURLsTable := `
	CREATE TABLE IF NOT EXISTS urls (
		short_url TEXT PRIMARY KEY,
		original_url TEXT NOT NULL
	);`

	createStatsTable := `
	CREATE TABLE IF NOT EXISTS url_stats (
		short_url TEXT PRIMARY KEY,
		count INTEGER,
		last_ips TEXT,
		referrers TEXT,
		last_geo_location TEXT,
		FOREIGN KEY (short_url) REFERENCES urls (short_url)
	);`

	if _, err := db.Exec(createURLsTable); err != nil {
		return err
	}
	if _, err := db.Exec(createStatsTable); err != nil {
		return err
	}

	return nil
}

// Save stores a shortened URL with its original URL and initializes statistics
func (s *URLStore) Save(shortURL, originalURL string) error {
	// Insert URL into urls table
	_, err := s.db.Exec("INSERT INTO urls (short_url, original_url) VALUES (?, ?)", shortURL, originalURL)
	if err != nil {
		return err
	}

	// Initialize stats for the URL
	_, err = s.db.Exec("INSERT INTO url_stats (short_url, count, last_ips, referrers, last_geo_location) VALUES (?, 0, '', '', '')", shortURL)
	return err
}

// Get retrieves the original URL from a shortened URL
func (s *URLStore) Get(shortURL string) (string, bool) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortURL).Scan(&originalURL)
	if err != nil {
		s.logger.Error("Failed to get original URL", slog.String("error", err.Error()))
		return "", false
	}
	return originalURL, true
}

// GetStats retrieves the statistics for a given shortened URL
func (s *URLStore) GetStats(shortURL string) (string, bool) {
	var stats URLStats
	var lastIPs, referrers string
	err := s.db.QueryRow("SELECT count, last_ips, referrers, last_geo_location FROM url_stats WHERE short_url = ?", shortURL).Scan(&stats.Count, &lastIPs, &referrers, &stats.LastGeoLocation)
	if err != nil {
		return "", false
	}

	stats.LastIPs = splitString(lastIPs)
	stats.Referrers = splitString(referrers)

	rr, err := json.Marshal(stats)
	if err != nil {
		return "", false
	}

	return string(rr), true
}

// UpdateURL updates the original URL for a given short URL
func (s *URLStore) UpdateURL(shortURL, newOriginalURL string) error {
	_, err := s.db.Exec("UPDATE urls SET original_url = ? WHERE short_url = ?", newOriginalURL, shortURL)
	return err
}

// UpdateStats updates the statistics for a given shortened URL
func (s *URLStore) UpdateStats(shortURL, ip, referrer, geoLocation string) error {
	// Retrieve existing stats
	statsStr, found := s.GetStats(shortURL)
	if !found {
		return fmt.Errorf("short URL not found")
	}

	var stats URLStats
	// Unmarshal the existing stats
	if err := json.Unmarshal([]byte(statsStr), &stats); err != nil {
		return err
	}

	// Update stats
	stats.Count++
	stats.LastIPs = appendWithLimit(stats.LastIPs, ip, 5)
	stats.Referrers = appendWithLimit(stats.Referrers, referrer, 5)
	stats.LastGeoLocation = geoLocation

	_, err := s.db.Exec("UPDATE url_stats SET count = ?, last_ips = ?, referrers = ?, last_geo_location = ? WHERE short_url = ?",
		stats.Count, joinStrings(stats.LastIPs), joinStrings(stats.Referrers), stats.LastGeoLocation, shortURL)
	return err
}

// Helper function to append an item to a list with a maximum limit
func appendWithLimit(list []string, item string, limit int) []string {
	if item == "" {
		return list
	}
	if len(list) >= limit {
		list = list[1:]
	}
	return append(list, item)
}

// Helper function to join a list of strings into a JSON-encoded string
func joinStrings(list []string) string {
	return strings.Join(list, ",")
}

// Helper function to split a JSON-encoded string into a list of strings
func splitString(data string) []string {
	return strings.Split(data, ",")
}

// Flush removes all key-value pairs from the URLStore and returns a backup JSON
func (s *URLStore) Flush() (map[string]string, error) {
	backup := make(map[string]string)
	rows, err := s.db.Query("SELECT short_url, original_url FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shortURL, originalURL string
		if err := rows.Scan(&shortURL, &originalURL); err != nil {
			return nil, err
		}
		backup[shortURL] = originalURL
	}

	// Delete all entries from urls and url_stats
	_, err = s.db.Exec("DELETE FROM urls")
	if err != nil {
		return nil, err
	}
	_, err = s.db.Exec("DELETE FROM url_stats")
	return backup, err
}

// Backup returns a copy of all stored URLs as a JSON string
func (s *URLStore) Backup() ([]byte, error) {
	rows, err := s.db.Query("SELECT short_url, original_url FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	backup := make(map[string]string)
	for rows.Next() {
		var shortURL, originalURL string
		if err := rows.Scan(&shortURL, &originalURL); err != nil {
			return nil, err
		}
		backup[shortURL] = originalURL
	}

	return json.Marshal(backup)
}

// Import loads URLs from a JSON string into the URLStore
func (s *URLStore) Import(data []byte) error {
	var importedURLs map[string]string
	if err := json.Unmarshal(data, &importedURLs); err != nil {
		return err
	}

	for shortURL, originalURL := range importedURLs {
		if err := s.Save(shortURL, originalURL); err != nil {
			log.Printf("Failed to import URL %s: %v", shortURL, err)
		}
	}
	return nil
}
