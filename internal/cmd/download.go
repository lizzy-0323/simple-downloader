package cmd

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [dataset]",
	Short: "download dataset",
	Long:  `Download dataset from the internet, e.g., sift1m, sift1b, deep1b.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dataset := args[0]
		url := getDatasetURL(dataset)
		if url == "" {
			log.Fatalf("Unknown dataset: %s", dataset)
		}

		err := downloadFile(dataset, url)
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
		log.Printf("Download success: %s", dataset)
	},
}

func init() {
	downloadCmd.Flags().StringP("dst", "d", ".", "destination directory")
	downloadCmd.Flags().StringP("src", "s", "", "source directory")
	downloadCmd.Flags().StringSlice("ops", []string{}, "operations to run")
}

func getDatasetURL(dataset string) string {
	switch dataset {
	case "sift1m":
		return "http://example.com/sift1m.zip"
	case "sift1b":
		return "http://example.com/sift1b.zip"
	case "deep1b":
		return "http://example.com/deep1b.zip"
	default:
		return ""
	}
}

func downloadFile(filename, url string) error {
	out, err := os.Create(filename + ".zip")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
