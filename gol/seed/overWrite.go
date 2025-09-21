//go:build ignore

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"fico/gol/db"
	seeder "fico/gol/seed"
)

// runOverWrite applies AutoMigrate and SeedDemo over all .db files in gol/data/
func runOverWrite() error {
	dataDir := filepath.FromSlash("data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}
	// Ensure demo.db exists
	demoPath := filepath.Join(dataDir, "demo.db")
	touchedDemo := false
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".db") {
			continue
		}
		abs := filepath.Join(dataDir, name)
		if err := processDB(abs); err != nil {
			return err
		}
		if name == "demo.db" {
			touchedDemo = true
		}
	}
	if !touchedDemo {
		if err := processDB(demoPath); err != nil {
			return err
		}
	}
	return nil
}

func processDB(path string) error {
	gdb, err := db.OpenSQLiteAt(path)
	if err != nil {
		return err
	}
	if err := seeder.AutoMigrate(gdb); err != nil {
		return err
	}
	if err := seeder.SeedDemo(gdb); err != nil {
		return err
	}
	log.Printf("OK: %s\n", path)
	return nil
}

func main() {
	if err := runOverWrite(); err != nil {
		log.Fatal(err)
	}
}
