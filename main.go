package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/felipehfs/godesafio2/models"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "filestask"
)

func main() {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	searchDir := flag.String("search", ".", "Search directory")
	flag.Parse()

	conn, err := sql.Open("postgres", sqlInfo)
	defer conn.Close()

	if err != nil {
		fmt.Println("A conexão com o banco de dados falhou.")
		log.Fatal(err)
	}

	inspector := &models.FileInspector{DB: conn}
	err = filepath.Walk(*searchDir, inspector.TracingFile)

	if err != nil {
		log.Fatal(err)
	}
}
