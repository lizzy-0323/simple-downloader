package downloader

import (
	"fmt"
	utils "go-downloader/internal"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SCPDownloader struct {
	client *sftp.Client
}

func NewSCPConfig(user, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func (d *SCPDownloader) GetConn(u *url.URL) (*sftp.Client, error) {
	var port int
	var err error
	if u.Hostname() == "" {
		return nil, fmt.Errorf("invalid hostname")
	}
	if u.Port() == "" {
		return nil, fmt.Errorf("port is not supported")
	}
	port, err = strconv.Atoi(u.Port())
	if err != nil {
		return nil, fmt.Errorf("failed to parse port: %v", err)
	}
	addr := fmt.Sprintf("%s:%d", u.Hostname(), port)
	password, _ := u.User.Password()
	conn, err := ssh.Dial("tcp", addr, NewSCPConfig(u.User.Username(), password))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create sftp client: %v", err)
	}

	return client, nil
}

func NewSCPDownloader() (*SCPDownloader, error) {
	return &SCPDownloader{client: nil}, nil
}

func (d *SCPDownloader) DownloadFile(remotePath, localPath string) error {
	fileName := utils.GetFileName(remotePath)
	if fileName == "" {
		return fmt.Errorf("invalid destination file")
	}
	err := utils.CreateDirIfNotExist(localPath)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// combine localPath and fileName
	localPath = path.Join(localPath, fileName)

	// check if localPath is available
	if _, err := os.Stat(localPath); err == nil {
		return fmt.Errorf("file already exists")
	}

	// Parse the remote URL
	u, err := url.Parse(remotePath)
	if err != nil {
		return fmt.Errorf("failed to parse SCP URL: %v", err)
	}

	d.client, err = d.GetConn(u)
	if err != nil {
		return err
	}
	defer d.client.Close()

	// path := strings.TrimPrefix(u.Path, "/")
	path := u.Path
	srcFile, err := d.client.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer srcFile.Close()

	localPathFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localPathFile.Close()

	// Create a progress bar
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	bar := SetBar(int(fileInfo.Size()))

	// Copy the file content with progress bar
	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(localPathFile, bar), srcFile, buf)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	return nil
}
