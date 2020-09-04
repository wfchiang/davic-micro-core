package main 

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/wfchiang/davic"
)

func pingHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	fmt.Fprintf(http_resp, "{}")	
}

func main () {
	log.Println("Init File Server...")
	file_server := http.FileServer(http.Dir("./static/"))

	log.Println("Starting Davic-Micro-Core...")
	mux_router := mux.NewRouter()

	mux_router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", file_server))
	mux_router.HandleFunc("/ping", homepageHandler)

	http.Handle("/", mux_router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Use defailt port %s", port)
	} 
	
	log.Fatal(http.ListenAndServe(":"+port, nil))

}