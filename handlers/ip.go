package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type IPResponse struct {
	IP string `json:"ip"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	City string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
}

func GetIP(w http.ResponseWriter, r *http.Request) {

	//test mode
	testMode := os.Getenv("GLOBE_TEST_MODE") == "true"

	if testMode {
		response := IPResponse {
			IP: "127.0.0.1",
			Latitude: 40.7128,
			Longitude: -74.0060,
			City: "New York",
			Country: "US",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	//normal implementation
	ip := GetClientIP(r)

	lat, lon, err := GetGeolocation(ip)

	if err != nil {
		log.Printf("Error getting geolocation: %v", err)
		lat, lon = 0.0, 0.0
	}
	
	response := IPResponse{
		IP: ip,
		Latitude: lat,
		Longitude: lon,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//check later - need better way to get client ip
func GetClientIP(r *http.Request) string {

	forwardedHeader := r.Header.Get("X-Forwarded-For")
	if forwardedHeader != "" {
		ips := strings.Split(forwardedHeader, ",")
		if len(ips) > 0 {
			return ips[0]
		}
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

func GetGeolocation(ip string) (float64, float64, error) {
	response, err := http.Get("https://ipinfo.io/" + ip + "/json")
	if err != nil {
		return 0.0, 0.0, err
	}
	defer response.Body.Close()

	var data struct {
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		City string `json:"city"`
		Country string `json:"country_name"`
	}
	
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0.0, 0.0, err
	}

	return data.Latitude, data.Longitude, nil
	
}