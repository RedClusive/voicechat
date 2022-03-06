package server

import (
	"log"
	"os"
)

type serverConfig struct {
	Port uint
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("not enough arguments")
		return
	}
}
