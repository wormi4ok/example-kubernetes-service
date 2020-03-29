package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTreeHandler(t *testing.T) {
	srv := httptest.NewServer(helloHandler(helloMsg))
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, srv.URL+"/hello", nil)
	req.RequestURI = "" // Request.RequestURI can't be set in client requests
	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that HTTP status code is 200
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response code = %d, actual = %d", http.StatusOK, res.StatusCode)
	}

	// Check Content-Type header
	responseCT := res.Header.Get("Content-Type")
	if responseCT != "application/json" {
		t.Errorf("Expected response content type = %s, actual = %s", "application/json", responseCT)
	}

	// Compare response body with expected
	expected := []byte("{\"helloMsg\":\"Infrastructure enables innovation\"}\n")
	b, err := ioutil.ReadAll(res.Body)
	if err != nil || bytes.Compare(expected, b) != 0 {
		t.Errorf("Expected response body = %s, actual = %s", expected, b)
	}

}

func TestPortFromEnv(t *testing.T) {
	expected := fmt.Sprint(rand.Int())
	err := os.Setenv("PORT", expected)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	actual := portFromEnv()
	if expected != actual {
		t.Errorf("Expected port = %s, actual = %s", expected, actual)
	}
}
