package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()

	port := "8080"
	filepathRoot := "/app/"

	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	serveMux.Handle(filepathRoot, http.FileServer(http.Dir(".")))
	serveMux.HandleFunc("/healthz", healthzHandler)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthzHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(200)
	rw.Write(([]byte)("OK"))
}
