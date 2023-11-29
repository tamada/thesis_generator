package thesis_generator

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type RepositoryType int

const (
	GitHub RepositoryType = iota + 1
	GitLab
)

type Repository struct {
	Owner          string         `json:"owner"`
	RepositoryName string         `json:"name"`
	HostName       string         `json:"type"`
	Type           RepositoryType `json:"-"`
	Url            string         `json:"-"` //  generate from owner, repository name, and host name
}

func InitializeGit(repo *Repository, fs billy.Filesystem) error {
	s := filesystem.NewStorage(fs, nil)
	r, err := git.Init(s, fs)
	// r, err := gitInit(fs)
	if err != nil {
		return err
	}
	var errs error
	if repo.Url != "" {
		_, err := r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{repo.Url}})
		if err != nil {
			errs = errors.Join(err)
		}
	}
	err2 := doFirstCommit(r, fs)
	return errors.Join(errs, err2)
}

func gitInit(fs billy.Filesystem) (*git.Repository, error) {
	err := fs.MkdirAll(".git", 0755)
	if err != nil {
		return nil, err
	}
	configFile := fs.Join(".git", "config")
	if _, err = fs.Create(configFile); err != nil {
		return nil, err
	}

	headFile := fs.Join(".git", "HEAD")
	headIn, err := fs.Create(headFile)
	if err != nil {
		return nil, err
	}
	defer headIn.Close()
	fmt.Fprintf(headIn, "ref: refs/heads/main")

	err1 := fs.MkdirAll(fs.Join(".git", "objects"), 0755)
	err2 := fs.MkdirAll(fs.Join(".git", "refs"), 0755)
	err3 := fs.MkdirAll(fs.Join(".git", "refs", "heads"), 0755)
	err4 := fs.MkdirAll(fs.Join(".git", "refs", "tags"), 0755)
	if err := errors.Join(err1, err2, err3, err4); err != nil {
		return nil, err
	}
	// gitDir := dotgit.New(fs)
	s := filesystem.NewStorage(fs, nil)
	return git.Open(s, fs)
}

func doFirstCommit(r *git.Repository, fs billy.Filesystem) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	infos, err := fs.ReadDir(".")
	if err != nil {
		return err
	}
	var errs []error
	for _, info := range infos {
		name := info.Name()
		if isGitSystemFileName(name) {
			continue
		}
		fmt.Printf("for git add: %s\n", info.Name())
		_, err := w.Add(info.Name())
		errs = append(errs, err)
	}
	// w.Add(".")
	_, err = w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Thesis Generator",
			Email: "thesis-generator@gmail.com",
			When:  time.Now(),
		},
	})
	errs = append(errs, err)
	errs = append(errs, moveFiles(fs, r))
	return errors.Join(errs...)
}

func moveFiles(fs billy.Filesystem, r *git.Repository) error {
	err1 := fs.Remove(".git") // remove `.git` file.
	err2 := fs.MkdirAll(".git", 0755)
	err3 := fs.Rename("HEAD", fs.Join(".git", "HEAD"))
	err4 := fs.Rename("config", fs.Join(".git", "config"))
	err5 := fs.Rename("objects", fs.Join(".git", "objects"))
	err6 := fs.Rename("refs", fs.Join(".git", "refs"))
	err7 := fs.Rename("index", fs.Join(".git", "index"))
	return errors.Join(err1, err2, err3, err4, err5, err6, err7)
}

func isGitSystemFileName(name string) bool {
	return name == ".git" || name == "objects" || name == "refs" ||
		name == "index" || name == "HEAD" || name == "config"
}
