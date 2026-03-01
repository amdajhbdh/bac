package bundle

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	WRAPXVersion = "1.0"
	MaxWRAPXSize = 100 * 1024 * 1024 // 100MB
)

type Metadata struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Subject     string                 `json:"subject"`
	Chapter     string                 `json:"chapter"`
	Language    string                 `json:"language"`
	Keywords    []string               `json:"keywords"`
	Authors     []string               `json:"authors"`
	CreatedAt   time.Time              `json:"created_at"`
	ModifiedAt  time.Time              `json:"modified_at"`
	License     string                 `json:"license"`
	Format      string                 `json:"format"`
	Size        int64                  `json:"size"`
	Checksum    string                 `json:"checksum"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
}

type WRAPX struct {
	Version  string
	Metadata Metadata
	Files    []FileEntry
	Reader   *zip.ReadCloser
}

type FileEntry struct {
	Path           string    `json:"path"`
	Name           string    `json:"name"`
	Size           int64     `json:"size"`
	CompressedSize int64     `json:"compressed_size"`
	ModTime        time.Time `json:"mod_time"`
	ContentType    string    `json:"content_type"`
	IsDirectory    bool      `json:"is_directory"`
}

func CreateWRAPX(outputPath string, sourceDir string, metadata Metadata) error {
	slog.Info("creating WRAPX bundle", "output", outputPath, "source", sourceDir)

	metadata.CreatedAt = time.Now()
	metadata.ModifiedAt = time.Now()
	metadata.Format = "WRAPX"

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	metaWriter, err := zipWriter.Create("_metadata.json")
	if err != nil {
		return fmt.Errorf("create metadata file: %w", err)
	}
	if _, err := metaWriter.Write(metadataBytes); err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	var totalSize int64
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		fileReader, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open file %s: %w", path, err)
		}
		defer fileReader.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("create header: %w", err)
		}
		header.Name = relPath

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("create zip header: %w", err)
		}

		n, err := io.Copy(writer, fileReader)
		if err != nil {
			return fmt.Errorf("copy file: %w", err)
		}

		totalSize += n
		slog.Debug("added file to bundle", "path", relPath, "size", n)

		return nil
	})

	if err != nil {
		return fmt.Errorf("walk source dir: %w", err)
	}

	metadata.Size = totalSize

	slog.Info("WRAPX bundle created", "output", outputPath, "size", totalSize)
	return nil
}

func OpenWRAPX(path string) (*WRAPX, error) {
	slog.Info("opening WRAPX bundle", "path", path)

	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	wrapx := &WRAPX{
		Reader:  reader,
		Files:   make([]FileEntry, 0),
		Version: WRAPXVersion,
	}

	for _, file := range reader.File {
		entry := FileEntry{
			Path:           file.Name,
			Name:           filepath.Base(file.Name),
			Size:           int64(file.UncompressedSize64),
			CompressedSize: int64(file.CompressedSize64),
			ModTime:        file.Modified,
			ContentType:    "application/octet-stream",
			IsDirectory:    file.FileInfo().IsDir(),
		}

		if file.Name == "_metadata.json" {
			rc, err := file.Open()
			if err != nil {
				reader.Close()
				return nil, fmt.Errorf("open metadata: %w", err)
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				reader.Close()
				return nil, fmt.Errorf("read metadata: %w", err)
			}
			if err := json.Unmarshal(data, &wrapx.Metadata); err != nil {
				reader.Close()
				return nil, fmt.Errorf("parse metadata: %w", err)
			}
		}

		wrapx.Files = append(wrapx.Files, entry)
	}

	slog.Info("WRAPX bundle opened", "files", len(wrapx.Files))
	return wrapx, nil
}

func (w *WRAPX) Extract(destDir string) error {
	slog.Info("extracting WRAPX bundle", "dest", destDir)

	for _, file := range w.Reader.File {
		path := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(path), 0755)

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("open file %s: %w", file.Name, err)
		}
		defer rc.Close()

		outFile, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create file %s: %w", file.Name, err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, rc); err != nil {
			return fmt.Errorf("write file %s: %w", file.Name, err)
		}
	}

	slog.Info("WRAPX bundle extracted")
	return nil
}

func (w *WRAPX) ReadFile(name string) ([]byte, error) {
	for _, file := range w.Reader.File {
		if file.Name == name {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("open file: %w", err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("file not found: %s", name)
}

func (w *WRAPX) ListFiles() []FileEntry {
	return w.Files
}

func (w *WRAPX) GetMetadata() Metadata {
	return w.Metadata
}

func (w *WRAPX) Close() error {
	return w.Reader.Close()
}

func ValidateWRAPX(path string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer reader.Close()

	hasMetadata := false
	for _, file := range reader.File {
		if file.Name == "_metadata.json" {
			hasMetadata = true
			break
		}
	}

	if !hasMetadata {
		return fmt.Errorf("missing _metadata.json")
	}

	return nil
}

func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func BundleFiles(files []string, output string, metadata Metadata) error {
	tempDir, err := os.MkdirTemp("", "wrapx-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read file %s: %w", file, err)
		}
		err = os.WriteFile(filepath.Join(tempDir, filepath.Base(file)), data, 0644)
		if err != nil {
			return fmt.Errorf("write file: %w", err)
		}
	}

	return CreateWRAPX(output, tempDir, metadata)
}

func Unbundle(bundlePath string, destDir string) error {
	wrapx, err := OpenWRAPX(bundlePath)
	if err != nil {
		return fmt.Errorf("open bundle: %w", err)
	}
	defer wrapx.Close()

	return wrapx.Extract(destDir)
}

func ListContents(bundlePath string) error {
	wrapx, err := OpenWRAPX(bundlePath)
	if err != nil {
		return fmt.Errorf("open bundle: %w", err)
	}
	defer wrapx.Close()

	fmt.Printf("WRAPX Bundle: %s\n", bundlePath)
	fmt.Printf("Version: %s\n", wrapx.Version)
	fmt.Printf("Title: %s\n", wrapx.Metadata.Title)
	fmt.Printf("Subject: %s\n", wrapx.Metadata.Subject)
	fmt.Printf("Size: %d bytes\n", wrapx.Metadata.Size)
	fmt.Printf("\nFiles:\n")

	for _, file := range wrapx.Files {
		fmt.Printf("  %s (%d bytes)\n", file.Path, file.Size)
	}

	return nil
}

var stdout io.Writer = os.Stdout

func SetOutput(w io.Writer) {
	stdout = w
}

func Printf(format string, args ...interface{}) {
	fmt.Fprintf(stdout, format, args...)
}

func Example() {
	metadata := Metadata{
		Title:       "Bac Math 2024",
		Description: "Mathematics exam questions 2024",
		Subject:     "Mathematics",
		Chapter:     "Algebra",
		Language:    "fr",
		Keywords:    []string{"bac", "math", "exam"},
		Authors:     []string{"Ministry of Education"},
		License:     "CC-BY-SA",
	}

	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(metadata)
	_ = buf
}
