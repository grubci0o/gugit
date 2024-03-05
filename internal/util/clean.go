package util

import (
	"log"
	"os"
)

func cleanDir() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir, err := os.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dir {
		if d.Name() == ".ugit" {
			continue
		}
		err := os.RemoveAll(d.Name())
		if err != nil {
			log.Println("Couldn't delete file " + d.Name())
		}
	}
}
