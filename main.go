package main

import (
	"net/http"
	"log"
)

func main() {
	config := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	appHandler := middlewareLog(config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/", appHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", config.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handleValidate)

	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)

	server := &http.Server{
		Addr: 		":" + port,
		Handler: 	mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}