package downloader_test

import (
	"go-downloader/internal/downloader"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSCPDownloader(t *testing.T) {
	downloader, err := downloader.NewSCPDownloader()
	assert.NoError(t, err)
	assert.NotNil(t, downloader)
}

func TestGetConn(t *testing.T) {
	downloader, err := downloader.NewSCPDownloader()
	assert.NoError(t, err)
	assert.NotNil(t, downloader)

	URL := "scp://liziyi:null@localhost:22/test-scp.txt"
	u, err := url.Parse(URL)
	_, err = downloader.GetConn(u)
	assert.NoError(t, err)
}

func TestDownloadFile(t *testing.T) {
	downloader, err := downloader.NewSCPDownloader()
	assert.NoError(t, err)
	assert.NotNil(t, downloader)

	// currently only support absolute path
	remotePath := "scp://liziyi:null@localhost:22/Users/liziyi/test-scp.txt"
	localPath := "."
	err = downloader.DownloadFile(remotePath, localPath)
	assert.NoError(t, err)
}
