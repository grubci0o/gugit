package cmd

import (
	"encoding/binary"
	"fmt"
	"gugit/internal"
	"gugit/internal/memory"
	"io"
	"log"
	"os"
)

func catFile(filePath string) {
	f, err := os.Open(memory.OBJECTS_DIR + "/" + filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fileContent, err := io.ReadAll(f)
	fileType := binary.LittleEndian.Uint32(fileContent[:4])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(internal.TypeToString(fileType))
	fmt.Println(string(fileContent[4:]))
}

func CatCMD(filePath string) {
	catFile(filePath)
}
