package handlers

import (
	"log"
	"net/http"
	"time"
)

func WithLogger(l *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // call original
		end := time.Now()
		l.Printf("%s %s %s %v\n",
			r.Method,
			r.RequestURI,
			r.UserAgent(),
			end.Sub(start))
	})
}
