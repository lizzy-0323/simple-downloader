package downloader_test

import (
	"downloader/internal/downloader"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHTTPDownloadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", "2048") // Simulate a 1KB file
		w.Write([]byte("test content"))
	}))
	defer server.Close()
	workers := 3
	resume := true
	d := downloader.NewHTTPDownloader(workers, resume)
	dst := "./data/test.zip"
	// url := "https://github.com/schollz/progressbar/archive/refs/tags/v3.17.1.zip"

	err := d.DownloadFile(server.URL, dst)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = os.Stat(dst)
	if os.IsNotExist(err) {
		t.Fatalf("expected file to exist, got %v", err)
	}

	// defer os.Remove(dst)
}
