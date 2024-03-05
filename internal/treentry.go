package internal

type TreeEntry struct {
	Typ_ FileType //blob or tree
	Name string   //filename
	Oid  string
}
