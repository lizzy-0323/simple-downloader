package cmd

import (
	"downloader/internal/conf"
	"downloader/internal/downloader"
	"log"
	"net/url"

	"github.com/spf13/cobra"
)

type DownloadConfig struct {
	Dst     string
	Src     string
	Workers int
	ops     []string
}

var datasets = []string{"sift1m", "sift10k", "sift1b", "deep1b"}

func NewDownloadCmd() *cobra.Command {
	downloadConfig := &DownloadConfig{}
	downloadCmd := &cobra.Command{
		Use:       "download [dataset]",
		Short:     "download dataset",
		Long:      `Download dataset from the internet, e.g., sift1m, sift1b, deep1b.`,
		ValidArgs: datasets,
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dataset := args[0]
			url := getDatasetURL(dataset)
			if url == "" {
				log.Fatalf("Unknown dataset: %s", dataset)
			}
			d := getDownloader(url, downloadConfig.Workers)
			err := d.DownloadFile(dataset, url)
			if err != nil {
				log.Fatalf("Download failed: %v", err)
			}
			log.Printf("Download success: %s", dataset)
		},
	}
	downloadCmd.Flags().StringVarP(&downloadConfig.Dst, "dst", "d", ".", "destination directory")
	downloadCmd.Flags().StringVarP(&downloadConfig.Src, "src", "s", "", "source directory")
	downloadCmd.Flags().IntVarP(&downloadConfig.Workers, "workers", "w", 16, "number of workers")
	downloadCmd.Flags().StringSliceVar(&downloadConfig.ops, "ops", []string{}, "operations to run")
	return downloadCmd
}

func getDownloader(rawURL string, workers int) downloader.Downloader {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	switch u.Scheme {
	case "http", "https":
		return downloader.NewHTTPDownloader(workers, true)
	case "ftp":
		return downloader.NewFTPDownloader()
	default:
		return nil
	}
}

func getDatasetURL(dataset string) string {
	switch dataset {
	case "sift1m":
		return conf.SIFT_1M_URL
	case "deep1b":
		return conf.DEEP_1B_URL
	default:
		return ""
	}
}
