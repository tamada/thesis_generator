package thesis_generator

import (
	"errors"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
)

func Generate(thesis *Thesis, templatePath string, fs billy.Filesystem) error {
	sourceFS, err := FindTemplate(templatePath)
	if err != nil {
		return err
	}
	return generateImpl(sourceFS, ".", thesis, fs)
}

func templating(sourceFS FS, destFS billy.Filesystem, source, dest string, thesis *Thesis) error {
	out, err := destFS.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	templ, err := template.New(source).Funcs(funcs()).ParseFS(sourceFS, source)
	if err != nil {
		return err
	}
	return templ.Execute(out, thesis)
}

func funcs() template.FuncMap {
	return map[string]any{
		"split": strings.Split,
		"trim":  strings.TrimSpace,
	}
}

func copyFile(sourceFS FS, destFS billy.Filesystem, source, toFile string) error {
	in, err := sourceFS.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := destFS.Create(toFile)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

func generateImpl(sourceFS FS, prefix string, thesis *Thesis, destFS billy.Filesystem) error {
	entries, err := sourceFS.OpenDir(prefix)
	if err != nil {
		return err
	}
	var errs []error
	for _, entry := range entries {
		if entry.IsDir() {
			if err := generateImpl(sourceFS, filepath.Join(prefix, entry.Name()), thesis, destFS); err != nil {
				errs = append(errs, err)
			}
		} else {
			toFile := filepath.Join(prefix, mapFileName(entry.Name(), thesis))
			fromFile := filepath.Join(prefix, entry.Name())
			if isTemplateTarget(entry.Name()) {
				templating(sourceFS, destFS, fromFile, toFile, thesis)
				errs = append(errs, err)
			} else {
				copyFile(sourceFS, destFS, fromFile, toFile)
				errs = append(errs, err)
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
		return filepath.Join(dir, thesis.Id+".bib")
	case "gitignore":
		return filepath.Join(dir, ".gitignore")
	default:
		return name
	}
}
