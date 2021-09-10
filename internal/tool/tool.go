package tool

import (
	"bird/internal/err"
	"os"
	"path/filepath"
)

const BIRD_HOME  = "BIRD_HOME"

func GetFolder(url string)  string {
	base := filepath.Base(url)
	home := os.Getenv(BIRD_HOME)
	path, e := filepath.Abs(filepath.Join(home, base))

	if e != nil{
		err.Handler(e)
	}
	return path
	
}

func IsDirExit(path string) bool{
	_ , e := os.Stat(path)
	return e != nil
}
