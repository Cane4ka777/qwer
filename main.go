package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"qwer-api/api"
)

func apiMiddleware(mux *http.ServeMux) http.Handler {
	// chain: rate limit -> json logger -> mux
	return api.JsonLogger(api.RateLimit(120, mux))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Static website (serves public/index.html at / and other assets)
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/", fs)

	// API routes
	mux.HandleFunc("/api", api.IndexHandler)
	mux.HandleFunc("/api/", api.IndexHandler)
	mux.HandleFunc("/api/band", api.BandHandler)
	mux.HandleFunc("/api/members", api.MembersHandler)
	mux.HandleFunc("/api/songs", api.SongsHandler)
	mux.HandleFunc("/api/albums", api.AlbumsHandler)
	mux.HandleFunc("/api/awards", api.AwardsHandler)
	mux.HandleFunc("/docs", api.DocsHandler)

	// Serve OpenAPI spec file
	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.yaml")
	})

	fmt.Printf("🎸 QWER Band API Server starting on port %s...\n", port)
	fmt.Printf("📋 Available endpoints:\n")
	fmt.Printf("  • GET http://localhost:%s/        - Website\n", port)
	fmt.Printf("  • GET http://localhost:%s/api    - API Documentation\n", port)
	fmt.Printf("  • GET http://localhost:%s/docs   - API Docs (ReDoc)\n", port)
	fmt.Printf("\n🚀 Server ready! Press Ctrl+C to stop.\n\n")

	handler := apiMiddleware(mux)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
