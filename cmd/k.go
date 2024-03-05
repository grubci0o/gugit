package cmd

import (
	"fmt"
)

func K() {

	oids := make(map[string]struct{})
	for _, ref := range genRefs() {
		resolved := ResolveName(ref)
		fmt.Println("Ref: " + ref)
		fmt.Println("GetRef: " + resolved)
		if _, ok := oids[resolved]; !ok {
			oids[resolved] = struct{}{}
		}
	}

	keys := make([]string, len(oids))

	i := 0
	for k := range oids {
		keys[i] = k
		i++
	}
	fmt.Println("Hey")
	fmt.Println(keys)
	for _, oid := range commitTree(keys) {
		com := getCommit(oid)
		fmt.Println(oid)

		if len(com.Parent) > 0 {
			fmt.Printf("Commit parents: %v\n", com.Parent)
		}
	}
}

func commitTree(oids []string) []string {
	visited := make(map[string]bool)
	tr := make([]string, len(oids))
	var oid string
	copy(tr, oids)

	for len(tr) != 0 {
		oid, tr = tr[len(tr)-1], tr[:len(tr)-1]
		if oid == "" || visited[oid] == true {
			continue
		}
		visited[oid] = true

		com := getCommit(oid)
		tr = append(tr, com.Parent...)
	}
	return tr
}
