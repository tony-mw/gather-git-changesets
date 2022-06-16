package gitActions

import "github.com/go-git/go-git/v5"

type GitEvent interface {
	GatherChangeset(r *git.Repository) []string
	OpenRepo() (*git.Repository, error)
	SetWorkingDirectories(f []string) []string
}

type RepoType struct {
	Kind        string
	DirsToCheck []string
}

type Repo struct {
	Url       string
	LocalPath string
	Branch    string
	RepoType  RepoType
}

type CommitEvent struct {
	Repo Repo
	True bool
}

type PREvent struct {
	Repo       Repo
	True       bool
	BaseBranch string
}

type Loglevel struct {
	Debug string
}
