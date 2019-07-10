package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/felipehfs/godesafio2/models"
	_ "github.com/lib/pq"
)

func main() {

	host := os.Getenv("HOST_PG")
	port := 5432
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := "filestask"

	sqlInfo := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, password, dbname)

	searchDir := flag.String("search", ".", "Search directory")
	flag.Parse()

	conn, err := sql.Open("postgres", sqlInfo)

	if err != nil {
		fmt.Println("A conexão com o banco de dados falhou.")
		log.Fatal(err)

	}
	defer conn.Close()

	inspector := &models.FileInspector{DB: conn}
	err = filepath.Walk(*searchDir, inspector.TracingFile)

	if err != nil {
		log.Fatal(err)
	}
}
