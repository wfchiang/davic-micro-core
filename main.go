package main 

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"encoding/json"
	"strconv"
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

// Recovering function 
func recoverFromPanic (http_resp http.ResponseWriter, http_reqt *http.Request) {
	if r := recover() ; r != nil {
		err_message := fmt.Sprintf("{\"message\": \"%v\"}", r)
		http.Error(http_resp, err_message, http.StatusInternalServerError)
	}
}

// Ping handler 
func pingHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	grantCORS(&http_resp); 
	fmt.Fprintf(http_resp, "{}")	
}

// Operation list handler
func optListHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, http_reqt) 
	grantCORS(&http_resp); 

	if resp_body, err := json.Marshal(OPT_CACHE); err == nil {
		fmt.Fprintf(http_resp, string(resp_body))
	} else {		
		panic(fmt.Sprintf("Response marshalling failed for optList: %v", err))
	}
}

// Operation appending handler 
func optAppendingHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, http_reqt) 
	grantCORS(&http_resp); 

	// Read the request body 
	var reqt_obj []interface{}
	if bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body); err == nil {
		reqt_obj = davic.CreateArrayFromBytes(bytes_reqt_body)
	} else {
		panic(fmt.Sprintf("Failed to read the request body: %v", err))
	}

	// Append the opt to OPT_CACHE
	OPT_CACHE = append(OPT_CACHE, reqt_obj)
	log.Println(fmt.Sprintf("Opt-append -- OPT_CACHE: %v", OPT_CACHE))

	// Response 
	fmt.Fprintf(http_resp, "{\"message\":\"appended\"}")	
}

// Remove an operation from the list 
func optRemoveHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, http_reqt) 
	grantCORS(&http_resp); 

	// Get the assigned opt-id 
	req_vars := mux.Vars(http_reqt)
	if opt_id, err := strconv.Atoi(req_vars["opt-id"]); err == nil {
		OPT_CACHE = append(OPT_CACHE[0:opt_id], OPT_CACHE[opt_id+1:]...)
	} else {
		log.Println(fmt.Sprintf("[ERROR] Invalid opt-id: %v", req_vars["opt-id"]))
		panic(fmt.Sprintf("Invalid opt-id"))
	}

	log.Println(fmt.Sprintf("Opt-remove -- OPT_CACHE: %v", OPT_CACHE))

	// Response 
	fmt.Fprintf(http_resp, "{\"message\":\"removed\"}")	
}

func runDavicHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, http_reqt) 
	grantCORS(&http_resp); 

	// Read the request body 
	var reqt_obj interface{}
	if bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body); err == nil {
		reqt_obj = davic.CreateObjFromBytes(bytes_reqt_body)
		log.Println(fmt.Sprintf("Run Davic with initial store: %v", string(bytes_reqt_body)))
	} else {
		panic(fmt.Sprintf("Failed to read the request body: %v", err))
	}

	// Setup the Davic environment 
	env := davic.CreateNewEnvironment()
	env.Store = davic.CastInterfaceToObj(reqt_obj)

	// Run the operations 
	env = davic.Execute(env, OPT_CACHE) 

	// Response with the store 
	if resp_body, err := json.Marshal(env.Store); err == nil { 
		fmt.Fprintf(http_resp, string(resp_body))
	} else {
		panic(fmt.Sprintf("Response marshalling failed: %v", err))
	} 
}

// Main 
func main () {
	log.Println("Starting Davic-Micro-Core...")
	mux_router := mux.NewRouter()

	mux_router.HandleFunc("/ping", pingHandler).Methods("GET")
	mux_router.HandleFunc("/opt-list", optListHandler).Methods("GET")
	mux_router.HandleFunc("/opt-append", optAppendingHandler).Methods("POST")
	mux_router.HandleFunc("/opt-remove/{opt-id}", optRemoveHandler).Methods("GET")
	mux_router.HandleFunc("/run", runDavicHandler).Methods("POST")

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
			handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "OPTIONS"}), 
			handlers.AllowedOrigins([]string{"*"}))(mux_router)))
}