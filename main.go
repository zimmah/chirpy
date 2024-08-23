package main

import (
	"fmt"
	"net/http"
)

const port = "8080"
type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		cfg.fileserverHits++
		fmt.Println("Request received")
		next.ServeHTTP(w, r)
	})
}

func main() {
	serverMux := http.NewServeMux()
	config := apiConfig{}

	

	serverMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

	serverMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf("Hits: %d", config.fileserverHits)))
    })

	serverMux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        config.fileserverHits = 0
    })
	
	fileServer := http.FileServer(http.Dir("."))
	
	serverMux.Handle("/app/", config.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	server := &http.Server{
		Addr: ":" + port,
		Handler: serverMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}