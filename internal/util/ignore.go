package util

const UGIT_DIR = ".ugit"
const OBJECTS_DIR = ".ugit/objects"

func Ignore(entry string) bool {
	if entry == UGIT_DIR || entry == OBJECTS_DIR {
		return true
	}
	return false
}
