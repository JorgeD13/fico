package main

import (
	"log"
	"path/filepath"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"

	// importa tus modelos para registrarlos
	"fico/gol/internal/models"
)

func main() {
	// Ruta de salida de los tipos en el frontend (desde el módulo gol)
	outFile := filepath.FromSlash("../frontend/src/types/models.gen.ts")

	t := typescriptify.New()
	// Opciones compatibles con la versión actual
	t.CreateInterface = true // generar interfaces en lugar de clases
	t.BackupDir = ""         // no crear backups
	t.Indent = "  "          // indentación de 2 espacios

	// Registrar aquí cada struct que quieras exportar
	t.Add(models.User{})

	if err := t.ConvertToFile(outFile); err != nil {
		log.Fatal(err)
	}
	log.Printf("Tipos TypeScript generados en %s\n", outFile)
}
