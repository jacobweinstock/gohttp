package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type loc struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type decodeError struct {
	Message string `json:"message"`
}

type health struct {
	Alive bool `json:"alive"`
}

// Log handle logging to stdout
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Echo this return what you send it
func Echo(w http.ResponseWriter, r *http.Request) {
	// First, we need a location struct to receive the decoded data.
	location := loc{}

	// The location data is inside the request body which is an io.ReadCloser,
	// but we need a byte slice for unmarshalling.
	// ReadAll from ioutil just comes in handy.
	// Note we use ReadAll here for simplicity. Be careful when using ReadAll in larger
	// projects, as reading large files can consume a lot of memory.
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading the body %s", err)
	}
	// Now we can decode the request data using the Unmarshal function.
	err = json.Unmarshal(jsn, &location)
	if err != nil {
		log.Printf("Decoding error: %s", err)
		errorData := decodeError{"bad json"}
		errorJSON, err := json.Marshal(errorData)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorJSON)
	} else {
		// We send a JSON response, so we need to set the Content-Type header accordingly.
		w.Header().Set("Content-Type", "application/json")

		// Sending the response is as easy as writing to the ResponseWriter object.
		returnJSON, err := json.Marshal(location)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(returnJSON)
		// To see if the request was correctly received, let's print it to the console.
		log.Printf("Received: %v\n", location)
	}
}

// Healthz basic health check
func Healthz(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	health, err := json.Marshal(health{true})
	if err != nil {
		log.Printf("Error: %s", err)
	}
	w.Write(health)
}

func main() {
	port := "8080"
	log.Println(fmt.Sprintf("starting service on http://localhost:%s", port))
	http.HandleFunc("/api", Echo)
	http.HandleFunc("/healthz", Healthz)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), Log(http.DefaultServeMux)))
}
