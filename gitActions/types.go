package gitActions

import "github.com/go-git/go-git/v5"

type GitEvent interface {
	GatherChangeset(r *git.Repository) []string
	OpenRepo() (*git.Repository, error)
}

type Repo struct {
	Url       string
	LocalPath string
	Branch    string
}

type CommitEvent struct {
	TerraformRepo Repo
	True          bool
}

type PREvent struct {
	TerraformRepo Repo
	True          bool
	BaseBranch    string
}

type Loglevel struct {
	Debug string
}
