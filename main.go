package main

import (
	"log"

	"github.com/juliankoehn/wetter-service/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
