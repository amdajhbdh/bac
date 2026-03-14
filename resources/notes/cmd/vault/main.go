package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: vault <command>")
		fmt.Println("Commands: sync, status, ocr")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "sync":
		syncVault()
	case "status":
		vaultStatus()
	case "ocr":
		runOCR()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func syncVault() {
	fmt.Println("Syncing vault...")

	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = "/home/med/Documents/bac/resources/notes"
	cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "vault: auto-sync")
	cmd.Dir = "/home/med/Documents/bac/resources/notes"
	cmd.Run()

	cmd = exec.Command("git", "push", "origin", "main")
	cmd.Dir = "/home/med/Documents/bac/resources/notes"
	cmd.Run()

	fmt.Println("✓ Vault synced")
}

func vaultStatus() {
	cmd := exec.Command("git", "status", "--short")
	cmd.Dir = "/home/med/Documents/bac/resources/notes"
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func runOCR() {
	vaultDir := "/home/med/Documents/bac/resources/notes"
	outputDir := filepath.Join(vaultDir, "03-Resources", "OCR")
	os.MkdirAll(outputDir, 0755)

	// Find PDFs
	var wg sync.WaitGroup
	pdfChan := make(chan string, 100)

	// Start workers
	workers := 8
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pdf := range pdfChan {
				ocrPDF(pdf, outputDir)
			}
		}()
	}

	// Walk directories
	dirs := []string{
		filepath.Join(vaultDir, "03-Resources"),
		filepath.Join(vaultDir, "07-Assets", "PDFs"),
		filepath.Join(vaultDir, "mauritanian-bac"),
	}

	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".pdf" {
				pdfChan <- path
			}
			return nil
		})
	}
	close(pdfChan)
	wg.Wait()

	fmt.Println("✓ OCR complete")
}

func ocrPDF(pdfPath, outputDir string) {
	baseName := filepath.Base(pdfPath[:len(pdfPath)-4])
	outputPath := filepath.Join(outputDir, baseName+".md")

	// Check if already done
	if _, err := os.Stat(outputPath); err == nil {
		return
	}

	// Convert to images
	imgPath := filepath.Join(outputDir, baseName)
	cmd := exec.Command("pdftoppm", "-r", "200", "-png", pdfPath, imgPath)
	cmd.Run()

	// OCR each image
	outFile, _ := os.Create(outputPath)
	defer outFile.Close()

	files, _ := filepath.Glob(imgPath + "-*.png")
	var wg sync.Mutex

	for _, img := range files {
		cmd := exec.Command("tesseract", "-l", "ara+fra+eng", img, "stdout")
		out, _ := cmd.Output()

		wg.Lock()
		outFile.Write(out)
		outFile.Write([]byte("\n"))
		wg.Unlock()

		os.Remove(img)
	}

	fmt.Printf("✓ OCR: %s\n", baseName)
}
