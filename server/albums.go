package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handler for GET /api/albums - Get all albums or specific album
func AlbumsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if requesting specific album by ID
	idParam := r.URL.Query().Get("id")
	titleParam := r.URL.Query().Get("title")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := APIResponse{Success: false, Message: "Invalid album ID"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		for _, album := range QWERBandData.Discography {
			if album.ID == id {
				json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Album retrieved", Data: album})
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Album not found"})
		return
	}

	if titleParam != "" {
		titleParam = strings.ToLower(titleParam)
		for _, album := range QWERBandData.Discography {
			if strings.ToLower(album.Title) == titleParam {
				json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Album retrieved", Data: album})
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Album not found"})
		return
	}

	// Enhanced list mode
	search := strings.ToLower(r.URL.Query().Get("search"))
	sortKey := strings.ToLower(r.URL.Query().Get("sort")) // title|date
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	list := QWERBandData.Discography
	if search != "" {
		var filtered []Album
		for _, a := range list {
			if strings.Contains(strings.ToLower(a.Title), search) ||
				strings.Contains(strings.ToLower(a.Type), search) {
				filtered = append(filtered, a)
			}
		}
		list = filtered
	}

	if sortKey == "title" || sortKey == "date" {
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				var ai, aj string
				if sortKey == "title" {
					ai, aj = strings.ToLower(list[i].Title), strings.ToLower(list[j].Title)
				} else {
					ai, aj = list[i].ReleaseDate, list[j].ReleaseDate
				}
				if ai > aj {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
	}

	page := 1
	limit := 50
	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 200 {
		limit = l
	}
	start := (page - 1) * limit
	if start > len(list) {
		start = len(list)
	}
	end := start + limit
	if end > len(list) {
		end = len(list)
	}
	paged := list[start:end]

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Albums retrieved", Data: paged})
}
