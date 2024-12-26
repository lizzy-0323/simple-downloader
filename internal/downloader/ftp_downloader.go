package downloader

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// FTPDownloader is a struct that holds the FTP client
type FTPDownloader struct {
	client *ftp.ServerConn
	bar    *progressbar.ProgressBar
}

func (d *FTPDownloader) SetBar(length int) {
	d.bar = progressbar.NewOptions(
		length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("downloading..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

// NewFTPDownloader creates a new FTPDownloader
func NewFTPDownloader(Host string, Port int) (*FTPDownloader, error) {
	URL := fmt.Sprintf("ftp://%s:%d", Host, Port)
	u, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FTP URL: %v", err)
	}

	client, err := ftp.Dial(u.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FTP server: %v", err)
	}

	if u.User != nil {
		password, _ := u.User.Password()
		err = client.Login(u.User.Username(), password)
		if err != nil {
			return nil, fmt.Errorf("failed to login to FTP server: %v", err)
		}
	} else {
		err = client.Login("anonymous", "anonymous")
		if err != nil {
			return nil, fmt.Errorf("failed to login to FTP server: %v", err)
		}
	}

	return &FTPDownloader{client: client}, nil
}

// DownloadFile downloads a file from the FTP server
func (d *FTPDownloader) DownloadFile(URL, Dst string) error {
	// Get the size of the file first, otherwise the server will end
	size, err := d.client.FileSize(URL)
	if err != nil {
		return fmt.Errorf("failed to get file size: %v", err)
	}

	resp, err := d.client.Retr(URL)
	if err != nil {
		return fmt.Errorf("failed to retrieve file from FTP server: %v", err)
	}
	defer resp.Close()

	localFile, err := os.Create(Dst)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	// Create a progress bar
	d.SetBar(int(size))

	// Copy the file content with progress bar
	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(localFile, d.bar), resp, buf)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}

// Close closes the FTP connection
func (d *FTPDownloader) Close() error {
	return d.client.Quit()
}
