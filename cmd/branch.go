package cmd

import (
	"fmt"
	"gugit/internal/memory"
	"gugit/internal/util"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func createBranch(name, oid string) {
	updateRef("refs/heads/"+name, RefValue{symbolic: false, value: oid}, true)
}

func Branch(name, oid string) {
	createBranch(name, oid)
	fmt.Println("Created branch " + name + "created from oid " + oid)
}

func isBranch(branch string) bool {
	// Remove any refs/heads/ prefix if it exists
	branch = strings.TrimPrefix(branch, "refs/heads/")

	// Check if the branch reference exists
	_, refVal, _ := getRef("refs/heads/"+branch, true)
	return refVal.value != ""
}

func getBranch() string {
	_, h, err := getRef("HEAD", false)
	if err != nil {
		return ""
	}
	if !h.symbolic {
		return ""
	}
	v := h.value
	if !strings.HasPrefix(v, "refs/heads/") {
		fmt.Println(v)
		log.Fatalln("Not a branch or file has been corrupted.")
	}
	p, err := filepath.Rel("refs/heads", v)
	util.Check(err)
	return p
}

func ListBranches() {
	var branches []string
	err := filepath.WalkDir(memory.UGIT_DIR+"/refs/heads/", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			branches = append(branches, d.Name())
		}
		return err
	})
	util.Check(err)

	headBranch := getBranch()
	for _, br := range branches {
		if br == headBranch {
			fmt.Printf("  *%s\n", br)
		}
		fmt.Printf("  %s\n", br)
	}
}
