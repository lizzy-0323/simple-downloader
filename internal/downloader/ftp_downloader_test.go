package downloader_test

import (
	"go-downloader/internal/downloader"
	"net/url"
	"testing"
)

func TestNewFTPDownloader(t *testing.T) {
	d, _ := downloader.NewFTPDownloader(21)

	URL := "ftp://ftp.irisa.fr/local/texmex/corpus/siftsmall.tar.gz"
	u, err := url.Parse(URL)
	if err != nil {
		t.Fatalf("failed to parse FTP URL: %v", err)
	}
	_, err = d.GetConn(u)
	if err != nil {
		t.Fatalf("failed to create FTP downloader: %v", err)
	}
	t.Log(d)
}

func TestFTPDownloadFile(t *testing.T) {
	d, _ := downloader.NewFTPDownloader(21)
	URL := "ftp://ftp.irisa.fr/local/texmex/corpus/siftsmall.tar.gz"
	Dst := "./data"
	err := d.DownloadFile(URL, Dst)
	if err != nil {
		t.Fatalf("failed to download file: %v", err)
	}
}
