package cmd

import (
	"fmt"
	"gugit/internal/memory"
	"gugit/internal/util"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RefValue struct {
	symbolic bool
	value    string
}

func updateRef(ref string, value RefValue, deref bool) {
	r, _ := getRef(ref, deref)

	//save symbolic ref
	if value.value == "" {
		return
	}
	var v string
	if value.symbolic {
		v = "ref: " + value.value
	} else {
		v = value.value
	}
	f, err := create(memory.UGIT_DIR + "/" + r)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(v))
	if err != nil {
		log.Println("Could not save head")
		log.Fatal(err)
	}
}

func getRef(oid string, deref bool) (string, RefValue) {
	path := memory.UGIT_DIR + "/" + oid
	if refNotExist(oid) {
		fmt.Println("Ref " + oid + " Does not exist ")
		return oid, RefValue{}
	}
	b, err := os.ReadFile(path)
	util.Check(err)
	s := string(b)

	//handle symbolic refs -> ref: <refname>
	split := strings.Split(s, ": ")
	split = deleteEmpty(split)

	var isSymbolic bool
	if strings.HasPrefix(s, "ref:") && len(split) == 2 {
		isSymbolic = true
	}

	if isSymbolic {
		if deref {
			//symbolic and recursively dereference it

			return getRef(split[1], true)
		}
		//case its symbolic, but we don't want to deref it
		return oid, RefValue{isSymbolic, split[1]}
	}
	//not symbolic so its just OID in file
	return oid, RefValue{symbolic: isSymbolic, value: split[0]}
}

func refNotExist(oid string) bool {
	_, err := os.Open(memory.UGIT_DIR + "/" + oid)
	return os.IsNotExist(err)
}

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func genRefs() []string {
	refs := []string{"HEAD"}
	err := filepath.WalkDir(memory.UGIT_DIR+"/"+"refs", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			refs = append(refs, d.Name())
		}
		return nil
	})
	util.Check(err)
	return refs
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func deleteRef(ref string, deref bool) {
	r, _ := getRef(ref, deref)
	err := os.Remove(memory.UGIT_DIR + "/" + r)
	util.Check(err)
}
