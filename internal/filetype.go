package internal

import (
	"bytes"
	"encoding/binary"
	"log"
)

const (
	BLOB FileType = iota
	TREE
	COMMIT
)

type FileType uint32

func (f FileType) ToBytes() []byte {
	u := uint32(f)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, u)
	if err != nil {
		log.Println("Could not convert FileType to byte array")
	}
	return buf.Bytes()
}

func TypeToString(typ uint32) string {
	switch typ {
	case 0:
		return "BLOB"
	case 1:
		return "TREE"
	case 2:
		return "COMMIT"

	}
	return ""
}
