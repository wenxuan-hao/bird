package task

import (
	"bird/internal/state"
	"os"
	"os/signal"
	"syscall"
)

func StartDownload(url string, path string, state *state.State){
	// 优雅关闭
	notifyChan := make(chan os.Signal, 1)
	signal.Notify(notifyChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for{
		select {
		case <-notifyChan:

		}
	}


}
