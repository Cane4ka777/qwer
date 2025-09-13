package server

import (
	"encoding/json"
	"net/http"
	"time"
)

// Handler for GET /api - API documentation and available endpoints
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		response := APIResponse{
			Success: false,
			Message: "Method not allowed",
			Data:    nil,
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	apiInfo := map[string]interface{}{
		"name":        "QWER Band API",
		"version":     "1.0.0",
		"description": "REST API for QWER band information, members, songs, albums, and awards",
		"band":        "QWER (쿼터)",
		"company":     "Tamago Production",
		"endpoints": map[string]interface{}{
			"GET /api": map[string]interface{}{
				"description": "Get API information and available endpoints",
			},
			"GET /docs": map[string]interface{}{
				"description": "OpenAPI (Swagger) spec rendered via ReDoc",
			},
			"GET /api/band": map[string]interface{}{
				"description": "Get complete band information",
			},
			"GET /api/members": map[string]interface{}{
				"description": "Get all members or specific member",
				"parameters": map[string]string{
					"id":     "Member ID (optional)",
					"name":   "Member name or stage name (optional)",
					"search": "Search in name/stage/position (optional)",
					"sort":   "Sort by name|stage_name (optional)",
					"page":   "Page number (optional)",
					"limit":  "Page size (optional)",
				},
			},
			"GET /api/songs": map[string]interface{}{
				"description": "Get all songs or specific song",
				"parameters": map[string]string{
					"id":     "Song ID (optional)",
					"title":  "Song title (optional)",
					"album":  "Album name (optional)",
					"search": "Search in title/album/genre (optional)",
					"sort":   "Sort by title|album|date (optional)",
					"page":   "Page number (optional)",
					"limit":  "Page size (optional)",
				},
			},
			"GET /api/albums": map[string]interface{}{
				"description": "Get all albums or specific album",
				"parameters": map[string]string{
					"id":     "Album ID (optional)",
					"title":  "Album title (optional)",
					"search": "Search in title/type (optional)",
					"sort":   "Sort by title|date (optional)",
					"page":   "Page number (optional)",
					"limit":  "Page size (optional)",
				},
			},
			"GET /api/awards": map[string]interface{}{
				"description": "Get awards with optional filters",
				"parameters": map[string]string{
					"year":  "Year (optional)",
					"event": "Event name (optional)",
				},
			},
		},
		"examples": map[string]interface{}{
			"Get band info":      "/api/band",
			"Get all members":    "/api/members",
			"Get member by name": "/api/members?name=chodan",
			"Get member by ID":   "/api/members?id=1",
			"Get all songs":      "/api/songs",
			"Search songs":       "/api/songs?search=discord&sort=title&page=1&limit=10",
			"Get songs by album": "/api/songs?album=manito",
			"Get all albums":     "/api/albums",
			"Search albums":      "/api/albums?search=blossom&sort=date",
			"Get awards":         "/api/awards?year=2024",
		},
	}

	resp := APIResponse{
		Success: true,
		Message: "QWER Band API is running successfully",
		Data:    apiInfo,
	}

	payload, _ := json.Marshal(resp)
	etag := computeETag(payload)
	last := time.Now().Add(-time.Minute) // stable-ish timestamp for cache
	if handleConditional(w, r, etag, last) {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
