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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test content"))
	}))
	defer server.Close()
	workers := 16
	downloader := downloader.NewHTTPDownloader(workers)
	dst := "./test.zip"
	url := "https://github.com/schollz/progressbar/archive/refs/tags/v3.17.1.zip"

	err := downloader.DownloadFile(url, dst)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = os.Stat(dst)
	if os.IsNotExist(err) {
		t.Fatalf("expected file to exist, got %v", err)
	}

	// os.Remove(dst)
}
