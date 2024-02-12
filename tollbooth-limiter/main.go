package main

import (
	"encoding/json"
	"net/http"

	tollbooth "github.com/didip/tollbooth/v7"
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

	msg := &Message{Status: "failed", Body: "too many request"}
	jsonMsg, _ := json.Marshal(msg)
	tollLimiter := tollbooth.NewLimiter(1, nil)
	tollLimiter.SetMessageContentType("application/json")
	tollLimiter.SetMessage(string(jsonMsg))
	http.Handle("/ping", tollbooth.LimitFuncHandler(tollLimiter, handler))
	http.ListenAndServe(":8080", nil)

}
