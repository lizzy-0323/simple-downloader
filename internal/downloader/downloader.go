package downloader

import (
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type Downloader interface {
	DownloadFile(URL, Dst string) error
}

type DefaultDownloader struct {
	bar *progressbar.ProgressBar
}

func NewDefaultDownloader() *DefaultDownloader {
	return &DefaultDownloader{}
}

func SetBar(length int) *progressbar.ProgressBar {
	return progressbar.NewOptions(
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
