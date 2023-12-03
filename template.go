package thesis_generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func isIn(word string, list []string) bool {
	for _, w := range list {
		if word == w {
			return true
		}
	}
	return false
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FindTemplate(keyword string) (FS, error) {
	if strings.HasPrefix(keyword, "embed:") {
		return &embedFS{base: "templates/" + keyword[6:], fs: templates}, nil
	} else if strings.HasPrefix(keyword, "file:") {
		return openDirFS(keyword[5:])
	} else if isIn(keyword, []string{"latex", "html", "markdown", "word"}) && exist(keyword) {
		return nil, fmt.Errorf("both embed:%s and file:%s found.  please specify either keyword", keyword, keyword)
	} else if isIn(keyword, []string{"latex", "html", "markdown", "word"}) {
		return &embedFS{base: "templates/" + keyword, fs: templates}, nil
	} else if exist(keyword) {
		return openDirFS(keyword)
	}
	return nil, fmt.Errorf("template %s not found", keyword)
}

type FS interface {
	OpenDir(path string) ([]fs.DirEntry, error)
	Open(path string) (fs.File, error)
}

func openDirFS(keyword string) (FS, error) {
	if _, err := os.Stat(keyword); err != nil {
		return nil, err
	}
	return &dirFS{base: keyword}, nil
}

type dirFS struct {
	base string
}

func (dfs *dirFS) OpenDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(dfs.base, path))
}

func (dfs *dirFS) Open(path string) (fs.File, error) {
	newPath := filepath.Join(dfs.base, path)
	if _, err := os.Stat(newPath); err != nil {
		return nil, err
	}
	return os.Open(newPath)
}

//go:embed templates
var templates embed.FS

type embedFS struct {
	base string
	fs   embed.FS
}

func (efs *embedFS) OpenDir(path string) ([]fs.DirEntry, error) {
	return efs.fs.ReadDir(path)
}

func (efs *embedFS) Open(path string) (fs.File, error) {
	return efs.fs.Open(path)
}
