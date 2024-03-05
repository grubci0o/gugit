package cmd

import (
	"fmt"
	"gugit/internal/memory"
	"gugit/internal/util"
	"io/fs"
	"path/filepath"
)

func Log(oid string) {
	refs := make(map[string][]string)
	err := filepath.WalkDir(memory.UGIT_DIR+"/"+"refs", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			base, err := filepath.Rel(memory.UGIT_DIR, path)
			util.Check(err)
			ref, refVal := getRef(base, true)
			ref = d.Name()
			refs[refVal.value] = append(refs[refVal.value], ref)
		}
		return err
	})

	for k, v := range refs {
		fmt.Printf("Commit OID %s, refs pointing to it -> %v\n", k, v)
	}

	util.Check(err)
	if oid == "" || oid == " " {
		_, rVal := getRef("HEAD", true)
		oid = rVal.value
	}

	//needed for multiple parents
	var comOIDS []string
	comOIDS = append(comOIDS, oid)

	fmt.Println("\nPrinting commit tree")
	for len(comOIDS) > 0 {
		//pop from left
		o := comOIDS[0]
		comOIDS = comOIDS[1:]
		com := getCommit(o)
		parents := com.Parent
		//prepend first commit
		//append rest
		if len(parents) > 0 {
			comOIDS = append([]string{parents[0]}, comOIDS...)
			comOIDS = append(comOIDS, parents[1:]...)
		}
		fmt.Println("\t\t|")
		fmt.Printf("%v\n", o)
	}
}
