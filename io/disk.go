package io

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mdw-tools/hugoinho/contracts"
)

type Disk struct{}

func (Disk) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
func (Disk) WriteFile(path string, content []byte, perm os.FileMode) error {
	return os.WriteFile(path, content, perm)
}
func (Disk) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
func (Disk) Walk(root string) (result chan contracts.FileSystemEntry) {
	result = make(chan contracts.FileSystemEntry)
	go func() {
		defer close(result)
		_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			result <- contracts.FileSystemEntry{
				Root:     root,
				Path:     path,
				DirEntry: entry,
				Error:    err,
			}
			return err
		})
	}()
	return result
}
