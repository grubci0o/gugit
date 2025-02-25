package cmd

import (
	"fmt"
	"gugit/internal"
	"gugit/internal/memory"
	"gugit/internal/util"
	"log"
	"os"
	"strings"
	"time"
)

func Commit(dirPath string) error {
	if dirPath == "" {
		dirPath = "." // Default to current directory if none specified
	}

	c := CreateCommit("My message", dirPath)
	cs := c.String()
	_, head, err := getRef("HEAD", true)
	if err != nil {
		return &util.GuGitError{
			Type:    util.ErrCommit,
			Message: "Failed to get HEAD reference",
			Err:     err,
		}
	}
	h := head.value

	_, mergeH, _ := getRef("MERGE_HEAD", true)
	if mergeH.value != "" {
		cs += "\nParent: " + mergeH.value
	}

	if len(h) > 0 {
		cs += "\nparent " + h + "\n"
	}
	oid := memory.HashObject([]byte(cs), internal.COMMIT)

	// Create commit object
	dest, err := os.Create(memory.OBJECTS_DIR + "/" + oid)
	if err != nil {
		return &util.GuGitError{
			Type:    util.ErrCommit,
			Message: "Failed to create commit object",
			Err:     err,
		}
	}
	defer dest.Close()

	_, err = dest.Write(internal.COMMIT.ToBytes())
	_, err = dest.Write([]byte(cs))
	if err != nil {
		log.Fatal(err)
	}

	updateRef("HEAD", RefValue{symbolic: false, value: oid}, true)
	fmt.Printf("\nCommitted with ID: %s\n", oid[:8])
	return nil
}

func CreateCommit(msg string, dirPath string) internal.Commit {
	oid := WriteTree(dirPath)
	author := "Test"
	t := time.Now()
	c := internal.Commit{Oid: oid, Author: author, Time: t, Msg: msg}
	return c
}

func getCommit(oid string) internal.Commit {
	com, typ_ := memory.GetObject(oid)
	if internal.TypeToString(uint32(typ_)) != "COMMIT" {
		log.Fatal("Returned file is not of commit type " + internal.TypeToString(uint32(typ_)))
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
