package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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

func (inspector FileInspector) UpdateFile(registry FileRegistry) (err error) {

	body, err := json.Marshal(registry.Items)

	if err != nil {
		return
	}

	query := fmt.Sprintf(`UPDATE files SET items = '%s'::jsonb WHERE id=$1`, string(body))

	stmt, err := inspector.DB.Prepare(query)

	if err != nil {
		return
	}
	_, err = stmt.Exec(query, registry.ID)

	return
}

func (inspector FileInspector) InsertFile(registry FileRegistry) (err error) {

	body, err := json.Marshal(registry.Items)

	if err != nil {
		return
	}

	query := fmt.Sprintf(`INSERT INTO files (id,updated_at,items) VALUES ($1,$2,'%s'::jsonb)`, string(body))

	smtp, err := inspector.DB.Prepare(query)

	if err != nil {
		return
	}

	id, _ := uuid.NewV4()

	var date = time.Now().Format("2006-01-02T15:04:05Z07:00")

	_, err = smtp.Exec(id.String(), date)

	return
}

func (inspector FileInspector) SearchFile(path string) (registry *FileRegistry, err error) {

	query := "SELECT id,items FROM files WHERE items->>'path'=$1"

	registry = &FileRegistry{}

	row := inspector.DB.QueryRow(query, path)

	var items string

	err = row.Scan(&registry.ID, &items)

	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(items), &registry.Items)

	return
}

func (inspector *FileInspector) createNewFile(path string, f os.FileInfo) (err error) {

	registry := &FileRegistry{}

	registry.Items = AddItems(f, path)

	err = inspector.InsertFile(*registry)

	return
}

func (inspector *FileInspector) TracingFile(path string, f os.FileInfo, err error) error {

	if f.IsDir() {
		return nil
	}

	file, err := inspector.SearchFile(path)

	if err != nil {

		if err == sql.ErrNoRows {
			fmt.Println("[saving newFile]:", path)

			err = inspector.createNewFile(path, f)
			if err != nil {
				log.Fatal("createNewFile", err)
			}

		} else {
			log.Fatal("searchFile", err.Error())
		}

	}
	if file.Items.Filename != "" && file.hasChanges(f.ModTime(), f.Name(), path, f.Size()) {
		// update values
		fmt.Println("[updating file]:", path)
		inspector.updateChanges(file, f, path)
	}
	return err
}

func (inspector FileInspector) updateChanges(file *FileRegistry, f os.FileInfo, path string) {

	err := inspector.UpdateFile(*file)
	if err != nil {
		log.Fatal("updateChanges", err)
	}
}

func AddItems(f os.FileInfo, path string) (items Items) {

	checksum := generateCheckSum(path)

	items.CheckSum = <-checksum
	items.Filename = f.Name()
	items.ModTime = f.ModTime()
	items.Size = f.Size()
	items.Path = path
	return

}
