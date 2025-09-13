package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// AlbumsHandler handles requests to /api/albums
func AlbumsHandler(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Only allow GET method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	// Get query parameters
	query := r.URL.Query()
	idParam := query.Get("id")
	titleParam := query.Get("title")
	searchParam := query.Get("search")
	sortParam := query.Get("sort")
	pageParam := query.Get("page")
	limitParam := query.Get("limit")

	// Convert pagination params
	page := 1
	limit := 10
	if pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Filter albums
	var filteredAlbums []Album
	for _, album := range QWERBandData.Discography {
		// Filter by ID
		if idParam != "" {
			if id, err := strconv.Atoi(idParam); err == nil {
				if album.ID != id {
					continue
				}
			}
		}

		// Filter by title
		if titleParam != "" {
			if !strings.Contains(strings.ToLower(album.Title), strings.ToLower(titleParam)) {
				continue
			}
		}

		// Filter by search (title or type)
		if searchParam != "" {
			searchLower := strings.ToLower(searchParam)
			if !strings.Contains(strings.ToLower(album.Title), searchLower) &&
				!strings.Contains(strings.ToLower(album.Type), searchLower) {
				continue
			}
		}

		filteredAlbums = append(filteredAlbums, album)
	}

	// Sort albums
	if sortParam != "" {
		switch sortParam {
		case "title":
			// Sort by title A-Z
			for i := 0; i < len(filteredAlbums)-1; i++ {
				for j := i + 1; j < len(filteredAlbums); j++ {
					if strings.ToLower(filteredAlbums[i].Title) > strings.ToLower(filteredAlbums[j].Title) {
						filteredAlbums[i], filteredAlbums[j] = filteredAlbums[j], filteredAlbums[i]
					}
				}
			}
		case "date":
			// Sort by release date (newest first)
			for i := 0; i < len(filteredAlbums)-1; i++ {
				for j := i + 1; j < len(filteredAlbums); j++ {
					if filteredAlbums[i].ReleaseDate < filteredAlbums[j].ReleaseDate {
						filteredAlbums[i], filteredAlbums[j] = filteredAlbums[j], filteredAlbums[i]
					}
				}
			}
		}
	}

	// Calculate pagination
	total := len(filteredAlbums)
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= total {
		filteredAlbums = []Album{}
	} else {
		if endIndex > total {
			endIndex = total
		}
		filteredAlbums = filteredAlbums[startIndex:endIndex]
	}

	response := APIResponse{
		Success: true,
		Message: "Albums retrieved successfully",
		Data:    filteredAlbums,
		Meta: map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": (total + limit - 1) / limit,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
