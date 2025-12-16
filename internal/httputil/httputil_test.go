package httputil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "downloaded.txt")

	err := DownloadFile(server.URL, destPath)
	if err != nil {
		t.Errorf("DownloadFile() error = %v", err)
		return
	}

	// Verify file exists and has correct content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Errorf("Failed to read downloaded file: %v", err)
		return
	}

	if string(content) != "test content" {
		t.Errorf("Downloaded content = %v, want %v", string(content), "test content")
	}
}

func TestDownloadFileError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	destPath := filepath.Join(tmpDir, "downloaded.txt")

	err := DownloadFile(server.URL, destPath)
	if err == nil {
		t.Error("DownloadFile() expected error for 404 response")
	}
}

func TestFetchJSON(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"test","value":42}`))
	}))
	defer server.Close()

	var data TestData
	err := FetchJSON(server.URL, &data)
	if err != nil {
		t.Errorf("FetchJSON() error = %v", err)
		return
	}

	if data.Name != "test" || data.Value != 42 {
		t.Errorf("FetchJSON() got %+v, want {Name:test Value:42}", data)
	}
}

func TestFetchBytes(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test bytes"))
	}))
	defer server.Close()

	content, err := FetchBytes(server.URL)
	if err != nil {
		t.Errorf("FetchBytes() error = %v", err)
		return
	}

	if string(content) != "test bytes" {
		t.Errorf("FetchBytes() = %v, want %v", string(content), "test bytes")
	}
}
