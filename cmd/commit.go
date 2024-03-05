package cmd

import (
	"gugit/internal"
	"gugit/internal/memory"
	"log"
	"os"
	"strings"
	"time"
)

func Commit() {
	c := CreateCommit("My message")
	cs := c.String()
	_, head := getRef("HEAD", true)
	h := head.value

	_, mergeH := getRef("MERGE_HEAD", true)

	if mergeH.value != "" {
		cs += "\nParent: " + mergeH.value
	}

	if len(h) > 0 {
		cs += "\nparent " + h + "\n"
	}
	oid := memory.HashObject([]byte(cs), internal.COMMIT)

	dest, err := os.Create(memory.OBJECTS_DIR + "/" + oid)
	if err != nil {
		log.Fatal(err)
	}
	defer dest.Close()

	_, err = dest.Write(internal.COMMIT.ToBytes())
	_, err = dest.Write([]byte(cs))
	if err != nil {
		log.Fatal(err)
	}

	updateRef("HEAD", RefValue{symbolic: false, value: oid}, true)
}

func CreateCommit(msg string) internal.Commit {

	oid := WriteTree(".")
	author := "Test"
	t := time.Now()
	c := internal.Commit{Oid: oid, Author: author, Time: t, Msg: msg}
	return c
}

func getCommit(oid string) internal.Commit {
	com, typ_ := memory.GetObject(oid)
	if internal.TypeToString(uint32(typ_)) != "COMMIT" {
		log.Fatal("Returned file is not of commit type")
	}

	content := string(com)

	var tree string
	var parent []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		l := strings.Split(line, " ")
		if len(l) < 2 {
			// empty line before commit message - we don't want to parse that
			continue
		}

		k, v := l[0], l[1]
		if k == "TREE" {
			tree = v
		} else if k == "parent" {
			parent = append(parent, v)
		}
	}

	msg := lines[len(lines)-1]
	//skipping parsing time FOR NOW
	return internal.Commit{Oid: tree, Parent: parent, Author: "test", Time: time.Now(), Msg: msg, Tree: tree}
}
