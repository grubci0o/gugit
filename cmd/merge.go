package cmd

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gugit/internal/memory"
	"gugit/internal/util"
	"log"
)

func Merge(otherBranch string) {
	_, r := getRef("HEAD", true)
	if r.value == "" {
		log.Fatalln("Could not merge -> HEAD not found.")
	}

	_, otherBranchRef := getRef(otherBranch, true)
	comHead := getCommit(r.value)
	comOther := getCommit(otherBranchRef.value)

	updateRef("MERGE_HEAD", RefValue{symbolic: false, value: otherBranch}, true)
	readTreeMerged(comHead.Oid, comOther.Oid)
	fmt.Println("Merged in working tree\nPleaseCommit")
}

func readTreeMerged(treeHead, treeOther string) {
	for path, blob := range mergeTrees(GetTree(treeHead, ""), GetTree(treeOther, "")) {
		f, err := create(memory.UGIT_DIR + "/" + path)
		util.Check(err)
		_, err = f.Write([]byte(blob))
		util.Check(err)
	}
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

	return dmp.DiffPrettyText(diffs)
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
