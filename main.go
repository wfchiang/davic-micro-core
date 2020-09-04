package main 

import (
	"fmt"
	"log"
	"os"
	"net/http"
	// "encoding/json"
	// "html/template"
	"github.com/gorilla/mux"
	// "io/ioutil"
	// "github.com/wfchiang/davic"
)

func pingHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	fmt.Fprintf(http_resp, "{}")	
}

func main () {
	log.Println("Starting Davic-Micro-Core...")
	mux_router := mux.NewRouter()

	mux_router.HandleFunc("/ping", pingHandler)

	http.Handle("/", mux_router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Use defailt port %s", port)
	} 
	
	log.Fatal(http.ListenAndServe(":"+port, nil))

}