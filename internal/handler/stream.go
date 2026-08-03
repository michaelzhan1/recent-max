package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/michaelzhan1/recent-max/internal/value/deque"
)

func StreamHandlerFactory(dq *deque.ValueDeque) http.HandlerFunc {
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

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		type response struct {
			MaxValue *float64 `json:"maxValue"`
		}

		for {
			select {
			case <-r.Context().Done():
				return // client disconnect
			case <-ticker.C:
				v, ok := dq.Peek()
				var resp response
				if ok {
					resp.MaxValue = &(v.Value)
				} else {
					resp.MaxValue = nil
				}
				payload, err := json.Marshal(resp)
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
