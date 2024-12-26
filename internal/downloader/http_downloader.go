package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type HTTPDownloader struct {
	workers int
	client  *http.Client
	wg      sync.WaitGroup
	bar     *progressbar.ProgressBar
}

func NewHTTPDownloader(workers int) *HTTPDownloader {
	return &HTTPDownloader{
		workers: workers,
		client:  &http.Client{},
		wg:      sync.WaitGroup{},
	}
}

func (d *HTTPDownloader) SetBar(length int) {
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

func (d *HTTPDownloader) downloadMulti(URL string, Dst string, totalSize int) error {
	log.Println("Using Multi Part Download")
	d.SetBar(totalSize)

	partSize := totalSize / d.workers
	// Create temporary directory to store part files
	partDir := d.getPartFileDir(Dst)
	os.Mkdir(partDir, 0755)
	defer os.RemoveAll(partDir)

	d.wg.Add(d.workers)

	rangeStart := 0
	for i := 0; i < d.workers; i++ {
		go func(i, rangeStart int) {
			rangeEnd := rangeStart + partSize
			if i == d.workers-1 {
				rangeEnd = totalSize
			}
			d.downloadPart(URL, partDir, rangeStart, rangeEnd)
		}(i, rangeStart)
		rangeStart += partSize + 1
	}

	d.wg.Wait()
	d.mergeParts(Dst)

	return nil
}

func (d *HTTPDownloader) downloadSingle(URL, Dst string) error {
	log.Println("Unsupport multi-part download, downloading in single thread")
	resp, err := d.client.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// set bar based on content length
	d.SetBar(int(resp.ContentLength))

	f, err := os.OpenFile(Dst, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(f, d.bar), resp.Body, buf)
	return err
}

func (d *HTTPDownloader) downloadPart(URL, Dst string, rangeStart, rangeEnd int) {
	defer d.wg.Done()
	if rangeStart >= rangeEnd {
		return
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	resp, err := d.client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	partFile, err := os.OpenFile(d.getPartFileName(Dst, rangeStart), flags, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	// use CopyBuffer to save memory
	_, err = io.CopyBuffer(partFile, resp.Body, buf)
	if err != nil {
		// ignore EOF error
		if err != io.EOF {
			return
		}
		log.Fatalln(err)
	}
}

func (d *HTTPDownloader) getPartFileName(Dst string, i int) string {
	return path.Join(Dst, fmt.Sprintf("%d.part", i))
}

func (d *HTTPDownloader) getPartFileDir(Dst string) string {
	// get the file name without extension
	return strings.SplitN(Dst, ".", 2)[0]
}

func (d *HTTPDownloader) DownloadFile(URL, Dst string) error {
	if Dst == "" {
		Dst = path.Base(URL)
	}
	resp, err := d.client.Head(URL)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		return d.downloadMulti(URL, Dst, int(resp.ContentLength))
	}
	return d.downloadSingle(URL, Dst)
}

func (d *HTTPDownloader) mergeParts(Dst string) error {
	destFile, err := os.OpenFile(Dst, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < d.workers; i++ {
		partFileName := d.getPartFileName(Dst, i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		io.Copy(destFile, partFile)
		partFile.Close()
		os.Remove(partFileName)
	}

	return nil
}
