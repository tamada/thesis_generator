package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"testing"
)

func TestZip(t *testing.T) {
	goMain([]string{"ttgenerator", "../../testdata/sample.json", "-o", "../../dest.zip"})
	defer os.Remove("../../dest.zip")

	r, _ := zip.OpenReader("../../dest.zip")
	defer r.Close()
	if len(r.File) != 31 {
		t.Errorf("len(r.File) = %d, want 31", len(r.File))
	}
	if r.File[0].Name != "2023bthesis_ykino/.git/HEAD" {
		t.Errorf("r.File[0].Name = \"%s\", wants 2023bthesis_ykino/.git/HEAD", r.File[0].Name)
	}
}

func TestTarGz(t *testing.T) {
	goMain([]string{"ttgenerator", "../../testdata/sample.json", "-o", "../../dest.tar.gz"})
	defer os.Remove("../../dest.tar.gz")

	in, _ := os.Open("../../dest.tar.gz")
	defer in.Close()
	gzipIn, _ := gzip.NewReader(in)
	tarIn := tar.NewReader(gzipIn)
	count := 0
	for {
		header, err := tarIn.Next()
		if err == io.EOF {
			break
		}
		count++
		if count == 1 && header.Name != "2023bthesis_ykino/.git/HEAD" {
			t.Errorf("header.Name = \"%s\", wants 2023bthesis_ykino/.git/HEAD", header.Name)

		}
	}
	if count != 31 {
		t.Errorf("count = %d, wants 31", count)
	}
}
