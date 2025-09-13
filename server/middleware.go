package server

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// jsonLogger logs requests with structured JSON
func JsonLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		entry := map[string]interface{}{
			"ts":     start.UTC().Format(time.RFC3339Nano),
			"method": r.Method,
			"path":   r.URL.Path,
			"ip":     r.RemoteAddr,
			"ua":     r.UserAgent(),
			"dur_ms": time.Since(start).Milliseconds(),
		}
		b, _ := json.Marshal(entry)
		log.Println(string(b))
	})
}

// simple in-memory token bucket per IP (best effort)
type rateBucket struct {
	tokens int
	last   time.Time
}

var rlStore = map[string]*rateBucket{}

func RateLimit(maxPerMinute int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		b := rlStore[ip]
		now := time.Now()
		if b == nil {
			b = &rateBucket{tokens: maxPerMinute, last: now}
			rlStore[ip] = b
		}
		// refill
		elapsed := now.Sub(b.last)
		refill := int(elapsed.Minutes()) * maxPerMinute
		if refill > 0 {
			b.tokens = maxPerMinute
			b.last = now
		}
		if b.tokens <= 0 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"success":false,"message":"rate limit"}`))
			return
		}
		b.tokens--
		next.ServeHTTP(w, r)
	})
}

// computeETag returns a weak ETag for a byte slice
func computeETag(data []byte) string {
	sum := sha1.Sum(data)
	return "W/\"" + hex.EncodeToString(sum[:]) + "\""
}

// handleConditional writes 304 if If-None-Match or If-Modified-Since matches
func handleConditional(w http.ResponseWriter, r *http.Request, etag string, last time.Time) bool {
	if inm := r.Header.Get("If-None-Match"); inm != "" && inm == etag {
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	if ims := r.Header.Get("If-Modified-Since"); ims != "" {
		if t, err := time.Parse(http.TimeFormat, ims); err == nil && !last.After(t) {
			w.WriteHeader(http.StatusNotModified)
			return true
		}
	}
	w.Header().Set("ETag", etag)
	w.Header().Set("Last-Modified", last.UTC().Format(http.TimeFormat))
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(60))
	return false
}
