package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveMultipartFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func CleanupFiles(paths []string) {
	for _, p := range paths {
		_ = os.Remove(p)
	}
}
