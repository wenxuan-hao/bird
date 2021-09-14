package cmd

import (
	"bird/internal/task"
	"errors"
	"log"
	"runtime"

	"github.com/spf13/cobra"
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

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().IntVarP(&cNum, "concurrency", "c", runtime.NumCPU(), "set the concurrency num")
}

func download(args []string) {
	// get url
	url := args[0]
	log.Printf("Strat downloading : %v \n", url)

	// task start
	task.StartDownload(url, cNum, nil)
}
