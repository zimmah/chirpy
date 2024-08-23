package main

import (
	"net/http"
	"log"
)

func main() {
	config := apiConfig{
		fileserverHits: 0,
		templatePath: 	"./admin/index.html",
		port:			"8080",
		filepathRoot:	".",
	}

	mux := http.NewServeMux()
	appHandler := middlewareLog(config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(config.filepathRoot)))))
	mux.Handle("/app/", appHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", config.handlerReset)
	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)

	server := &http.Server{
		Addr: 		":" + config.port,
		Handler: 	mux,
	}

	log.Printf("Serving files from %s on port: %s\n", config.filepathRoot, config.port)
	log.Fatal(server.ListenAndServe())

}