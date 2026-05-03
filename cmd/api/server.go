package api

import (
	"net/http"
)

func StartHealthServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
	})

	go http.ListenAndServe(":5002", mux)
}