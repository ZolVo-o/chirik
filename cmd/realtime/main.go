package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type hub struct {
	mu          sync.Mutex
	subscribers map[chan []byte]struct{}
}

func newHub() *hub {
	return &hub{subscribers: make(map[chan []byte]struct{})}
}

func (h *hub) subscribe() chan []byte {
	channel := make(chan []byte, 16)
	h.mu.Lock()
	h.subscribers[channel] = struct{}{}
	h.mu.Unlock()
	return channel
}

func (h *hub) unsubscribe(channel chan []byte) {
	h.mu.Lock()
	delete(h.subscribers, channel)
	close(channel)
	h.mu.Unlock()
}

func (h *hub) publish(payload []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for subscriber := range h.subscribers {
		select {
		case subscriber <- payload:
		default:
			// A slow client must not block the rest of the broadcast fan-out.
		}
	}
}

func main() {
	secret := os.Getenv("CHIRIK_REALTIME_SECRET")
	if secret == "" {
		log.Fatal("CHIRIK_REALTIME_SECRET is required")
	}

	h := newHub()
	allowedOrigin := os.Getenv("CHIRIK_ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		subscriber := h.subscribe()
		defer h.unsubscribe(subscriber)
		fmt.Fprint(w, ": connected\n\n")
		flusher.Flush()
		keepAlive := time.NewTicker(20 * time.Second)
		defer keepAlive.Stop()

		for {
			select {
			case payload := <-subscriber:
				fmt.Fprintf(w, "data: %s\n\n", payload)
				flusher.Flush()
			case <-keepAlive.C:
				fmt.Fprint(w, ": ping\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})
	mux.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("X-Internal-Secret") != secret {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var incoming event
		if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil || incoming.Type == "" {
			http.Error(w, "Invalid event", http.StatusBadRequest)
			return
		}
		payload, err := json.Marshal(incoming)
		if err != nil {
			http.Error(w, "Invalid event", http.StatusBadRequest)
			return
		}
		h.publish(payload)
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := os.Getenv("CHIRIK_REALTIME_ADDR")
	if addr == "" {
		addr = ":8090"
	}
	srv := &http.Server{Addr: addr, Handler: mux}
	log.Printf("realtime server listening on %s", addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
