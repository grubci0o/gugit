package util

import (
	"log"
	"os"
)

func CleanDir() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := os.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		name := entry.Name()
		// Skip .ugit directory and other ignored files
		if Ignore(name) {
			continue
		}

		if err := os.RemoveAll(name); err != nil {
			log.Printf("Warning: couldn't delete %s: %v\n", name, err)
		}
	}
}
