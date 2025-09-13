package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// AwardsHandler serves GET /api/awards with optional filters: year, event
func AwardsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	yearParam := r.URL.Query().Get("year")
	eventParam := strings.ToLower(r.URL.Query().Get("event"))

	awards := QWERBandData.Awards
	var filtered []Award

	for _, a := range awards {
		if yearParam != "" {
			y, err := strconv.Atoi(yearParam)
			if err != nil || a.Year != y {
				continue
			}
		}
		if eventParam != "" {
			if !strings.Contains(strings.ToLower(a.Event), eventParam) {
				continue
			}
		}
		filtered = append(filtered, a)
	}

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Awards retrieved", Data: filtered})
}
