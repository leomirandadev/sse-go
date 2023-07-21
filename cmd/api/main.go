package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var chanDashBoard = make(chan DashBoard, 10)

func main() {
	go sendMessages()

	http.HandleFunc("/sse", dashboardHandler)
	log.Println("serving on localhost:8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Println("connected")
	defer func() {
		log.Println("disconnected")
	}()

	receiveMessages(w, r)
}

type DashBoard struct {
	User uint `json:"user"`
}

func receiveMessages(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("flusher error")
		return
	}

	for {
		select {
		case ev, openChannel := <-chanDashBoard:
			if !openChannel {
				return
			}

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(ev)
			fmt.Fprintln(w, "data:", buf.String())
			f.Flush()

		case <-r.Context().Done():
			return
		}
	}

}

func sendMessages() {
	for range time.Tick(2 * time.Second) {
		chanDashBoard <- DashBoard{
			User: uint(rand.Uint32()),
		}
	}
}
