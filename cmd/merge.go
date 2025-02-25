package cmd

import (
	"fmt"
	"gugit/internal"
	"gugit/internal/memory"
	"gugit/internal/util"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Merge(otherBranch string) error {
	// Get current HEAD commit
	_, head, err := getRef("HEAD", true)
	if err != nil {
		return &util.GuGitError{
			Type:    util.ErrMerge,
			Message: "Failed to get HEAD reference",
			Err:     err,
		}
	}

	// Get other branch commit
	_, otherBranchRef, err := getRef(otherBranch, true)
	if err != nil {
		return &util.GuGitError{
			Type:    util.ErrMerge,
			Message: "Branch " + otherBranch + " not found",
			Err:     err,
		}
	}

	// Get common ancestor
	ancestor := getCommonAncestor(head.value, otherBranchRef.value)
	if ancestor == "" {
		log.Println("Warning: no common ancestor found, doing recursive merge")
	}

	// Store MERGE_HEAD for the commit
	if err := updateRef("MERGE_HEAD", RefValue{
		symbolic: false,
		value:    otherBranchRef.value,
	}, true); err != nil {
		return &util.GuGitError{
			Type:    util.ErrMerge,
			Message: "Failed to update MERGE_HEAD",
			Err:     err,
		}
	}

	// Read and merge trees
	comHead := getCommit(head.value)
	comOther := getCommit(otherBranchRef.value)

	// Print merge information
	fmt.Printf("Merging branch '%s' into '%s'\n", otherBranch, getBranch())

	changes := readTreeMerged(comHead.Tree, comOther.Tree)
	if len(changes) > 0 {
		fmt.Println("\nChanged files:")
		for path, action := range changes {
			fmt.Printf("%s: %s\n", action, path)
		}
	}

	// Create merge commit
	mergeCommit("Merge branch '" + otherBranch + "'")

	// Properly close and remove MERGE_HEAD
	if err := os.Remove(filepath.Join(memory.UGIT_DIR, "MERGE_HEAD")); err != nil {
		// If file is in use, schedule removal for program exit
		defer func() {
			_ = os.Remove(filepath.Join(memory.UGIT_DIR, "MERGE_HEAD"))
		}()
	}

	fmt.Printf("\nMerge completed successfully\n")
	return nil
}

func readTreeMerged(treeHead, treeOther string) map[string]string {
	changes := make(map[string]string)
	mergedTree := mergeTrees(GetTree(treeHead, ""), GetTree(treeOther, ""))

	for path, blob := range mergedTree {
		fullPath := filepath.Join(".", path)

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			log.Printf("Warning: couldn't create directory for %s: %v\n", path, err)
			continue
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(blob), 0644); err != nil {
			log.Printf("Warning: couldn't write file %s: %v\n", path, err)
			continue
		}

		changes[path] = "modified"
	}

	return changes
}

func mergeTrees(treeHead, treeOther map[string]string) map[string]string {
	tree := make(map[string]string)
	entries := compareTrees([]map[string]string{treeHead, treeOther})

	for path, oids := range entries {
		tree[path] = mergeBlobs(oids[0], oids[1])
	}

	return tree
}

func mergeBlobs(objHead, objOther string) string {
	var objH, objO []byte
	if objHead != "" {
		objH, _ = memory.GetObject(objHead)
	}
	if objOther != "" {
		objO, _ = memory.GetObject(objOther)
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(objH), string(objO), true)
	diffs = dmp.DiffCleanupMerge(diffs)

	// Print diff information
	if len(diffs) > 1 { // More than one diff means there are changes
		fmt.Printf("\nConflict resolution for changes:\n")
		fmt.Println(dmp.DiffPrettyText(diffs))
	}

	merged := dmp.DiffText2(diffs)
	return merged
}

func getCommonAncestor(oid1, oid2 string) string {
	//build a visited "set" of ancestors from first oid
	comOIDS := []string{oid1}
	visited := make(map[string]bool)
	for len(comOIDS) > 0 {
		o := comOIDS[0]
		if _, ok := visited[o]; !ok {
			visited[o] = true
		}
		comOIDS = comOIDS[1:]
		com := getCommit(o)
		parents := com.Parent
		if len(parents) > 0 {
			comOIDS = append([]string{parents[0]}, comOIDS...)
			comOIDS = append(comOIDS, parents[1:]...)
		}

	}

	//iterate over oid2 ancestors
	//if oid of ancestor is in visited set its LCA
	//else return empty string
	comOIDS2 := []string{oid2}

	for len(comOIDS2) > 0 {
		o := comOIDS2[0]
		if _, ok := visited[o]; ok {
			return o
		}
		comOIDS2 = comOIDS2[1:]
		com := getCommit(o)
		parents := com.Parent
		if len(parents) > 0 {
			comOIDS2 = append([]string{parents[0]}, comOIDS2...)
			comOIDS2 = append(comOIDS2, parents[1:]...)
		}

	}
	return ""
}

func mergeCommit(message string) {
	// Get current working directory state
	oid := WriteTree(".")

	// Get HEAD and MERGE_HEAD
	_, head, err := getRef("HEAD", true)
	_, mergeHead, err := getRef("MERGE_HEAD", true)

	// Create commit object
	author := "Test"
	t := time.Now()
	c := internal.Commit{
		Oid:    oid,
		Author: author,
		Time:   t,
		Msg:    message,
		Parent: []string{head.value, mergeHead.value}, // Both parents
		Tree:   oid,
	}

	// Write commit object
	cs := c.String()
	cs += "\nparent " + head.value + "\n"
	cs += "parent " + mergeHead.value + "\n"

	commitOid := memory.HashObject([]byte(cs), internal.COMMIT)

	dest, err := os.Create(memory.OBJECTS_DIR + "/" + commitOid)
	if err != nil {
		log.Fatal(err)
	}
	defer dest.Close()

	_, err = dest.Write(internal.COMMIT.ToBytes())
	_, err = dest.Write([]byte(cs))
	if err != nil {
		log.Fatal(err)
	}

	// Update HEAD to point to new commit
	updateRef("HEAD", RefValue{
		symbolic: false,
		value:    commitOid,
	}, true)

	// In mergeCommit function:
	if err := deleteRef("MERGE_HEAD", false); err != nil {
		log.Printf("Warning: could not delete MERGE_HEAD: %v\n", err)
	}
}
