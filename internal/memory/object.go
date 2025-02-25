package memory

import (
	"encoding/binary"
	"gugit/internal"
	"gugit/internal/util"
	"os"
)

func GetObject(oid string) ([]byte, internal.FileType) {
	// probably need to handle ugit/HEAD as well ?
	data, err := os.ReadFile(OBJECTS_DIR + "/" + oid)
	util.Check(err)

	typ_ := data[:4]
	data = data[4:]
	return data, internal.FileType(binary.LittleEndian.Uint32(typ_))
}
