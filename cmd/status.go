package cmd

import (
	"fmt"
)

func Status() {
	h := ResolveName("@")
	br := getBranch()

	// Print branch status
	if br != "" {
		fmt.Printf("On branch %s\n", br)
	} else if h != "" {
		fmt.Printf("HEAD detached at %s\n", h[:8])
	} else {
		fmt.Println("No commits yet")
		fmt.Println("(use \"gugit commit\" to create the first commit)")
		return
	}

	// Print merge status if in merge
	_, mergeH, _ := getRef("MERGE_HEAD", true)
	if mergeH.value != "" {
		fmt.Printf("Merging with %s\n", mergeH.value[:8])
	}

	// Get trees for comparison
	var treeEntries map[string]string
	if h != "" {
		commit := getCommit(h)
		if commit.Tree != "" {
			treeEntries = GetTree(commit.Tree, "")
		}
	}
	if treeEntries == nil {
		treeEntries = make(map[string]string)
	}

	workingTree, err := getWorkingTree()
	if err != nil {
		fmt.Printf("Error reading working directory: %v\n", err)
		return
	}

	changes := diffFiles(treeEntries, workingTree)

	if len(changes) == 0 {
		fmt.Println("\nNothing to commit, working tree clean")
		return
	}

	fmt.Println("\nChanges:")
	printChangesByType(changes)
}

func printChangesByType(changes map[string]string) {
	var newFiles, modifiedFiles, deletedFiles []string
	for path, action := range changes {
		switch action {
		case "new file":
			newFiles = append(newFiles, path)
		case "modified":
			modifiedFiles = append(modifiedFiles, path)
		case "deleted":
			deletedFiles = append(deletedFiles, path)
		}
	}

	if len(newFiles) > 0 {
		fmt.Println("\nUntracked files:")
		for _, file := range newFiles {
			fmt.Printf("\t%s\n", file)
		}
		fmt.Println("\nnothing added to commit but untracked files present")
	}

	if len(modifiedFiles) > 0 {
		fmt.Println("\nModified files:")
		for _, file := range modifiedFiles {
			fmt.Printf("\t%s\n", file)
		}
	}

	if len(deletedFiles) > 0 {
		fmt.Println("\nDeleted files:")
		for _, file := range deletedFiles {
			fmt.Printf("\t%s\n", file)
		}
	}
}
