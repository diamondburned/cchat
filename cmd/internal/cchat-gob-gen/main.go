package main

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/diamondburned/cchat/repository"
)

const output = "repository.gob"

func main() {
	f, err := os.Create(output)
	if err != nil {
		log.Fatalln("Failed to create file:", err)
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(repository.Main); err != nil {
		os.Remove(output)
		log.Fatalln("Failed to gob encode:", err)
	}
}
