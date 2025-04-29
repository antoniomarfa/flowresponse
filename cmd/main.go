package main

import (
	"log"
	"net/http"

	"flowresponse/handles"
)

func main() {

	// Iniciar el servidor HTTP en una goroutine
	http.HandleFunc("/token", handles.HandleToken)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
