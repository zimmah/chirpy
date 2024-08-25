package router

import (
	"log"
	"net/http"
	"os"
)

func Router() {
	jwtSecret := os.Getenv("JWT_SECRET")
	config := apiConfig{
		fileserverHits: 0,
		jwtSecret: []byte(jwtSecret),
	}

	mux := http.NewServeMux()
	appHandler := middlewareLog(config.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FilepathRoot)))))
	mux.Handle("/app/", appHandler)

	// /api
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", config.handlerReset)
	// /api/chirps
	mux.Handle("GET /api/chirps", middlewareLog(http.HandlerFunc(handleGetChirps)))
	mux.Handle("GET /api/chirps/{chirpID}", middlewareLog(http.HandlerFunc(handleGetChirpByID)))
	mux.Handle("POST /api/chirps", middlewareLog(http.HandlerFunc(handlePostChirps)))
	// /api/users
	mux.Handle("GET /api/users", middlewareLog(http.HandlerFunc(handleGetUsers)))
	mux.Handle("GET /api/users/{userID}", middlewareLog(http.HandlerFunc(handleGetUserByID)))
	mux.Handle("POST /api/users", middlewareLog(http.HandlerFunc(handlePostUsers)))
	mux.Handle("PUT /api/users", middlewareLog(http.HandlerFunc(config.handlePutUsers)))
	// /api/login
	mux.Handle("POST /api/login", middlewareLog(http.HandlerFunc(config.handleLogin)))

	// /admin
	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)

	server := &http.Server{
		Addr: 		":" + port,
		Handler: 	mux,
	}

	log.Printf("Serving files from %s on port: %s\n", FilepathRoot, port)
	log.Fatal(server.ListenAndServe())
}