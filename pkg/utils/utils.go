package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// generateShortURL creates a random string to use as a short URL
func GenerateShortURL() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:6], nil // Return a 6-character string
}

// getGeoLocation fetches the geolocation of an IP address using ip-api.com
func GetGeoLocation(ip string) string {
	value := "Unknown"
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return value
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return value
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return value
	}

	// Extract city and country
	if result["status"] == "success" {
		city := result["city"].(string)
		country := result["country"].(string)
		return fmt.Sprintf("%s, %s", city, country)
	}

	return value
}
