package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handler for GET /api/songs - Get all songs or specific song
func SongsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Collect all songs from all albums
	var allSongs []Song
	for _, album := range QWERBandData.Discography {
		allSongs = append(allSongs, album.Songs...)
	}

	// Check if requesting specific song by ID
	idParam := r.URL.Query().Get("id")
	titleParam := r.URL.Query().Get("title")
	albumParam := r.URL.Query().Get("album")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := APIResponse{
				Success: false,
				Message: "Invalid song ID",
				Data:    nil,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		for _, song := range allSongs {
			if song.ID == id {
				response := APIResponse{
					Success: true,
					Message: "Song information retrieved successfully",
					Data:    song,
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		response := APIResponse{
			Success: false,
			Message: "Song not found",
			Data:    nil,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if titleParam != "" {
		titleParam = strings.ToLower(titleParam)
		for _, song := range allSongs {
			if strings.ToLower(song.Title) == titleParam {
				response := APIResponse{
					Success: true,
					Message: "Song information retrieved successfully",
					Data:    song,
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		response := APIResponse{
			Success: false,
			Message: "Song not found",
			Data:    nil,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filter by album if specified
	if albumParam != "" {
		albumParam = strings.ToLower(albumParam)
		var albumSongs []Song
		for _, song := range allSongs {
			if strings.ToLower(song.Album) == albumParam {
				albumSongs = append(albumSongs, song)
			}
		}

		if len(albumSongs) == 0 {
			response := APIResponse{
				Success: false,
				Message: "No songs found for the specified album",
				Data:    nil,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := APIResponse{
			Success: true,
			Message: "Songs from album retrieved successfully",
			Data:    albumSongs,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Enhanced list: search/sort/pagination
	search := strings.ToLower(r.URL.Query().Get("search"))
	sortKey := strings.ToLower(r.URL.Query().Get("sort")) // title|album|date
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	list := allSongs
	if search != "" {
		var filtered []Song
		for _, s := range list {
			if strings.Contains(strings.ToLower(s.Title), search) ||
				strings.Contains(strings.ToLower(s.Album), search) ||
				strings.Contains(strings.ToLower(s.Genre), search) {
				filtered = append(filtered, s)
			}
		}
		list = filtered
	}

	if sortKey == "title" || sortKey == "album" || sortKey == "date" {
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				var ai, aj string
				switch sortKey {
				case "title":
					ai, aj = strings.ToLower(list[i].Title), strings.ToLower(list[j].Title)
				case "album":
					ai, aj = strings.ToLower(list[i].Album), strings.ToLower(list[j].Album)
				default:
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

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Songs retrieved", Data: paged})
}
