package cmd

import (
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"gugit/internal/memory"
	"gugit/internal/util"
)

func Show(oid string) {
	com := getCommit(oid)
	fmt.Println(com.String())
	parent := com.Parent
	var parTree string
	if len(parent) > 0 {
		parTree = getCommit(parent[0]).Tree
	}
	result := diffTrees(GetTree(parTree, ""), GetTree(com.Tree, ""))
	fmt.Println(result)
}

func compareTrees(trees []map[string]string) map[string][]string {
	//entries are map path:list of oids
	//basically later for each file we can its ids and if they are different
	entries := make(map[string][]string, len(trees))

	for i, tree := range trees {
		for path, oid := range tree {
			if len(entries[path]) == 0 {
				entries[path] = make([]string, len(trees))
			}
			entries[path][i] = oid
		}
	}
	return entries
}

func diffTrees(treeFrom, treeTo map[string]string) string {
	var output string
	entries := compareTrees([]map[string]string{treeFrom, treeTo})
	for path, oids := range entries {
		oidFrom, oidTo := oids[0], oids[1]
		if oidFrom != oidTo {
			out, err := difflib.GetUnifiedDiffString(diffBlobs(oidFrom, oidTo, path))
			util.Check(err)
			output += out
		}
	}
	return output
}

func diffBlobs(oidFrom, oidTo, path string) difflib.UnifiedDiff {
	blobFrom, _ := memory.GetObject(oidFrom)
	blobTo, _ := memory.GetObject(oidTo)
	unidiff := difflib.UnifiedDiff{A: difflib.SplitLines(string(blobFrom)),
		B:        difflib.SplitLines(string(blobTo)),
		FromFile: oidFrom,
		ToFile:   oidTo,
		Context:  3}
	return unidiff
}

func diffFiles(treeFrom, treeTo map[string]string) map[string]string {
	entries := compareTrees([]map[string]string{treeFrom, treeTo})

	changes := make(map[string]string, len(entries))
	for path, oids := range entries {
		objFrom, objTo := oids[0], oids[1]
		if objFrom != objTo {
			var action string
			if objFrom == "" {
				action = "new file"
			} else if objTo == "" {
				action = "deleted"
			} else {
				action = "modified"
			}
			changes[path] = action
		}
	}
	return changes
}
