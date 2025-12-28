package middlewares

import (
	"log"
	"net"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new responseWriter
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// StatusCode returns the captured status code
func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

// LoggingMiddleware logs HTTP requests with IP, latency, URL, method, and status code
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timing
		start := time.Now()

		// Wrap the response writer to capture status code
		wrappedWriter := NewResponseWriter(w)

		// Get client IP address
		clientIP := getClientIP(r)

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// Calculate latency
		latency := time.Since(start)

		// Log the request details
		log.Printf("[%s] %s %s %s %d %v",
			clientIP,
			r.Method,
			r.URL.Path,
			r.Proto,
			wrappedWriter.StatusCode(),
			latency,
		)
	})
}

// getClientIP extracts the real client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if len(xff) > 0 {
			return xff
		}
	}

	// Check X-Real-IP header (for nginx)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Check X-Forwarded-Proto header
	if xfp := r.Header.Get("X-Forwarded-Proto"); xfp != "" {
		// This is typically used for protocol, but some setups use it for IP
		if xfp != "" {
			return xfp
		}
	}

	// Fall back to RemoteAddr
	if r.RemoteAddr != "" {
		// RemoteAddr includes port, so we need to extract just the IP
		// Format is typically "IP:port"
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			return host
		}
		return r.RemoteAddr
	}

	return "unknown"
}
