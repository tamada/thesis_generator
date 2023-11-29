package thesis_generator

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

type ArchiveType int

const (
	Zip ArchiveType = iota + 1
	TarGz
	Dir
)

func Archive(thesisId string, fs billy.Filesystem, archiveType ArchiveType, dest io.Writer) error {
	switch archiveType {
	case Zip:
		return zipArchive(fs, dest, thesisId)
	case TarGz:
		return tarGzArchive(fs, dest, thesisId)
	}
	return fmt.Errorf("%d: unknown archive type", archiveType)
}

func tarGzArchive(fs billy.Filesystem, dest io.Writer, thesisId string) error {
	gzipWriter := gzip.NewWriter(dest)
	tarWriter := tar.NewWriter(gzipWriter)
	var errs error
	infos, err := fs.ReadDir(".")
	if err != nil {
		return err
	}
	for _, info := range infos {
		err := writeToTarEntry(fs, info, tarWriter, thesisId)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	err1 := tarWriter.Close()
	err2 := gzipWriter.Close()
	return errors.Join(errs, err1, err2)
}

func writeToTarEntry(fs billy.Filesystem, info os.FileInfo, tarWriter *tar.Writer, thesisId string) error {
	header := &tar.Header{
		Name: fs.Join(thesisId, info.Name()),
		Size: info.Size(),
		Mode: 0600,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	in, err := fs.Open(info.Name())
	if err != nil {
		return err
	}
	defer in.Close()
	if _, err := io.Copy(tarWriter, in); err != nil {
		return err
	}
	return nil
}

func zipArchive(fs billy.Filesystem, dest io.Writer, thesisId string) error {
	zipWriter := zip.NewWriter(dest)
	defer zipWriter.Close()
	return zipArchiveImpl(fs, zipWriter, ".", thesisId)
}

func zipArchiveImpl(fs billy.Filesystem, zipWriter *zip.Writer, fromDir, toDir string) error {
	// fmt.Printf("zipArchiveImpl: %s -> %s\n", fromDir, toDir)
	infos, err := fs.ReadDir(fromDir)
	if err != nil {
		return err
	}
	var errs error
	for _, info := range infos {
		fmt.Printf("zipArchiveImpl:  %s/%s\n", fromDir, info.Name())
		fromFile := fs.Join(fromDir, info.Name())
		if info.IsDir() {
			err := zipArchiveImpl(fs, zipWriter, fromFile, fs.Join(toDir, info.Name()))
			errs = errors.Join(errs, err)
			continue
		}
		err := writeToZipEntry(fs, zipWriter, fromFile, toDir)
		errs = errors.Join(errs, err)
	}
	return errs
}

func writeToZipEntry(fs billy.Filesystem, zipWriter *zip.Writer, path string, toDir string) error {
	fmt.Printf("writeToZipEntry: %s\n", path)
	out, err := zipWriter.Create(fs.Join(toDir, filepath.Base(path)))
	if err != nil {
		return errors.Join(err, fmt.Errorf("%s: create error", fs.Join(toDir, filepath.Base(path))))
	}
	in, err := fs.Open(path)
	if err != nil {
		return errors.Join(err, fmt.Errorf("%s: open error", path))
	}
	defer in.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
