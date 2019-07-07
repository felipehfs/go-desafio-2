package models

import (
	"fmt"
	"log"
	"time"
)

type FileRegistry struct {
	ID       int
	Filename string
	Path     string
	ModTime  time.Time
	CheckSum []byte
	Size     int64
	UUID     string
}

func (f *FileRegistry) hasChanges(modTime time.Time, filename string, path string, size int64) bool {
	if f.Filename != filename {
		log.Println(fmt.Sprintf("[%s filename changed]: old=%s new=%s", filename, f.Filename, filename))
		return true
	}
	if f.Path != path {
		log.Println(fmt.Sprintf("[%s path changed]: old=%s new=%s", filename, f.Path, path))
		return true
	}

	if f.Size != size {
		log.Println(fmt.Sprintf("[%s size changed]: old=%d new=%d", filename, f.Size, size))
		return true
	}

	return false
}
