package bundle

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	ARXVersion = "1.0"
)

type ARX struct {
	Metadata Metadata
	Files    []FileEntry
	TarPath  string
}

func CreateARX(outputPath string, sourceDir string, metadata Metadata) error {
	slog.Info("creating ARX archive", "output", outputPath, "source", sourceDir)

	metadata.CreatedAt = time.Now()
	metadata.ModifiedAt = time.Now()
	metadata.Format = "ARX"

	tarPath := outputPath + ".tar"

	tarFile, err := os.Create(tarPath)
	if err != nil {
		return fmt.Errorf("create tar file: %w", err)
	}
	defer tarFile.Close()

	gzWriter := gzip.NewWriter(tarFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	err = tarWriter.WriteHeader(&tar.Header{
		Name: "_metadata.json",
		Mode: 0644,
		Size: int64(len(metadataBytes)),
	})
	if err != nil {
		return fmt.Errorf("write metadata header: %w", err)
	}
	if _, err := tarWriter.Write(metadataBytes); err != nil {
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

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("create header: %w", err)
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("write header: %w", err)
		}

		if !info.IsDir() {
			fileReader, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer fileReader.Close()

			n, err := io.Copy(tarWriter, fileReader)
			if err != nil {
				return fmt.Errorf("copy file: %w", err)
			}
			totalSize += n
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walk source: %w", err)
	}

	metadata.Size = totalSize

	slog.Info("ARX archive created", "output", outputPath, "size", totalSize)
	return nil
}

func OpenARX(path string) (*ARX, error) {
	slog.Info("opening ARX archive", "path", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	arx := &ARX{
		Files:   make([]FileEntry, 0),
		TarPath: path,
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tar: %w", err)
		}

		entry := FileEntry{
			Path:        header.Name,
			Name:        filepath.Base(header.Name),
			Size:        header.Size,
			ModTime:     header.ModTime,
			IsDirectory: header.Typeflag == tar.TypeDir,
		}

		if header.Name == "_metadata.json" {
			data, err := io.ReadAll(tarReader)
			if err != nil {
				return nil, fmt.Errorf("read metadata: %w", err)
			}
			if err := json.Unmarshal(data, &arx.Metadata); err != nil {
				return nil, fmt.Errorf("parse metadata: %w", err)
			}
		}

		arx.Files = append(arx.Files, entry)
	}

	file.Close()

	slog.Info("ARX archive opened", "files", len(arx.Files))
	return arx, nil
}

func (a *ARX) Extract(destDir string) error {
	slog.Info("extracting ARX archive", "dest", destDir)

	file, err := os.Open(a.TarPath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar: %w", err)
		}

		path := filepath.Join(destDir, header.Name)

		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(path, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(path), 0755)

		outFile, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}

		_, err = io.Copy(outFile, tarReader)
		outFile.Close()
		if err != nil {
			return fmt.Errorf("copy file: %w", err)
		}
	}

	slog.Info("ARX archive extracted")
	return nil
}

func (a *ARX) ListFiles() []FileEntry {
	return a.Files
}

func (a *ARX) GetMetadata() Metadata {
	return a.Metadata
}

func ShouldUseARX(fileCount int, totalSize int64) bool {
	return fileCount > 1000 || totalSize > MaxWRAPXSize
}

func CreateBundle(outputPath string, sourceDir string, metadata Metadata) error {
	var err error

	if ShouldUseARX(0, 0) {
		err = CreateARX(outputPath, sourceDir, metadata)
	} else {
		err = CreateWRAPX(outputPath, sourceDir, metadata)
	}

	return err
}
