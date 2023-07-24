package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModeType)
	if err != nil {
		return err
	}
	defer func(fromFile *os.File) {
		err := fromFile.Close()
		if err != nil {
			fmt.Println("ошибка при закрытии файла fromPath: ", err)
		}
	}(fromFile)

	fileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	if fileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func(toFile *os.File) {
		err := toFile.Close()
		if err != nil {
			fmt.Println("ошибка при закрытии файла toPath: ", err)
		}
	}(toFile)

	if limit == 0 || limit > fileInfo.Size()-offset {
		limit = fileInfo.Size() - offset
	}

	if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if _, err := io.CopyN(toFile, fromFile, limit); err != nil {
		return err
	}

	return nil
}
