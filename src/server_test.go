package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func httpTEST(handler http.HandlerFunc, method string, endpoint string, statusCODE int, postDATA string, expectedRESP string, t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	var jsonStr = []byte(postDATA)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != statusCODE {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, statusCODE)
	}

	// Check the response body is what we expect.
	expected := expectedRESP
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthz(t *testing.T) {
	httpTEST(http.HandlerFunc(Healthz), "GET", "/healthz", 200, "", `{"alive":true}`, t)
}

func TestEcho(t *testing.T) {
	httpTEST(http.HandlerFunc(Echo), "POST", "/api", 200, `{"Lat":35.14326, "Lon":-116.104}`, `{"lat":35.14326,"lon":-116.104}`, t)
}

func TestEcho400(t *testing.T) {
	httpTEST(http.HandlerFunc(Echo), "POST", "/api", 400, `{"lat":35.14326,"lon":"-116.104"}`, `{"message":"bad json"}`, t)
}
