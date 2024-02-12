package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := &Message{
		Status: "success",
		Body:   "Successfully reached out service",
	}
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		return
	}
}

func main() {
	http.Handle("/ping", rateLimiter(handler))
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println("error listen and serve at port 8080:", err)
	}
}
