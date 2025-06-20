package server

import (
	"fmt"
	"github.com/rafa-mori/golife/internal"
	"github.com/rafa-mori/logz"
	"net/http"
)

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Set appropriate headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a channel to receive events
	events := make(chan internal.IManagedProcessEvents)

	// Register the channel to receive events from the lifecycle manager
	//internal.RegisterEventChannel(events)

	// Handle client disconnection
	ctx := r.Context()
	go func() {
		<-ctx.Done()
		close(events)
		logz.Info("Client disconnected", nil)
	}()

	// Stream events to the client
	for event := range events {
		fmt.Fprintf(w, "data: %v\n\n", event)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func RegisterSSEEndpoint(mux *http.ServeMux) {
	mux.HandleFunc("/events", SSEHandler)
}
