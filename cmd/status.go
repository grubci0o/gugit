package cmd

import (
	"fmt"
)

func Status() {
	h := ResolveName("@")
	br := getBranch()
	if br != "" {
		fmt.Println("You are on branch " + br)
	} else {
		fmt.Println("Detached HEAD at " + h[:10])
	}

	_, mergeH := getRef("MERGE_HEAD", true)
	if mergeH.value != "" {
		fmt.Printf("Merging with %v", mergeH.value[:10])
	}
	fmt.Println("\nChanged files:")

	hTree := getCommit(h).Tree
	if hTree != "" {
		entries := diffFiles(GetTree(hTree, ""), getWorkingTree())
		for path, action := range entries {
			fmt.Println(action + ": " + path)
		}
	}
}
