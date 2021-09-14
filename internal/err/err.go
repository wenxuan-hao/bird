package err

import (
	"log"
	"runtime"
)

func Handler(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Fatalf("ERROR [%s: %d]:%v \n", file, line, err)
	} else {
		log.Fatalln(err)
	}
}
