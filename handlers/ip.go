package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type IPResponse struct {
	IP        string  `json:"ip"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
}

func GetIP(w http.ResponseWriter, r *http.Request) {
	// Test mode (uncomment if needed)
	/*testMode := os.Getenv("GLOBE_TEST_MODE") == "true"
	if testMode {
		response := IPResponse{
			IP:       "127.0.0.1",
			Latitude: 40.7128,
			Longitude: -74.0060,
			City:     "New York",
			Country:  "US",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}*/


	ip := GetClientIP(r)
	log.Printf("Detected IP: %s", ip)
	log.Printf("Headers - X-Forwarded-For: %s, X-Real-IP: %s",
		r.Header.Get("X-Forwarded-For"),
		r.Header.Get("X-Real-IP"))

	lat, lon, err := GetGeolocation(ip)
	if err != nil {
		log.Printf("Error getting geolocation for IP %s: %v", ip, err)
		// Return IP even if geolocation fails
		response := IPResponse{
			IP:        ip,
			Latitude:  0.0,
			Longitude: 0.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := IPResponse{
		IP:        ip,
		Latitude:  lat,
		Longitude: lon,
	}

	log.Printf("Success - IP: %s, LAT: %f, LON: %f", response.IP, response.Latitude, response.Longitude)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetClientIP(r *http.Request) string {

	forwardedHeader := r.Header.Get("X-Forwarded-For")
	if forwardedHeader != "" {
		ips := strings.Split(forwardedHeader, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return strings.TrimSpace(realIP)
	}

	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

func GetGeolocation(ip string) (float64, float64, error) {

	if ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") {
		return 0.0, 0.0, fmt.Errorf("private IP address")
	}

	lat, lon, err := getGeolocationFromIPInfo(ip)
	if err == nil {
		return lat, lon, nil
	}

	log.Printf("ipinfo.io failed, trying ip-api.com: %v", err)

	lat, lon, err = getGeolocationFromIPAPI(ip)
	if err == nil {
		return lat, lon, nil
	}

	return 0.0, 0.0, fmt.Errorf("all geolocation services failed")
}

func getGeolocationFromIPInfo(ip string) (float64, float64, error) {
	response, err := http.Get("https://ipinfo.io/" + ip + "/json")
	if err != nil {
		return 0.0, 0.0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0.0, 0.0, fmt.Errorf("ipinfo.io returned status %d", response.StatusCode)
	}

	var data struct {
		Loc     string `json:"loc"`
		City    string `json:"city"`
		Country string `json:"country"`
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0.0, 0.0, err
	}

	if data.Loc == "" {
		return 0.0, 0.0, fmt.Errorf("no location data from ipinfo.io")
	}

	// Parse "lat,lon" string
	coords := strings.Split(data.Loc, ",")
	if len(coords) != 2 {
		return 0.0, 0.0, fmt.Errorf("invalid location format from ipinfo.io")
	}

	var lat, lon float64
	if _, err := fmt.Sscanf(coords[0], "%f", &lat); err != nil {
		return 0.0, 0.0, err
	}
	if _, err := fmt.Sscanf(coords[1], "%f", &lon); err != nil {
		return 0.0, 0.0, err
	}

	return lat, lon, nil
}

func getGeolocationFromIPAPI(ip string) (float64, float64, error) {

	response, err := http.Get("http://ip-api.com/json/" + ip + "?fields=status,lat,lon")
	if err != nil {
		return 0.0, 0.0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0.0, 0.0, fmt.Errorf("ip-api.com returned status %d", response.StatusCode)
	}

	var data struct {
		Status string  `json:"status"`
		Lat    float64 `json:"lat"`
		Lon    float64 `json:"lon"`
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0.0, 0.0, err
	}

	if data.Status != "success" {
		return 0.0, 0.0, fmt.Errorf("ip-api.com returned status: %s", data.Status)
	}

	if data.Lat == 0 && data.Lon == 0 {
		return 0.0, 0.0, fmt.Errorf("no location data from ip-api.com")
	}

	return data.Lat, data.Lon, nil
}
