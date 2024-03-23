package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		fmt.Fprintf(os.Stdout, "%s \"%s %s %s\" from %s - %d in %s\n", start.Format("2006-01-02 15:04:05"), r.Method, r.URL.Path, r.Proto, r.RemoteAddr, wrapped.statusCode, time.Since(start))
	})
}
