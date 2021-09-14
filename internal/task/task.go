package task

import (
	"bird/internal/download"
	"bird/internal/err"
	"bird/internal/state"
	"bird/internal/tool"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func StartDownload(url string, conc int, st *state.State) {
	// 优雅关闭
	notifyChan := make(chan os.Signal, 1)
	signal.Notify(notifyChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var downloader *download.HttpDownloader
	if st == nil {
		downloader = download.NewHttpDownloader(url, conc)
	} else {
		// downloader = download.RecoverDownloader(st)
	}

	// 错误同步channel
	errChan := make(chan error, 1)

	// 中断同步channel
	interruptChan := make(chan bool, downloader.Conc)

	// 状态同步channel
	stateChan := make(chan *state.State, downloader.Conc)
	states := []*state.State{}

	// 结束channel
	downChan := make(chan bool, 1)

	// 分片文件处理完成channel
	fileChan := make(chan string, downloader.Conc)
	files := []string{}

	downloader.Download(errChan, interruptChan, stateChan, downChan, fileChan)

	for {
		select {
		case <-notifyChan:
			for i := 0; i < downloader.Conc; i++ {
				interruptChan <- true
			}
			downloader.Interrupte = true
		case e := <-errChan:
			err.Handler(e)
		case item := <-stateChan:
			states = append(states, item)
		// case item := <-fileChan:
		// 	files = append(files, item)
		case <-downChan:
			// 若被中断，保存各部分状态
			if downloader.Interrupte {
				log.Println("Be interrupteed, start saving state.")
				state.Save(states)
				return
			} else {
				// 合并文件
				for _, slice := range downloader.DownloadSlice {
					files = append(files, slice.Path)
				}
				tool.MergeFile(files, filepath.Base(url))
				// success
				log.Println("DOWNLOAD SUCCESS!")
				return
			}
		}

	}

}
