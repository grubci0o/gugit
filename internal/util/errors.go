package util

import (
	"fmt"
	"log"
	"os"
)

type ErrType int

const (
	ErrIO ErrType = iota
	ErrRef
	ErrCommit
	ErrMerge
	ErrBranch
)

type GuGitError struct {
	Type    ErrType
	Message string
	Err     error
}

func (e *GuGitError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func HandleError(err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *GuGitError:
		switch e.Type {
		case ErrIO:
			log.Printf("IO Error: %v\n", e)
		case ErrRef:
			log.Printf("Reference Error: %v\n", e)
		case ErrCommit:
			log.Printf("Commit Error: %v\n", e)
		case ErrMerge:
			log.Printf("Merge Error: %v\n", e)
		case ErrBranch:
			log.Printf("Branch Error: %v\n", e)
		}
		os.Exit(1)
	default:
		log.Printf("Unexpected error: %v\n", err)
		os.Exit(1)
	}
}
