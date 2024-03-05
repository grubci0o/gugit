package cmd

import (
	"fmt"
	"gugit/internal/memory"
	"log"
)

func createTag(name, oid string) {
	updateRef("refs/tags/"+name, RefValue{
		symbolic: false,
		value:    oid,
	}, true)
}

func Tag(name, oid string) {
	if oid == "" {
		_, val := getRef("HEAD", true)
		createTag(name, val.value)
	} else {
		createTag(name, oid)
	}
}

func ResolveName(name string) string {
	if name == "@" {
		_, rVal := getRef("HEAD", true)
		return rVal.value
	}
	lookDirs := []string{name, "refs/" + name, "refs/tags/" + name, "refs/heads/" + name}
	for _, dir := range lookDirs {
		if _, ref := getRef(dir, false); ref.value != "" {
			fmt.Println("Object with that tag: " + ref.value)
			_, rVal := getRef(dir, true)
			return rVal.value
		}
	}
	if memory.ValidOID(name) {
		return name
	}
	log.Fatalln("Unrecognized name. Not a tag or OID")
	//unreachable
	return ""
}
