package downloader_test

import (
	"downloader/internal/downloader"
	"testing"
)

func TestNewFTPDownloader(t *testing.T) {
	d, err := downloader.NewFTPDownloader("ftp.irisa.fr", 21)
	if err != nil {
		t.Fatalf("failed to create FTP downloader: %v", err)
	}
	t.Log(d)
}

func TestFTPDownloadFile(t *testing.T) {
	d, err := downloader.NewFTPDownloader("ftp.irisa.fr", 21)
	if err != nil {
		t.Fatalf("failed to create FTP downloader: %v", err)
	}
	defer d.Close()

	err = d.DownloadFile("local/texmex/corpus/siftsmall.tar.gz", "./siftsmall.tar.gz")
	if err != nil {
		t.Fatalf("failed to download file: %v", err)
	}
}
