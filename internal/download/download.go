package download

import (
	"bird/internal/err"
	"bird/internal/state"
	"bird/internal/tool"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type HttpDownloader struct {
	Url           string
	Conc          int
	Path          string
	DownloadSlice []*downloadSlice
	Resumeable    bool
	Clen          int
	Interrupte    bool
}

type downloadSlice struct {
	Url     string
	Path    string
	FromIdx int
	ToIdx   int
	Num     int
}

var (
	client = http.Client{}
)

func NewHttpDownloader(url string, conc int) *HttpDownloader {
	resume := true
	// 发送head包，根据Accept-Ranges获得server是否支持并发下载
	resp, e := http.Head(url)
	if e != nil {
		err.Handler(e)
	}
	ar := resp.Header.Get("Accept-Ranges")
	if ar != "bytes" {
		conc = 1
		resume = false
	}
	log.Printf("concurency num: %v \n", conc)
	// 获取body长度
	len := resp.Header.Get("Content-Length")
	log.Printf("total size: %v \n", len)
	vlen, e := strconv.Atoi(len)
	if e != nil {
		err.Handler(e)
	}
	slices := []*downloadSlice{}

	path := tool.GetCacheFolder(url)
	log.Printf("cache folder: %v \n", path)

	for i := 0; i < conc; i++ {
		rg := vlen / conc
		fromidx := i * rg
		toidx := (i+1)*rg - 1
		if i == conc-1 {
			toidx = vlen
		}
		ds := &downloadSlice{
			Url:     url,
			Path:    filepath.Join(path, fmt.Sprintf("%v_%v", filepath.Base(url), i)),
			FromIdx: fromidx,
			ToIdx:   toidx,
			Num:     i,
		}
		slices = append(slices, ds)
	}
	dl := &HttpDownloader{
		Url:           url,
		Conc:          conc,
		Path:          path,
		DownloadSlice: slices,
		Resumeable:    resume,
		Clen:          vlen,
		Interrupte:    false,
	}
	return dl

}

// func RecoverDownloader(state *state.State) *HttpDownloader {

// }

func (downloader *HttpDownloader) Download(errChan chan error, interruptChan chan bool, stateChan chan *state.State, downChan chan bool, fileChan chan string) {
	wg := sync.WaitGroup{}
	for i := 0; i < downloader.Conc; i++ {
		wg.Add(1)
		go func(slice *downloadSlice, downloader *HttpDownloader) {
			log.Printf("download: %v \n", slice)
			defer wg.Done()
			req, e := http.NewRequest(http.MethodGet, slice.Url, nil)
			if e != nil {
				errChan <- e
				return
			}
			if downloader.Conc > 1 {
				req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", slice.FromIdx, slice.ToIdx))
			}
			resp, e := client.Do(req)
			if e != nil {
				errChan <- e
				return
			}
			// io.ReadCloser
			defer resp.Body.Close()
			// 将内存中的数据拷贝到目标文件中。需要写 & 操作权限。遵循最少权限原则。
			file, e := os.OpenFile(slice.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if e != nil {
				errChan <- e
				return
			}
			defer file.Close()
			// 每次拷贝100bytes
			offset := 0

			for {
				select {
				// 被中断，保存状态
				case <-interruptChan:
					st := &state.State{
						Url:     downloader.Url,
						FromIdx: slice.FromIdx,
						Offset:  offset,
						Path:    slice.Path,
					}
					stateChan <- st
					return
				default:
					n, e := io.CopyN(file, resp.Body, 100)
					if e != nil {
						if e == io.EOF {
							// fileChan <- slice.path
							log.Printf("download success : %v", slice.Path)
							return
						}
						errChan <- e
						return
					}
					offset = offset + int(n)
				}

			}

		}(downloader.DownloadSlice[i], downloader)
	}
	wg.Wait()
	downChan <- true
}
