// +build ignore

package main

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/diamondburned/cchat/repository"
)

func main() {
	f, err := os.Create("repository.gob")
	if err != nil {
		log.Fatalln("Failed to create file:", err)
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(repository.Main); err != nil {
		os.Remove("repository.gob")
		log.Fatal("Failed to gob encode:", err)
	}
}
