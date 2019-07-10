package models

import (
	"fmt"
	"log"
	"time"
)

type FileRegistry struct {
	ID    string
	Items Items
}

type Items struct {
	Filename string    `json:"filename"`
	Path     string    `json:"path"`
	ModTime  time.Time `json:"time"`
	CheckSum []byte    `json:"check"`
	Size     int64     `json:"size"`
}

func (f *FileRegistry) hasChanges(modTime time.Time, filename string, path string, size int64) bool {

	if f.Items.Filename != filename {
		log.Println(fmt.Sprintf("[%s filename changed]: old=%s new=%s", filename, f.Items.Filename, filename))
		return true
	}
	if f.Items.Path != path {
		log.Println(fmt.Sprintf("[%s path changed]: old=%s new=%s", filename, f.Items.Path, path))
		return true
	}

	if f.Items.Size != size {
		log.Println(fmt.Sprintf("[%s size changed]: old=%d new=%d", filename, f.Items.Size, size))
		return true
	}

	return false
}
