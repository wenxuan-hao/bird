package task

import (
	"bird/internal/download"
	"bird/internal/err"
	"bird/internal/state"
	"os"
	"os/signal"
	"syscall"
)

func StartDownload(url string, path string, conc int, st *state.State) {
	// 优雅关闭
	notifyChan := make(chan os.Signal, 1)
	signal.Notify(notifyChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var downloader *download.HttpDownloader
	if st == nil {
		downloader = download.NewHttpDownloader(url, conc, path)
	} else {
		downloader = download.RecoverDownloader(st)
	}

	// 错误同步channel
	errChan := make(chan error, 1)

	// 中断同步channel
	interruptChan := make(chan bool, downloader.Conc)

	// 状态同步channel
	stateChan := make(chan *state.State, downloader.Conc)
	states := make([]*state.State, downloader.Conc)

	// 结束channel
	downChan := make(chan bool, 1)

	downloader.Download(errChan, interruptChan, stateChan, downChan)

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
		case <-downChan:
			// 若被中断，保存各部分状态
			if downloader.Interrupte {

			} else {
				// 合并文件

				// success

			}
		}

	}

}
