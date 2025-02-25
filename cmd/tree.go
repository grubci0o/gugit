package cmd

import (
	"bufio"
	"gugit/internal"
	"gugit/internal/memory"
	"gugit/internal/util"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WriteTree(dirPath string) string {
	dirEntries, err := os.ReadDir(dirPath)

	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	treeEntries := make([]internal.TreeEntry, 0)
	for _, entry := range dirEntries {
		var te internal.TreeEntry
		if util.Ignore(entry.Name()) {
			continue
		}
		if entry.IsDir() {
			typ_ := internal.TREE
			name := filepath.Join(dirPath, entry.Name())
			oid := WriteTree(name)
			te = internal.TreeEntry{
				Typ_: typ_,
				Name: name,
				Oid:  oid,
			}

		} else {
			typ_ := internal.BLOB
			name := filepath.Join(cwd, dirPath, entry.Name())
			oid := memory.HashFile(name, typ_)
			te = internal.TreeEntry{
				Typ_: typ_,
				Name: name,
				Oid:  oid,
			}
			memory.StoreObject(te.Name, te.Typ_)
		}
		treeEntries = append(treeEntries, te)
	}

	var sb strings.Builder
	for _, te := range treeEntries {
		sb.WriteString(" " + internal.TypeToString(uint32(te.Typ_)) + " " + te.Oid + " " + te.Name + "\n")
	}

	//skip empty directories
	ws := sb.String()
	if ws == "" {
		return ""
	}
	oid := memory.StoreObjectBytes([]byte(ws), internal.TREE)
	return oid
}

func GetTree(oid string, basePath string) map[string]string {
	// Handle empty oid
	if oid == "" {
		return make(map[string]string)
	}

	f, err := os.Open(memory.OBJECTS_DIR + "/" + oid)
	if err != nil {
		// Return empty map instead of fatal error
		return make(map[string]string)
	}
	defer f.Close()

	files := make(map[string]string)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var content []byte
	for scanner.Scan() {
		content = append(content, scanner.Bytes()...)
	}

	if len(content) < 4 {
		return files
	}

	content = content[4:]
	if len(content) == 0 {
		return files
	}

	entriesString := strings.Split(string(content), " ")
	if len(entriesString) <= 1 {
		return files
	}
	entriesString = entriesString[1:] // Skip first empty element after split

	for i := 0; i < len(entriesString); i += 3 {
		if i+2 >= len(entriesString) {
			break
		}

		typ_, hash, name := entriesString[i], entriesString[i+1], entriesString[i+2]

		cleanPath := filepath.ToSlash(filepath.Clean(name))

		wd, err := os.Getwd()
		if err == nil {
			if relPath, err := filepath.Rel(wd, cleanPath); err == nil {
				cleanPath = filepath.ToSlash(relPath)
			}
		}

		if typ_ == "BLOB" {
			files[cleanPath] = hash
		} else if typ_ == "TREE" {
			subFiles := GetTree(hash, cleanPath+"/")
			for k, v := range subFiles {
				files[filepath.ToSlash(filepath.Clean(k))] = v
			}
		}
	}
	return files
}

func ReadTree(oid string) {
	util.CleanDir()

	for path, oid := range GetTree(oid, "./") {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		data, typ := memory.GetObject(oid)
		if typ != internal.BLOB {
			log.Fatal("Expected blob, got " + internal.TypeToString(uint32(typ)))
		}

		// Write file contents (skip the type bytes)
		if _, err := f.Write(data); err != nil {
			log.Fatal(err)
		}
	}
}

func getWorkingTree() (map[string]string, error) {
	entries := make(map[string]string)
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Convert to forward slashes and clean path
		cleanPath := filepath.ToSlash(filepath.Clean(path))

		// Skip directories and ignored paths
		if d.IsDir() || util.Ignore(cleanPath) {
			return nil
		}

		// Make path relative to current directory
		relPath, err := filepath.Rel(".", cleanPath)
		if err != nil {
			return err
		}

		// Normalize to forward slashes
		normalizedPath := filepath.ToSlash(relPath)

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		entries[normalizedPath] = memory.HashObject(data, internal.BLOB)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}
