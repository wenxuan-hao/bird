package cmd

import (
	"bird/internal/err"
	"bird/internal/task"
	"bird/internal/tool"
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download file form a url",
	Run: func(cmd *cobra.Command, args []string) {
		download(args)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		return errors.New("wrong with download args, should only have one url")
	},

}

var cNum int

func init(){
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().IntVarP(&cNum, "concurrency", "c", runtime.NumCPU(), "set the concurrency num")
}

func download(args []string){
	// get url
	url := args[0]
	log.Printf("Strat downloading : %v \n", url)
	path := tool.GetFolder(url)
	if tool.IsDirExit(path){
		os.RemoveAll(path)
	}
	e := os.MkdirAll(path, os.ModeDir)
	if e != nil{
		err.Handler(e)
	}
	// task start
	task.StartDownload(url, path, nil)
}
