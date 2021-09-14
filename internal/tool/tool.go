package tool

import (
	"bird/internal/err"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
)

const BIRD_HOME = "BIRD_HOME"

// 缓存下载分片文件的文件夹。默认为 ~/.bird/文件名
func GetCacheFolder(url string) string {
	base := filepath.Base(url)
	home := os.Getenv(BIRD_HOME)
	if home == "" {
		var e error
		home, e = GetHome()
		if e != nil {
			err.Handler(e)
		}
		home = filepath.Join(home, ".bird")
	}
	path, e := filepath.Abs(filepath.Join(home, base))
	if e != nil {
		err.Handler(e)
	}
	if IsDirExit(path) {
		os.RemoveAll(path)
	}
	e = os.MkdirAll(path, 0755)
	if e != nil {
		err.Handler(e)
	}

	return path
}

func GetHome() (string, error) {
	u, e := user.Current()
	if e != nil {
		return "", e
	}
	return u.HomeDir, nil
}

func IsDirExit(path string) bool {
	_, e := os.Stat(path)
	return e != nil
}

func MergeFile(files []string, dest string) {
	log.Printf("merge file: %v \n", dest)
	// 文件排序
	sort.Strings(files)

	f, e := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if e != nil {
		err.Handler(e)
	}
	defer f.Close()

	// 合并为目标文件
	for _, file := range files {
		src, e := os.Open(file)
		if e != nil {
			err.Handler(e)
		}
		defer src.Close()
		io.Copy(f, src)
	}
}
