package thesis_generator

import (
	"errors"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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
	// s := memory.NewStorage()
	dot, err := fs.Chroot(".git")
	if err != nil {
		return err
	}
	s := filesystem.NewStorage(dot, nil)
	r, err := git.InitWithOptions(s, fs, git.InitOptions{DefaultBranch: plumbing.Main})
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
	return errors.Join(errs...)
}
