package download

type HttpDownloader struct {
	url string
	conc int
	path string
	downloadSlice []downloadSlice
}

type downloadSlice struct {
	url string
	path string
	fromIdx int
	toIdx int
	num int
}

func NewHttpDownloader(url string, conc int, path string) *HttpDownloader{
	// 发送head包，根据Accept-Ranges获得server是否支持并发下载

}
