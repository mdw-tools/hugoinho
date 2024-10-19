package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mdw-tools/hugoinho/contracts"
)

type InMemoryFileSystem struct {
	ModTime      time.Time
	lock         *sync.RWMutex
	Files        map[string]*MemoryFile
	ErrReadFile  map[string]error
	ErrWriteFile map[string]error
	ErrMkdirAll  map[string]error
	ErrWalkFunc  map[string]error
}

func NewInMemoryFileSystem() *InMemoryFileSystem {
	return &InMemoryFileSystem{
		ModTime:      time.Now(),
		lock:         new(sync.RWMutex),
		Files:        make(map[string]*MemoryFile),
		ErrReadFile:  make(map[string]error),
		ErrWriteFile: make(map[string]error),
		ErrMkdirAll:  make(map[string]error),
		ErrWalkFunc:  make(map[string]error),
	}
}

func (this *InMemoryFileSystem) ReadFile(path string) ([]byte, error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	err := this.ErrReadFile[path]
	if err != nil {
		return nil, err
	}
	return this.Files[path].content, nil
}

func (this *InMemoryFileSystem) WriteFile(path string, content []byte, perm os.FileMode) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	err := this.ErrWriteFile[path]
	if err != nil {
		return err
	}
	this.Files[path] = &MemoryFile{
		content: content,
		name:    filepath.Base(path),
		mode:    perm,
		modTime: this.ModTime,
	}
	return nil
}

func (this *InMemoryFileSystem) MkdirAll(path string, perm os.FileMode) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	err := this.ErrMkdirAll[path]
	if err != nil {
		return err
	}
	components := strings.Split(path, string(os.PathSeparator))
	for x := 1; x < len(components)+1; x++ {
		p := strings.Join(components[:x], string(os.PathSeparator))
		this.Files[p] = &MemoryFile{
			name:    filepath.Base(p),
			mode:    perm,
			modTime: this.ModTime,
			isDir:   true,
		}
	}
	return nil
}

func (this *InMemoryFileSystem) Walk(root string) (result chan contracts.FileSystemEntry) {
	result = make(chan contracts.FileSystemEntry)
	go func() {
		this.lock.RLock()
		defer this.lock.RUnlock()
		defer close(result)
		root = filepath.Clean(root)
		var paths []string
		for path := range this.Files {
			if path == root || strings.HasPrefix(path, root+"/") {
				paths = append(paths, path)
			}
		}
		sort.Strings(paths)
		for _, path := range paths {
			result <- contracts.FileSystemEntry{
				Root:     root,
				Path:     path,
				DirEntry: this.Files[path],
				Error:    this.ErrWalkFunc[path],
			}
		}
	}()
	return result
}

///////////////////////////////////////////////////////////////////

type MemoryFile struct {
	content []byte
	name    string
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (this *MemoryFile) Content() string            { return string(this.content) }
func (this *MemoryFile) Name() string               { return this.name }
func (this *MemoryFile) Size() int64                { return int64(len(this.content)) }
func (this *MemoryFile) Mode() os.FileMode          { return this.mode }
func (this *MemoryFile) ModTime() time.Time         { return this.modTime }
func (this *MemoryFile) IsDir() bool                { return this.isDir }
func (this *MemoryFile) Sys() any                   { panic("NOT NEEDED") }
func (this *MemoryFile) Type() fs.FileMode          { return this.mode }
func (this *MemoryFile) Info() (fs.FileInfo, error) { return this, nil }
