package cmd

import (
	"fmt"
	"gugit/internal/memory"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func Log(oid string) {
	// Get refs map first
	refs := make(map[string][]string)
	err := filepath.WalkDir(filepath.Join(memory.UGIT_DIR, "refs"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !d.IsDir() {
			relPath, err := filepath.Rel(memory.UGIT_DIR, path)
			if err != nil {
				return nil
			}
			_, refVal, _ := getRef(relPath, true)
			if refVal.value != "" {
				refs[refVal.value] = append(refs[refVal.value], relPath)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("Warning: error walking refs: %v\n", err)
	}

	// Print refs
	for oid, refNames := range refs {
		prettyRefs := make([]string, 0)
		for _, ref := range refNames {
			ref = strings.TrimPrefix(ref, "refs/heads/")
			ref = strings.TrimPrefix(ref, "refs/tags/")
			prettyRefs = append(prettyRefs, ref)
		}
		fmt.Printf("Commit %s, refs: [%s]\n", oid[:8], strings.Join(prettyRefs, ", "))
	}

	// Get starting commit
	if oid == "" {
		_, head, _ := getRef("HEAD", true)
		oid = head.value
	}

	fmt.Println("\nPrinting commit tree")

	// Track visited commits to avoid duplicates
	visited := make(map[string]bool)
	var commits []string
	commits = append(commits, oid)

	for len(commits) > 0 {
		current := commits[0]
		commits = commits[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		commit := getCommit(current)

		// Add parents to the queue
		if len(commit.Parent) > 0 {
			commits = append(commits, commit.Parent...)
		}

		fmt.Println("\t\t|")
		fmt.Printf("%v\n", current)
	}
}
