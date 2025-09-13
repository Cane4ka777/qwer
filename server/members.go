package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handler for GET /api/members - Get all members or specific member
func MembersHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if requesting specific member by ID
	idParam := r.URL.Query().Get("id")
	nameParam := r.URL.Query().Get("name")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response := APIResponse{
				Success: false,
				Message: "Invalid member ID",
				Data:    nil,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		for _, member := range QWERBandData.Members {
			if member.ID == id {
				response := APIResponse{
					Success: true,
					Message: "Member information retrieved successfully",
					Data:    member,
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		response := APIResponse{
			Success: false,
			Message: "Member not found",
			Data:    nil,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if nameParam != "" {
		nameParam = strings.ToLower(nameParam)
		for _, member := range QWERBandData.Members {
			if strings.ToLower(member.StageName) == nameParam || strings.ToLower(member.Name) == nameParam {
				response := APIResponse{
					Success: true,
					Message: "Member information retrieved successfully",
					Data:    member,
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		response := APIResponse{
			Success: false,
			Message: "Member not found",
			Data:    nil,
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Enhanced list mode: search, sort, pagination
	search := strings.ToLower(r.URL.Query().Get("search"))
	sortKey := strings.ToLower(r.URL.Query().Get("sort")) // name|stage_name
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	list := QWERBandData.Members

	if search != "" {
		var filtered []Member
		for _, m := range list {
			if strings.Contains(strings.ToLower(m.Name), search) ||
				strings.Contains(strings.ToLower(m.StageName), search) ||
				strings.Contains(strings.ToLower(m.Position), search) {
				filtered = append(filtered, m)
			}
		}
		list = filtered
	}

	if sortKey == "name" || sortKey == "stage_name" {
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				ai := strings.ToLower(func() string {
					if sortKey == "name" {
						return list[i].Name
					}
					return list[i].StageName
				}())
				aj := strings.ToLower(func() string {
					if sortKey == "name" {
						return list[j].Name
					}
					return list[j].StageName
				}())
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

	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Members retrieved", Data: paged})
}
