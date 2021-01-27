package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRenderProgressBadge(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/42/", ts.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	content_type, ok := resp.Header["Content-Type"]

	if !ok {
		t.Fatalf("Expected Content-Type header to be set")
	}

	if content_type[0] != "image/svg+xml" {
		t.Fatalf("Expected Content-Type=\"image/svg+xml\", got %s", content_type[0])
	}

}

func TestBadgeColors(t *testing.T) {
	ts := httptest.NewServer(setupServer())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/15/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(string(body), "#d9534f") {
		t.Fatalf("Expected color #d9534f to be present in SVG")
	}
	resp.Body.Close()

	resp, err = http.Get(fmt.Sprintf("%s/42/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(string(body), "#f0ad4e") {
		t.Fatalf("Expected color #f0ad4e to be present in SVG")
	}
	resp.Body.Close()

	resp, err = http.Get(fmt.Sprintf("%s/91/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(string(body), "#5cb85c") {
		t.Fatalf("Expected color #5cb85c to be present in SVG")
	}
	resp.Body.Close()
}
