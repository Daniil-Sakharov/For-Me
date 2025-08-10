package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	ts := httptest.NewServer(NewMux())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/hello")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != "Hello, World!\n" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestNotFound(t *testing.T) {
	ts := httptest.NewServer(NewMux())

	resp, err := http.Get(ts.URL + "/notfound")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
