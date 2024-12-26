package downloader

type downloader interface {
	DownloadFile(URL, Dst string) error
	SetBar(length int)
}
