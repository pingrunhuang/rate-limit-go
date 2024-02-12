package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimiter(next func(http.ResponseWriter, *http.Request)) http.Handler {
	var limit rate.Limit = 2
	var b int = 4
	limiter := rate.NewLimiter(limit, b)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "failed",
				Body:   fmt.Sprintf("Reached rate %v and permits bursts of at most %v tokens.", limit, b),
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		} else {
			next(w, r)
		}
	})
}
