package main

import (
	"bird/internal/task"
	"runtime"
)

func main() {
	//cmd.Execute()
	src := "https://lf1-ttcdn-tos.pstatp.com/obj/mubu-assets/client/Mubu-3.6.0.dmg"
	task.StartDownload(src, runtime.NumCPU(), nil)

}
