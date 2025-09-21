package api

import (
	"log"
	"os"
	"path/filepath"

	"fico/gol/db"
)

func Start() {
	_ = os.MkdirAll(filepath.FromSlash("data"), 0755)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/demo.db"
	}
	gdb, err := db.OpenSQLiteAt(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	SetUpServer(gdb)
}
