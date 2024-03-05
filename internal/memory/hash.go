package memory

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"gugit/internal"
	"io"
	"log"
	"os"
	"regexp"
)

func HashObject(data []byte, typ_ internal.FileType) string {
	h := sha1.New()
	h.Write(typ_.ToBytes())
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func HashFile(filePath string, typ_ internal.FileType) string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	h.Write(binary.LittleEndian.AppendUint32(make([]byte, 0), uint32(typ_)))
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	hashString := hex.EncodeToString(h.Sum(nil))
	return hashString
}

func ValidOID(oid string) bool {
	r, _ := regexp.Compile("^[a-fA-F0-9]{40}$")
	return r.MatchString(oid)
}
