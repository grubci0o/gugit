package cmd

import (
	"gugit/internal/memory"
	"log"
	"os"
)

func Init() {
	err := os.Mkdir(memory.UGIT_DIR, 0700)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir(memory.OBJECTS_DIR, 0700)
	if err != nil {
		log.Fatal(err)
	}
	println("Initialized ugit directory")

	updateRef("HEAD", RefValue{symbolic: true, value: "refs/heads/master"}, true)
}
