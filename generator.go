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

func copyFile(fs billy.Filesystem, source string) error {
	in, err := templates.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := fs.Create(filepath.Base(source))
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
	err1 := templating(fs, "templates/latex/thesis.tex", thesis.Id+".tex", thesis)
	err2 := templating(fs, "templates/latex/thesis.bib", thesis.Id+".bib", thesis)
	err3 := templating(fs, "templates/latex/llmk.toml", "llmk.toml", thesis)
	err4 := templating(fs, "templates/latex/README.md", "README.md", thesis)
	err5 := templating(fs, "templates/latex/gitignore", ".gitignore", thesis)
	err6 := copyFile(fs, "templates/latex/csg-thesis.bst")
	err7 := copyFile(fs, "templates/latex/csg-thesis.sty")
	return errors.Join(err1, err2, err3, err4, err5, err6, err7)
}
