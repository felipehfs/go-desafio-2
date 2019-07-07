package models

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	uuid "github.com/nu7hatch/gouuid"
)

type FileInspector struct {
	DB *sql.DB
}

func generateCheckSum(path string) chan []byte {
	ch := make(chan []byte)
	go func() {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			log.Fatal("checksumError", err)
		}
		ch <- hash.Sum(nil)
	}()

	return ch
}

func (inspector FileInspector) UpdateFile(registry FileRegistry) error {
	query := `
		UPDATE files SET filename=$2, modtime=$3,
			size=$4, path=$5, checksum=$6 WHERE id=$1
	`
	_, err := inspector.DB.Exec(query, registry.ID,
		registry.Filename, registry.ModTime, registry.Size, registry.Path, registry.CheckSum)
	return err
}

func (inspector FileInspector) InsertFile(registry FileRegistry) error {
	query := `
		INSERT INTO files (filename, modtime, size, path, Uuid, checksum)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := inspector.DB.Exec(query,
		registry.Filename, registry.ModTime,
		registry.Size, registry.Path, registry.UUID, registry.CheckSum)
	return err
}

func (inspector FileInspector) SearchFile(path string) (*FileRegistry, error) {
	query := `SELECT id, filename, modtime, size, path, Uuid FROM files WHERE path=$1`
	registry := &FileRegistry{}
	row := inspector.DB.QueryRow(query, path)
	err := row.Scan(&registry.ID, &registry.Filename,
		&registry.ModTime, &registry.Size,
		&registry.Path, &registry.UUID)

	return registry, err
}

func (inspector *FileInspector) createNewFile(path string, f os.FileInfo) {
	registry := &FileRegistry{}
	registry.Filename = f.Name()
	registry.ModTime = f.ModTime()
	registry.Size = f.Size()
	registry.Path = path
	u, err := uuid.NewV4()

	if err != nil {
		log.Fatal("uuid", err)
	}

	registry.UUID = u.String()
	checksum := generateCheckSum(path)
	registry.CheckSum = <-checksum

	err = inspector.InsertFile(*registry)

	if err != nil {
		log.Fatal("insertFile", err)
	}
}

func (inspector *FileInspector) TracingFile(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return nil
	}
	file, err := inspector.SearchFile(path)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("[saving newFile]:", path)
			inspector.createNewFile(path, f)
		} else {
			log.Fatal("searchFile", err)
		}

	}
	if file.Filename != "" && file.hasChanges(f.ModTime(), f.Name(), path, f.Size()) {
		// update values
		fmt.Println("[updating file]:", path)
		inspector.updateChanges(file, f, path)
	}
	return nil
}

func (inspector FileInspector) updateChanges(file *FileRegistry, f os.FileInfo, path string) {
	file.Filename = f.Name()
	file.ModTime = f.ModTime()
	file.Size = f.Size()
	file.Path = path
	file.CheckSum = <-generateCheckSum(file.Path)
	err := inspector.UpdateFile(*file)
	if err != nil {
		log.Fatal("updateChanges", err)
	}
}
