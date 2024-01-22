package copy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func ScanFileForEmail(filePath, email string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// ตรวจสอบว่าเนื้อหาของไฟล์มีอีเมลหรือไม่
	if strings.Contains(string(content), email) {
		fmt.Println("File contains email:", filePath)
	}

	return nil
}
func ZipFiles(folderPath, zipFileName string) error {
	// List all files in the folder
	dirEntries, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	// Create a zip file
	zipFileName = filepath.Join(folderPath, zipFileName)
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Iterate through files and add them to the zip archive
	for _, dirEntry := range dirEntries {
		filePath := filepath.Join(folderPath, dirEntry.Name())
		if !dirEntry.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			// Get the file information
			fileInfo, err := file.Stat()
			if err != nil {
				return err
			}

			// Create a new file header
			fileHeader, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				return err
			}

			// Set the file header name to the original file name
			fileHeader.Name = dirEntry.Name()

			// Create a new zip file entry
			writer, err := zipWriter.CreateHeader(fileHeader)
			if err != nil {
				return err
			}

			// Copy the file content to the zip file entry
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
