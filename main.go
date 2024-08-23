package main

import (
	"fmt"
	"net/http"
)

const port = "8080"

func main() {
	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
	
	fileServer := http.FileServer(http.Dir("."))
	
	serverMux.Handle("/app/", http.StripPrefix("/app", fileServer))

	server := &http.Server{
		Addr: ":" + port,
		Handler: serverMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}