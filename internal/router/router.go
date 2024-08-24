package router

import (
	"log"
	"net/http"
)

func Router() {
	config := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	appHandler := middlewareLog(config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FilepathRoot)))))
	mux.Handle("/app/", appHandler)

	// /api
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", config.handlerReset)
	mux.Handle("GET /api/chirps", middlewareLog(http.HandlerFunc(handleGetChirps)))
	mux.Handle("GET /api/chirps/{chirpID}", middlewareLog(http.HandlerFunc(handleGetChirpByID)))
	mux.Handle("POST /api/chirps", middlewareLog(http.HandlerFunc(handlePostChirps)))

	// /admin
	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)

	server := &http.Server{
		Addr: 		":" + port,
		Handler: 	mux,
	}

	log.Printf("Serving files from %s on port: %s\n", FilepathRoot, port)
	log.Fatal(server.ListenAndServe())
}