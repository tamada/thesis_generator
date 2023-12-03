package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	flag "github.com/spf13/pflag"
	tg "github.com/tamadalab/thesis_generator"
	"golang.org/x/term"
)

type options struct {
	template string
	output   string
	helpFlag bool
}

func helpMessage(name string, flags *flag.FlagSet) string {
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 80
	}
	return fmt.Sprintf(`Usage: %s [OPTIONS] [SETTING_JSON_FILE]]
OPTIONS
%s
SETTING_JSON_FILE
    JSON file that contains the setting of the thesis.`, filepath.Base(name), flags.FlagUsagesWrapped(width))
}

func buildOptions(args []string) (*flag.FlagSet, *options) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	opts := &options{}
	flags.Usage = func() { fmt.Println(helpMessage(args[0], flags)) }
	flags.StringVarP(&opts.template, "template", "t", "latex", "specify the template file. available: latex, markdown, html, and word. default: latex")
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this help message")
	flags.StringVarP(&opts.output, "output", "o", "dest", "specify the destination directory or archive file (zip or tar.gz). default: dest")
	return flags, opts
}

func performImpl(thesis *tg.Thesis, opts *options, fs billy.Filesystem) error {
	if err := tg.Generate(thesis, opts.template, fs); err != nil {
		return err
	}

	return tg.InitializeGit(thesis.Repository, fs)
}

func writeResult(thesis *tg.Thesis, opts *options, fs billy.Filesystem, at tg.ArchiveType) error {
	if at != tg.Dir {
		writer, err := os.Create(opts.output)
		if err != nil {
			return err
		}
		defer writer.Close()
		if err := tg.Archive(thesis.Id, fs, at, writer); err != nil {
			return err
		}
	}
	return nil
}

func perform(opts *options, settingJson string) int {
	thesis, err := loadJson(settingJson)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	fs, at := getOutputFs(opts.output)
	if err := performImpl(thesis, opts, fs); err != nil {
		fmt.Println(err.Error())
		return 2
	}
	if err := writeResult(thesis, opts, fs, at); err != nil {
		fmt.Println(err.Error())
		return 2
	}
	return 0
}

func getOutputFs(output string) (billy.Filesystem, tg.ArchiveType) {
	var t tg.ArchiveType = tg.Dir
	if strings.HasSuffix(output, ".zip") {
		return memfs.New(), tg.Zip
	} else if strings.HasSuffix(output, ".tar.gz") {
		return memfs.New(), tg.TarGz
	}
	os.Mkdir(output, 0755)
	return osfs.New(output), t
}

func loadJson(settingJson string) (*tg.Thesis, error) {
	if settingJson == "" {
		return loadJsonImpl(os.Stdin)
	}
	file, err := os.Open(settingJson)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return loadJsonImpl(file)
}

func appendData(thesis *tg.Thesis) *tg.Thesis {
	if thesis.Author.StudentId == "" {
		thesis.Author.StudentId = strings.Split(thesis.Author.Email, "@")[0][1:]
	}
	thesis.Id = fmt.Sprintf("%4d%cthesis_%s", thesis.Year, thesis.Degree[0], thesis.Repository.Owner)
	thesis.Repository.RepositoryName = thesis.Id
	thesis.Repository.Url = fmt.Sprintf("https://%s.com/%s/%s", thesis.Repository.HostName, thesis.Repository.Owner, thesis.Repository.RepositoryName)
	return thesis
}

func validateString(target, message string) error {
	if target == "" {
		return fmt.Errorf("%s is empty", message)
	}
	return nil
}

func validateEmail(email string) error {
	terms := strings.Split(email, "@")
	if len(terms) != 2 || len(terms[0]) == 0 || len(terms[1]) == 0 {
		return fmt.Errorf("%s: invalid email", email)
	}
	return nil
}

func validate(thesis *tg.Thesis) (*tg.Thesis, error) {
	var errs []error
	if thesis.Degree != "bachelor" && thesis.Degree != "master" && thesis.Degree != "doctoral" {
		errs = append(errs, fmt.Errorf("%s: unknown degree", thesis.Degree))
	}
	errs = append(errs, validateString(thesis.Title, "title"))

	errs = append(errs, validateString(thesis.Author.Name, "author name"))
	errs = append(errs, validateString(thesis.Author.Name, "author email"))
	errs = append(errs, validateEmail(thesis.Author.Email))
	if thesis.Year == 0 {
		now := time.Now()
		year := now.Year()
		if now.Month() < 4 {
			year = year - 1
		}
		thesis.Year = year
	}
	return appendData(thesis), errors.Join(errs...)
}

func loadJsonImpl(r io.Reader) (*tg.Thesis, error) {
	thesis := &tg.Thesis{}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, thesis); err != nil {
		return nil, err
	}
	return validate(thesis)
}

func goMain(args []string) int {
	flags, opts := buildOptions(args)
	if err := flags.Parse(args); err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if opts.helpFlag {
		fmt.Println(helpMessage(args[0], flags))
		return 0
	}
	newArgs := flags.Args()
	if len(newArgs) < 2 {
		fmt.Println("no setting json file specified")
		fmt.Println(helpMessage(args[0], flags))
		return 1
	}
	return perform(opts, newArgs[1])
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
