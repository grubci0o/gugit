package memory

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gugit/internal"
	"gugit/internal/util"
	"io"
	"log"
	"os"
	"strconv"
)

const UGIT_DIR = ".ugit"
const OBJECTS_DIR = ".ugit/objects"

func StoreObject(filePath string, objType internal.FileType) error {
	fmt.Println("Storage test: path is " + filePath)
	type_ := make([]byte, 4)
	binary.LittleEndian.PutUint32(type_, uint32(objType))
	OID := HashFile(filePath, objType)
	dest, err := os.Create(OBJECTS_DIR + "/" + OID)
	if err != nil {
		return &util.GuGitError{
			Type:    util.ErrIO,
			Message: "Failed to create object file",
			Err:     err,
		}
	}
	defer dest.Close()

	if _, err := dest.Write(type_); err != nil {
		return &util.GuGitError{
			Type:    util.ErrIO,
			Message: "Failed to write object type",
			Err:     err,
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return &util.GuGitError{
			Type:    util.ErrIO,
			Message: "Failed to open source file",
			Err:     err,
		}
	}
	defer f.Close()

	nBytes, err := io.Copy(dest, f)
	if err != nil {
		return &util.GuGitError{
			Type:    util.ErrIO,
			Message: "Failed to copy file contents",
			Err:     err,
		}
	}
	fmt.Println("Written this much bytes: " + strconv.Itoa(int(nBytes)))
	return nil
}

func StoreObjectBytes(data []byte, typ internal.FileType) string {
	typ_ := make([]byte, 4)
	binary.LittleEndian.PutUint32(typ_, uint32(typ))
	h := sha1.New()
	h.Write(typ_)
	h.Write(data)
	oid := hex.EncodeToString(h.Sum(nil))
	obj := OBJECTS_DIR + "\\" + oid

	f, err := os.Create(obj)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(typ_)
	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	return oid
}
