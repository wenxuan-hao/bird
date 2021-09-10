package err

import (
	"log"
	"os"
)

func Handler(err error){
	log.Println(err)
	os.Exit(1)
}