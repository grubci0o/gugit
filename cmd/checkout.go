package cmd

import "fmt"

func Checkout(name string) {
	oid := ResolveName(name)
	com := getCommit(oid)
	fmt.Println("At checkout and ill read tree with oid " + com.Tree)
	ReadTree(com.Tree)

	var h RefValue
	fmt.Println("Is " + name + "a branch")
	if isBranch(name) {
		fmt.Println("It's a branch")
		h = RefValue{
			symbolic: true,
			value:    "refs/heads/" + name,
		}
	} else {
		h = RefValue{
			symbolic: false,
			value:    oid,
		}
	}

	updateRef("HEAD", h, false)
}
