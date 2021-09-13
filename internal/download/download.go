package download

import (
	"bird/internal/err"
	"bird/internal/state"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type HttpDownloader struct {
	url           string
	Conc          int
	path          string
	downloadSlice []*downloadSlice
	resumeable    bool
	clen          int
	Interrupte    bool
}

type downloadSlice struct {
	url     string
	path    string
	fromIdx int
	toIdx   int
	num     int
}

var (
	client = http.Client{}
)

func NewHttpDownloader(url string, conc int, path string) *HttpDownloader {
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
	// 获取body长度
	len := resp.Header.Get("Content-Length")
	vlen, e := strconv.Atoi(len)
	if e != nil {
		err.Handler(e)
	}

	slices := make([]*downloadSlice, conc)

	for i := 0; i < conc; i++ {
		rg := vlen / conc
		fromidx := i * rg
		toidx := (i+1)*rg - 1
		if toidx > vlen {
			toidx = vlen
		}
		ds := &downloadSlice{
			url:     url,
			path:    filepath.Join(path, fmt.Sprintf("path %v", i)),
			fromIdx: fromidx,
			toIdx:   toidx,
			num:     i,
		}
		slices = append(slices, ds)
	}
	dl := &HttpDownloader{
		url:           url,
		Conc:          conc,
		path:          path,
		downloadSlice: slices,
		resumeable:    resume,
		clen:          vlen,
		Interrupte:    false,
	}
	return dl

}

func RecoverDownloader(state *state.State) *HttpDownloader {

}

func (downloader *HttpDownloader) Download(errChan chan error, interruptChan chan bool, stateChan chan *state.State, downChan chan bool) {
	wg := sync.WaitGroup{}
	for i := 0; i < downloader.Conc; i++ {
		wg.Add(1)
		go func(slice *downloadSlice, downloader *HttpDownloader) {
			defer wg.Done()
			req, e := http.NewRequest(http.MethodGet, slice.url, nil)
			if e != nil {
				errChan <- e
				return
			}
			if downloader.Conc > 1 {
				req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", slice.fromIdx, slice.toIdx))
			}
			resp, e := client.Do(req)
			if e != nil {
				errChan <- e
				return
			}
			// io.ReadCloser
			defer resp.Body.Close()
			// 将内存中的数据拷贝到目标文件中。需要写 & 操作权限。遵循最少权限原则。
			file, e := os.OpenFile(slice.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			defer file.Close()
			if e != nil {
				errChan <- e
				return
			}
			// 每次拷贝100bytes
			offset := 0

			for {
				select {
				// 被中断，保存状态
				case <-interruptChan:
					st := &state.State{
						Url:     downloader.url,
						FromIdx: slice.fromIdx,
						Offset:  offset,
						Path:    slice.path,
					}
					stateChan <- st
					break

				default:
					n, e := io.CopyN(file, resp.Body, 100)
					if e != nil {
						if e == io.EOF {
							break
						}
						errChan <- e
						return
					}
					offset = offset + int(n)
				}

			}

		}(downloader.downloadSlice[i], downloader)
	}
	wg.Wait()
	downChan <- true
}
