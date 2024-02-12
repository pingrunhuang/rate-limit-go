package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func perClientLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client, 0)
	)

	// remove client
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, cli := range clients {
				if time.Since(cli.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				message := &Message{
					Status: "failed",
					Body:   "ip limit error",
				}
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(message)
				mu.Unlock()
				return
			} else {
				next(w, r)
			}
			mu.Unlock()
		}

	})
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := &Message{Status: "success", Body: "successfully reached end point"}
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
	w.WriteHeader(http.StatusOK)
	message := &Message{Status: "success", Body: "hello world!"}
	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		return
	}
}

func main() {
	http.Handle("/ping", perClientLimiter(endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
