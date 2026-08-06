package handler

import "net/http"

func PauseHandlerFactory(pauseChan chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pauseChan <- true
		w.WriteHeader(http.StatusOK)
	}
}

func ResumeHandlerFactory(pauseChan chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pauseChan <- false
		w.WriteHeader(http.StatusOK)
	}
}
