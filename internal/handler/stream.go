package handler

import (
	"encoding/json"
	"net/http"

	"github.com/michaelzhan1/recent-max/internal/value"
)

func DataStreamHandlerFactory(dataChan chan value.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		for {
			select {
			case <-r.Context().Done():
				return // client disconnect
			case msg := <-dataChan:
				payload, err := json.Marshal(msg)
				if err != nil {
					http.Error(w, "Error encoding response", http.StatusInternalServerError)
					return
				}

				if _, err := w.Write([]byte("data: ")); err != nil {
					return
				}
				if _, err := w.Write(payload); err != nil {
					return
				}
				if _, err := w.Write([]byte("\n\n")); err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}
