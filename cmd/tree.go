package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"gugit/internal"
	"gugit/internal/memory"
	"gugit/internal/util"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WriteTree(dirPath string) string {
	dirEntries, err := os.ReadDir(dirPath)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	treeEntries := make([]internal.TreeEntry, 0)
	for _, entry := range dirEntries {
		var te internal.TreeEntry
		if util.Ignore(entry.Name()) {
			continue
		}
		if entry.IsDir() {
			typ_ := internal.TREE
			name := dirPath + "\\" + entry.Name()
			oid := WriteTree(name)
			te = internal.TreeEntry{
				Typ_: typ_,
				Name: name,
				Oid:  oid,
			}

		} else {
			typ_ := internal.BLOB
			name := cwd + "\\" + dirPath + "\\" + entry.Name()
			oid := memory.HashFile(name, typ_)
			te = internal.TreeEntry{
				Typ_: typ_,
				Name: name,
				Oid:  oid,
			}
			memory.StoreObject(te.Name, te.Typ_)
		}
		treeEntries = append(treeEntries, te)
	}

	var sb strings.Builder
	for _, te := range treeEntries {
		sb.WriteString(" " + internal.TypeToString(uint32(te.Typ_)) + " " + te.Oid + " " + te.Name + "\n")
	}
	//skip empty directories
	ws := sb.String()
	oid := memory.StoreObjectBytes([]byte(ws), internal.TREE)
	fmt.Println("Tree oid " + oid)
	return oid
}

func GetTree(oid string, basePath string) map[string]string {
	f, err := os.Open(memory.OBJECTS_DIR + "/" + oid)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//path to oid
	files := make(map[string]string)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	lines := make([]byte, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Bytes()...)
	}
	//skip this file type
	lines = lines[4:]
	entriesString := strings.Split(string(lines), " ")[1:]

	for i := 0; i < len(entriesString); i = i + 3 {
		typ_, hash, name := entriesString[i], entriesString[i+1], entriesString[i+2]
		wd, err := os.Getwd()
		util.Check(err)
		path, err := filepath.Rel(wd, name)
		util.Check(err)
		if typ_ == "BLOB" {
			files[path] = hash
		} else if typ_ == "TREE" {
			for k, v := range GetTree(hash, path+"/") {
				//k is path so add /dir/ to path to it won't create just a file
				files[k] = v
			}
		}
	}
	return files
}

func ReadTree(oid string) {
	//Not sure about using this yet
	//cleanDir()
	for path, oid := range GetTree(oid, "./") {
		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		data, err := os.ReadFile(memory.OBJECTS_DIR + "/" + oid)

		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(f, bytes.NewReader(data[4:]))
		if err != nil {
			log.Fatal(err)
		}
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getWorkingTree() map[string]string {
	entries := make(map[string]string)
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if util.Ignore(path) || d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)

		if err != nil {
			return err
		}
		entries[path] = memory.HashObject(data, internal.BLOB)
		return nil
	})
	util.Check(err)
	return entries
}
