package main 

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"
	// "html/template"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"io/ioutil"
	"github.com/wfchiang/davic"
)

// Global variables 
var OPT_CACHE = []interface{}{} 

// Helper functions 
func grantCORS (http_resp *http.ResponseWriter) {
	(*http_resp).Header().Set("Access-Control-Allow-Origin", "*"); 
}

// Ping handler 
func pingHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	grantCORS(&http_resp); 
	fmt.Fprintf(http_resp, "{}")	
}

// Operation list handler
func optListHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	grantCORS(&http_resp); 

	if resp_body, err := json.Marshal(OPT_CACHE); err == nil {
		fmt.Fprintf(http_resp, string(resp_body))
	} else {		
		panic(fmt.Sprintf("Response marshalling failed for optList: %v", err))
	}
}

// Operation appending handler 
func optAppendingHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	grantCORS(&http_resp); 

	// Read the request body 
	var reqt_obj interface{}
	if bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body); err == nil {
		reqt_obj = davic.CreateArrayFromBytes(bytes_reqt_body)
		log.Println(fmt.Sprintf("Opt-append requested: \n    string: %v\n    obj: %v", string(bytes_reqt_body), reqt_obj))
	} else {
		panic(fmt.Sprintf("Failed to read the request body: %v", err))
	}

	// Append the opt to OPT_CACHE
	OPT_CACHE = append(OPT_CACHE, reqt_obj)
	log.Println(fmt.Sprintf("Opt-append -- OPT_CACHE: %v", OPT_CACHE))

	// Response 
	fmt.Fprintf(http_resp, "{\"message\":\"appended\"}")	
}

// Main 
func main () {
	log.Println("Starting Davic-Micro-Core...")
	mux_router := mux.NewRouter()

	mux_router.HandleFunc("/ping", pingHandler).Methods("GET")
	mux_router.HandleFunc("/opt-list", optListHandler).Methods("GET")
	mux_router.HandleFunc("/opt-append", optAppendingHandler).Methods("POST")

	// http.Handle("/", mux_router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Use defailt port %s", port)
	} 
	
	log.Fatal(http.ListenAndServe(
		":"+port, 
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), 
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), 
			handlers.AllowedOrigins([]string{"*"}))(mux_router)))
}