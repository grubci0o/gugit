package cmd

import "gugit/internal"
import "gugit/internal/memory"

func FileCMD(filePath string, objType internal.FileType) {
	memory.StoreObject(filePath, objType)
}
