package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"imagecopier/internal/copier"
	"imagecopier/internal/scanner"
	"imagecopier/internal/utils"
)

var workers = flag.Int("w", 8, "Number of concurrent workers")
var noScreenshots = flag.Bool("n", false, "Skip screenshot images")

func main() {
	flag.Parse()

	if *noScreenshots {
		utils.NoScreenshots = true
	}

	args := flag.Args()
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-w N] [-n] <destino> <modo> <carpeta1> [carpeta2] ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Modo: copy | skip\n")
		os.Exit(1)
	}

	destDir := args[0]
	mode := args[1]
	sourceDirs := args[2:]

	if mode != "copy" && mode != "skip" {
		fmt.Fprintf(os.Stderr, "Error: modo debe ser 'copy' o 'skip'\n")
		os.Exit(1)
	}

	for _, dir := range sourceDirs {
		if !utils.DirExists(dir) {
			fmt.Fprintf(os.Stderr, "Error: la carpeta '%s' no existe\n", dir)
			os.Exit(1)
		}
	}

	if err := utils.CreateDirIfNotExists(destDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: no se pudo crear la carpeta destino: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("escaneando carpetas...\n")
	start := time.Now()

	files, err := scanner.New().Scan(sourceDirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al escanear: %v\n", err)
		os.Exit(1)
	}

	totalFiles := len(files)
	elapsed := time.Since(start)
	rate := 0.0
	if elapsed.Seconds() > 0 {
		rate = float64(totalFiles) / elapsed.Seconds()
	}

	fmt.Printf("Encontrados %d archivos de imagen en %.2f segundos (%.1f files/s)\n\n", totalFiles, elapsed.Seconds(), rate)

	if totalFiles == 0 {
		fmt.Println("No se encontraron imagenes para copiar.")
		os.Exit(0)
	}

	fmt.Printf("Copiando archivos con %d workers...\n", *workers)

	c := copier.New(destDir, mode, *workers)
	c.SetTotal(totalFiles)

	copyStart := time.Now()
	copied, err := c.CopySimple(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al copiar: %v\n", err)
		os.Exit(1)
	}

	copyElapsed := time.Since(copyStart)
	copyRate := 0.0
	if copyElapsed.Seconds() > 0 {
		copyRate = float64(copied) / copyElapsed.Seconds()
	}

	totalElapsed := time.Since(start)
	fmt.Printf("\n=== Resumen ===\n")
	fmt.Printf("Archivos encontrados: %d\n", totalFiles)
	fmt.Printf("Archivos copiados: %d\n", copied)
	fmt.Printf("Tiempo total: %.2f segundos\n", totalElapsed.Seconds())
	if copyRate > 0 {
		fmt.Printf("Velocidad: %.1f files/s\n", copyRate)
	}
}