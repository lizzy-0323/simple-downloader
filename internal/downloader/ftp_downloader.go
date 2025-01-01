package downloader

import (
	"fmt"
	utils "go-downloader/internal"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/jlaffaye/ftp"
	"github.com/schollz/progressbar/v3"
)

// FTPDownloader is a struct that holds the FTP client
type FTPDownloader struct {
	Port   int
	bar    *progressbar.ProgressBar
	client *ftp.ServerConn
}

// NewFTPDownloader creates a new FTPDownloader
func NewFTPDownloader(Port int) (*FTPDownloader, error) {
	return &FTPDownloader{
		Port: Port,
	}, nil
}

func (d *FTPDownloader) GetConn(u *url.URL) (*ftp.ServerConn, error) {
	addr := fmt.Sprintf("%s:%d", u.Host, d.Port)
	client, err := ftp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FTP server: %v", err)
	}

	return client, nil
}

// DownloadFile downloads a file from the FTP server
func (d *FTPDownloader) DownloadFile(URL, Dst string) error {
	fileName := utils.GetFileName(URL)
	if fileName == "" {
		return fmt.Errorf("invalid destination file")
	}
	err := utils.CreateDirIfNotExist(Dst)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// combine Dst and fileName
	Dst = path.Join(Dst, fileName)

	// Parse the remote URL
	u, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("failed to parse FTP URL: %v", err)
	}

	// Create a new FTP client
	d.client, err = d.GetConn(u)
	if err != nil {
		return err
	}

	if u.User != nil {
		password, _ := u.User.Password()
		err = d.client.Login(u.User.Username(), password)
		if err != nil {
			return fmt.Errorf("failed to login to FTP server: %v", err)
		}
	} else {
		err = d.client.Login("anonymous", "anonymous")
		if err != nil {
			return fmt.Errorf("failed to login to FTP server: %v", err)
		}
	}
	defer d.client.Quit()

	// Get the size of the file first, otherwise the server will end
	size, err := d.client.FileSize(u.Path)
	if err != nil {
		return fmt.Errorf("failed to get file size: %v", err)
	}

	resp, err := d.client.Retr(u.Path)
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
	d.bar = SetBar(int(size))

	// Copy the file content with progress bar
	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(localFile, d.bar), resp, buf)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}
