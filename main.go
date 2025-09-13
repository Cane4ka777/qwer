package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"qwer-api/server"
)

func serverMiddleware(mux *http.ServeMux) http.Handler {
	// chain: rate limit -> json logger -> mux
	return server.JsonLogger(server.RateLimit(120, mux))
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

	// API routes (primary)
	mux.HandleFunc("/api", server.IndexHandler)
	mux.HandleFunc("/api/", server.IndexHandler)
	mux.HandleFunc("/api/band", server.BandHandler)
	mux.HandleFunc("/api/members", server.MembersHandler)
	mux.HandleFunc("/api/songs", server.SongsHandler)
	mux.HandleFunc("/api/albums", server.AlbumsHandler)
	mux.HandleFunc("/api/awards", server.AwardsHandler)

	// Aliases for compatibility
	mux.HandleFunc("/server", server.IndexHandler)
	mux.HandleFunc("/server/", server.IndexHandler)
	mux.HandleFunc("/server/band", server.BandHandler)
	mux.HandleFunc("/server/members", server.MembersHandler)
	mux.HandleFunc("/server/songs", server.SongsHandler)
	mux.HandleFunc("/server/albums", server.AlbumsHandler)
	mux.HandleFunc("/server/awards", server.AwardsHandler)

	mux.HandleFunc("/docs", server.DocsHandler)

	// Serve OpenAPI spec file under both paths
	mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.yaml")
	})
	mux.HandleFunc("/openserver.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.yaml")
	})

	fmt.Printf("🎸 QWER Band API Server starting on port %s...\n", port)
	fmt.Printf("📋 Available endpoints:\n")
	fmt.Printf("  • GET http://localhost:%s/        - Website\n", port)
	fmt.Printf("  • GET http://localhost:%s/api     - API Index\n", port)
	fmt.Printf("  • GET http://localhost:%s/docs    - API Docs (ReDoc)\n", port)
	fmt.Printf("\n🚀 Server ready! Press Ctrl+C to stop.\n\n")

	handler := serverMiddleware(mux)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
