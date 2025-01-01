package downloader

import (
	"fmt"
	utils "go-downloader/internal"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type HTTPDownloader struct {
	workers int
	client  *http.Client
	resume  bool
	wg      sync.WaitGroup
	bar     *progressbar.ProgressBar
}

func NewHTTPDownloader(workers int, resume bool) (*HTTPDownloader, error) {
	return &HTTPDownloader{
		workers: workers,
		client:  &http.Client{},
		wg:      sync.WaitGroup{},
		resume:  resume,
	}, nil
}

func (d *HTTPDownloader) downloadMulti(URL string, Dst string, totalSize int) error {
	log.Println("Using Multi Part Download")
	d.bar = SetBar(totalSize)

	partSize := totalSize / d.workers
	// Create temporary directory to store part files
	partDir := fmt.Sprintf("%s/parts/", utils.GetFileDir(Dst))
	os.MkdirAll(partDir, 0777)
	defer os.RemoveAll(partDir)

	d.wg.Add(d.workers)

	rangeStart := 0
	for i := 0; i < d.workers; i++ {
		go func(i, rangeStart int) {
			rangeEnd := rangeStart + partSize

			// the last part should download the remaining bytes
			if i == d.workers-1 {
				rangeEnd = totalSize
			}

			downloadedSize := 0

			partFileName := d.getPartFileName(partDir, Dst, i)
			// use resume download
			if d.resume {
				fileInfo, err := os.Stat(partFileName)
				if err == nil {
					downloadedSize = int(fileInfo.Size())
				}
				if err != nil && !os.IsNotExist(err) {
					log.Fatalln(err)
				}
				d.bar.Add(downloadedSize)
			}

			d.downloadPart(URL, partFileName, rangeStart+downloadedSize, rangeEnd)
		}(i, rangeStart)

		rangeStart += partSize + 1
	}

	d.wg.Wait()
	d.mergeParts(partDir, Dst)

	return nil
}

// Do not support resume download from breakpoint, because range header is not supported
func (d *HTTPDownloader) downloadSingle(URL, Dst string) error {
	log.Println("Do not support multi-part download, downloading in single thread")
	resp, err := d.client.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// set bar based on content length
	d.bar = SetBar(int(resp.ContentLength))

	f, err := os.OpenFile(Dst, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(f, d.bar), resp.Body, buf)
	return err
}

func (d *HTTPDownloader) downloadPart(URL, partFileName string, rangeStart, rangeEnd int) {
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
	partFile, err := os.OpenFile(partFileName, flags, 0666)
	if err != nil {
		log.Fatalln("Download part err:", err)
	}
	defer partFile.Close()

	// use CopyBuffer to save memory
	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(partFile, resp.Body, buf)
	if err != nil {
		// ignore EOF error
		if err != io.EOF {
			return
		}
		log.Fatalln(err)
	}
}

func (d *HTTPDownloader) getPartFileName(partDir, Dst string, i int) string {
	filename := utils.GetFileName(Dst)
	return fmt.Sprintf("%s%s-%d.part", partDir, filename, i)
}

func (d *HTTPDownloader) DownloadFile(URL, Dst string) error {
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

	// check if Dst is available
	if _, err := os.Stat(Dst); err == nil {
		return fmt.Errorf("file already exists")
	}

	// check if the server supports multi-part download
	resp, err := d.client.Head(URL)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		return d.downloadMulti(URL, Dst, int(resp.ContentLength))
	}
	return d.downloadSingle(URL, Dst)
}

func (d *HTTPDownloader) mergeParts(partDir, Dst string) error {
	destFile, err := os.OpenFile(Dst, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < d.workers; i++ {
		partFileName := d.getPartFileName(partDir, Dst, i)
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
