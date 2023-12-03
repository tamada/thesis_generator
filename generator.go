package thesis_generator

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

//go:embed templates
var templates embed.FS

func Generate(thesis *Thesis, format Format, fs billy.Filesystem) error {
	switch format {
	case LaTeX:
		return generateLaTeX(thesis, fs)
	case HTML, Markdown, MicrosoftWord:
		return errors.New("not implemented yet")
	default:
		return fmt.Errorf("%d: unknown format", format)
	}
}

func templating(fs billy.Filesystem, source, dest string, thesis *Thesis) error {
	out, err := fs.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	templ, err := template.ParseFS(templates, source)
	if err != nil {
		return err
	}
	return templ.Execute(out, thesis)
}

func copyFile(fs billy.Filesystem, source string, toFile string) error {
	in, err := templates.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := fs.Create(toFile)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

func generateLaTeX(thesis *Thesis, fs billy.Filesystem) error {
	return generateImpl("templates/latex", "", thesis, fs)
}

func generateImpl(source, prefix string, thesis *Thesis, fs billy.Filesystem) error {
	entries, err := templates.ReadDir(source)
	if err != nil {
		return err
	}
	var errs []error
	for _, entry := range entries {
		if entry.IsDir() {
			if err := generateImpl(filepath.Join(source, entry.Name()), filepath.Join(prefix, entry.Name()), thesis, fs); err != nil {
				errs = append(errs, err)
			}
		} else {
			toFile := filepath.Join(prefix, mapFileName(entry.Name(), thesis))
			fromFile := filepath.Join(source, entry.Name())
			if isTemplateTarget(entry.Name()) {
				if err := templating(fs, fromFile, toFile, thesis); err != nil {
					errs = append(errs, err)
				}
			} else {
				if err := copyFile(fs, fromFile, toFile); err != nil {
					errs = append(errs, err)
				}
			}

		}
	}
	return errors.Join(errs...)
}

func isTemplateTarget(name string) bool {
	switch filepath.Ext(name) {
	case ".png", ".jpg", ".jpeg", ".pdf", ".gif", ".svg", ".sty", ".bst":
		return false
	default:
		return true
	}
}

func mapFileName(name string, thesis *Thesis) string {
	dir := filepath.Dir(name)
	switch filepath.Base(name) {
	case "thesis.tex":
		return filepath.Join(dir, thesis.Id+".tex")
	case "thesis.bib":
		return filepath.Join(thesis.Id + ".bib")
	default:
		return name
	}
}
