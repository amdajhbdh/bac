package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Document struct {
	Pages    []Page
	Metadata Metadata
	Tables   []Table
	Images   []Image
}

type Page struct {
	Number int
	Text   string
	Images []Image
	Tables []Table
}

type Metadata struct {
	Title        string
	Author       string
	Subject      string
	Creator      string
	Producer     string
	CreationDate string
	ModDate      string
	PageCount    int
}

type Table struct {
	Page    int
	X, Y    float64
	Width   float64
	Height  float64
	Headers []string
	Rows    [][]string
}

type Image struct {
	Page   int
	X, Y   float64
	Width  float64
	Height float64
	Format string
	Data   []byte
}

type Extractor struct {
	usePDFCPU  bool
	usePoppler bool
}

func NewExtractor() *Extractor {
	return &Extractor{
		usePDFCPU:  true,
		usePoppler: true,
	}
}

func (e *Extractor) ExtractText(ctx context.Context, inputPath string) (string, error) {
	slog.Info("extracting text from PDF", "path", inputPath)

	var output string
	var err error

	if e.usePoppler {
		output, err = e.extractWithPoppler(ctx, inputPath)
		if err == nil {
			return output, nil
		}
		slog.Warn("poppler extraction failed, trying pdfcpu", "error", err)
	}

	if e.usePDFCPU {
		output, err = e.extractWithPDFCPU(ctx, inputPath)
		if err == nil {
			return output, nil
		}
		slog.Warn("pdfcpu extraction failed", "error", err)
	}

	return "", fmt.Errorf("all extractors failed: %w", err)
}

func (e *Extractor) extractWithPoppler(ctx context.Context, inputPath string) (string, error) {
	cmd := exec.CommandContext(ctx, "pdftotext", "-layout", inputPath, "-")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pdftotext: %w", err)
	}
	return stdout.String(), nil
}

func (e *Extractor) extractWithPDFCPU(ctx context.Context, inputPath string) (string, error) {
	cmd := exec.CommandContext(ctx, "pdfcpu", "text", "extract", inputPath, "-")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pdfcpu: %w", err)
	}
	return stdout.String(), nil
}

func (e *Extractor) ExtractMetadata(ctx context.Context, inputPath string) (*Metadata, error) {
	slog.Info("extracting metadata from PDF", "path", inputPath)

	cmd := exec.CommandContext(ctx, "pdfinfo", inputPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("pdfinfo: %w", err)
	}

	metadata := &Metadata{}
	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Title:") {
			metadata.Title = strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
		} else if strings.HasPrefix(line, "Author:") {
			metadata.Author = strings.TrimSpace(strings.TrimPrefix(line, "Author:"))
		} else if strings.HasPrefix(line, "Subject:") {
			metadata.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		} else if strings.HasPrefix(line, "Creator:") {
			metadata.Creator = strings.TrimSpace(strings.TrimPrefix(line, "Creator:"))
		} else if strings.HasPrefix(line, "Producer:") {
			metadata.Producer = strings.TrimSpace(strings.TrimPrefix(line, "Producer:"))
		} else if strings.HasPrefix(line, "CreationDate:") {
			metadata.CreationDate = strings.TrimSpace(strings.TrimPrefix(line, "CreationDate:"))
		} else if strings.HasPrefix(line, "ModDate:") {
			metadata.ModDate = strings.TrimSpace(strings.TrimPrefix(line, "ModDate:"))
		} else if strings.HasPrefix(line, "Pages:") {
			fmt.Sscanf(strings.TrimSpace(strings.TrimPrefix(line, "Pages:")), "%d", &metadata.PageCount)
		}
	}

	return metadata, nil
}

func (e *Extractor) ExtractImages(ctx context.Context, inputPath string, outputDir string) ([]Image, error) {
	slog.Info("extracting images from PDF", "path", inputPath)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	cmd := exec.CommandContext(ctx, "pdfimages", "-list", inputPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pdfimages list: %w", err)
	}

	cmd = exec.CommandContext(ctx, "pdfimages", "-png", inputPath, filepath.Join(outputDir, "img"))
	if err := cmd.Run(); err != nil {
		slog.Warn("pdfimages extraction failed", "error", err)
	}

	images := []Image{}
	files, _ := os.ReadDir(outputDir)
	for i, file := range files {
		if !file.IsDir() {
			data, _ := os.ReadFile(filepath.Join(outputDir, file.Name()))
			images = append(images, Image{
				Page:   i + 1,
				Format: "png",
				Data:   data,
			})
		}
	}

	return images, nil
}

func (e *Extractor) ExtractTables(ctx context.Context, inputPath string) ([]Table, error) {
	slog.Info("extracting tables from PDF", "path", inputPath)

	images, err := e.ExtractImages(ctx, inputPath, "/tmp/tables")
	if err != nil {
		return nil, err
	}

	var tables []Table
	for _, img := range images {
		tabulaData, err := e.extractWithTabula(img.Data)
		if err != nil {
			continue
		}
		tables = append(tables, Table{
			Page:    img.Page,
			Headers: tabulaData.Headers,
			Rows:    tabulaData.Rows,
		})
	}

	return tables, nil
}

type TabulaResult struct {
	Headers []string
	Rows    [][]string
}

func (e *Extractor) extractWithTabula(imageData []byte) (*TabulaResult, error) {
	tmpImg := "/tmp/tabula_input.png"
	os.WriteFile(tmpImg, imageData, 0644)

	cmd := exec.Command("python3", "-c", `
import subprocess
import sys
result = subprocess.run(['tabula-py', '-g', tmpImg], capture_output=True, text=True)
print(result.stdout)
`)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Run()

	lines := strings.Split(output.String(), "\n")
	result := &TabulaResult{
		Rows: [][]string{},
	}

	for i, line := range lines {
		if i == 0 {
			result.Headers = strings.Split(strings.TrimSpace(line), "\t")
		} else if strings.TrimSpace(line) != "" {
			result.Rows = append(result.Rows, strings.Split(strings.TrimSpace(line), "\t"))
		}
	}

	return result, nil
}

func ExtractTextFromFile(ctx context.Context, filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		extractor := NewExtractor()
		return extractor.ExtractText(ctx, filePath)
	case ".txt":
		data, err := os.ReadFile(filePath)
		return string(data), err
	case ".docx":
		return extractFromDOCX(filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

func extractFromDOCX(path string) (string, error) {
	cmd := exec.Command("pandoc", "-t", "plain", path)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pandoc: %w", err)
	}
	return stdout.String(), nil
}

type TextChunk struct {
	Text   string
	Page   int
	Source string
}

func ChunkText(text string, chunkSize int, overlap int) []TextChunk {
	var chunks []TextChunk
	words := strings.Fields(text)

	for i := 0; i < len(words); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, TextChunk{
			Text:   chunk,
			Source: "pdf",
		})
	}

	return chunks
}

func CountPages(reader io.ReaderAt, size int64) (int, error) {
	data := make([]byte, 2000)
	n, err := reader.ReadAt(data, 0)
	if err != nil && err != io.EOF {
		return 0, err
	}

	content := string(data[:n])
	idx := strings.Index(content, "/Type /Page")
	if idx == -1 {
		return 0, nil
	}

	count := 0
	for i := 0; i < len(content)-10; i++ {
		if strings.HasPrefix(content[i:], "/Type /Page") {
			count++
		}
	}

	return count, nil
}
