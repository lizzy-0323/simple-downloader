package cmd

import (
	"downloader/internal/downloader"
	"errors"
	"log"
	"net/url"

	"github.com/spf13/cobra"
)

type DownloadConfig struct {
	Dst     string
	Workers int
}

// NewDownloadCmd returns a new download command
func NewDownloadCmd() *cobra.Command {
	downloadConfig := &DownloadConfig{}
	downloadCmd := &cobra.Command{
		Use:   "download [URL]",
		Short: "download files with ftp or http protocol",
		Long:  "download files with ftp or http protocol",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			d, err := getDownloader(url, downloadConfig.Workers)
			if err != nil {
				log.Fatalf("Failed to get downloader: %v", err)
			}
			err = d.DownloadFile(url, downloadConfig.Dst)
			if err != nil {
				log.Fatalf("Download failed: %v", err)
			}
			log.Printf("Download success")
		},
	}
	downloadCmd.Flags().StringVarP(&downloadConfig.Dst, "dst", "d", ".", "destination directory to save the file")
	downloadCmd.Flags().IntVarP(&downloadConfig.Workers, "workers", "w", 16, "number of workers")
	return downloadCmd
}

func getDownloader(rawURL string, workers int) (downloader.Downloader, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http":
		return downloader.NewHTTPDownloader(workers, true), nil
	case "ftp":
		return downloader.NewFTPDownloader(), nil
	default:
		return nil, errors.New("unsupported protocol")
	}
}
