package gitActions

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var logger = Loglevel{
	Debug: os.Getenv("DEBUG"),
}

func Do(g GitEvent) []string {

	repoObject, _ := g.OpenRepo()
	filesChanged := g.GatherChangeset(repoObject)
	tfDirs := SetTerraformWorkingDirectories(filesChanged)

	return tfDirs
}

func SetTerraformWorkingDirectories(f []string) []string {

	var tfWorkingDirs []string
	var dup bool = false

	for _, v := range f {

		if strings.Split(v, "/")[0] != "terraform" {
			continue
		}

		for _, directory := range tfWorkingDirs {
			if directory == filepath.Dir(v) {
				dup = true
				break
			} else {
				dup = false
			}
		}

		if dup != true {
			tfWorkingDirs = append(tfWorkingDirs, filepath.Dir(v))
		}
	}

	if len(tfWorkingDirs) == 0 {
		tfWorkingDirs = append(tfWorkingDirs, "no changes")
	}

	log.Println("Working Dirs are - ", tfWorkingDirs)

	return tfWorkingDirs
}

func OpenRepoCommon(RepoPath string) (*git.Repository, error) {

	repo, err := git.PlainOpen(RepoPath)

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r CommitEvent) OpenRepo() (*git.Repository, error) {

	log.Println("Opening repo for Commit Event")

	repo, err := OpenRepoCommon(r.TerraformRepo.LocalPath)
	if err != nil {
		log.Fatal(err)
	}

	return repo, nil
}

func (r PREvent) OpenRepo() (*git.Repository, error) {

	log.Println("Opening Repo for PR Event")

	repo, err := OpenRepoCommon(r.TerraformRepo.LocalPath)
	if err != nil {
		log.Fatal(err)
	}

	return repo, nil
}

func RetrieveFiles(filesPatched []diff.FilePatch) []string {

	var filesChanged []string

	for _, fileList := range filesPatched {
		from, to := fileList.Files()

		if fileList.IsBinary() {
			log.Println("Ignoring binary files")
			continue
		}

		if to == nil {
			log.Println("A new file was created.")
			filesChanged = append(filesChanged, from.Path())
			continue
		}

		if from == nil {
			log.Println("A file was deleted.")
			filesChanged = append(filesChanged, to.Path())
			continue
		}

		if from.Path() == to.Path() {
			filesChanged = append(filesChanged, from.Path())
		} else {
			log.Println("err")
		}
	}

	return filesChanged
}

func (t CommitEvent) GatherChangeset(r *git.Repository) []string {

	var currentCommitHash *object.Commit
	var previousCommitHash *object.Commit

	ref, _ := r.Log(&git.LogOptions{})
	counter := 0

	ref.ForEach(func(c *object.Commit) error {
		if counter == 0 {
			currentCommitHash = c
		} else if counter == 1 {
			previousCommitHash = c
		}
		counter += 1

		return nil
	})

	if logger.Check() {
		log.Println("Current Commit Hash is: ", currentCommitHash.Hash, "Previous Commit Hash is: ", previousCommitHash.Hash)
	}
	diff, err := currentCommitHash.Patch(previousCommitHash)
	if err != nil {
		log.Fatal(err)
	}

	filesPatched := diff.FilePatches()
	filesChanged := RetrieveFiles(filesPatched)

	return filesChanged
}

func GetBaseCommit(repo *git.Repository, baseBranch string) (*object.Commit, error) {

	var currentCommitHash *object.Commit
	var myBranchRef plumbing.Hash

	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if strings.Contains(string(ref.Name()), fmt.Sprintf("refs/heads/%s", baseBranch)) {
				log.Println(ref.Name(), ref.Hash())
				myBranchRef = ref.Hash()
			}
		}

		return nil
	})
	w, _ := repo.Worktree()

	err := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", baseBranch)),
		Force:  true,
	})
	if err != nil {
		log.Println("Error checking out branch")
		log.Println(err)
	}

	ref, _ := repo.Log(&git.LogOptions{
		From: myBranchRef,
	})
	counter := 0

	ref.ForEach(func(c *object.Commit) error {
		if counter == 0 {
			log.Println("From Main\n", c)
			currentCommitHash = c
		}
		counter += 1
		return nil
	})
	return currentCommitHash, nil
}

func (r PREvent) GatherChangeset(repo *git.Repository) []string {

	var filesChanged []string
	var currentCommit *object.Commit
	var myBranchRef plumbing.Hash
	var commits []*object.Commit
	var appendOn bool = true

	currentCommit, err := GetBaseCommit(repo, r.BaseBranch)
	log.Printf("Base Branch Ref Hash is %s\n", currentCommit.Hash)
	if err != nil {
		log.Fatal(err)
	}

	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if strings.Contains(string(ref.Name()), r.TerraformRepo.Branch) {
				log.Printf("Current Branch Ref name is: %s\n Current Branch Ref Ref Hash is %s\n", ref.Name(), ref.Hash())
				myBranchRef = ref.Hash()
			}
		}

		return nil
	})

	ref, _ := repo.Log(&git.LogOptions{
		From: myBranchRef,
	})

	ref.ForEach(func(c *object.Commit) error {

		ancestor, _ := c.IsAncestor(currentCommit)

		if logger.Check() {
			log.Printf("Is ancestor: %t", ancestor)
		}
		if ancestor {
			appendOn = false
			return nil
		}
		if appendOn {
			commits = append(commits, c)
		}

		return nil
	})

	if logger.Check() {
		log.Println(commits)
	}

	for _, v := range commits {

		d, err := v.Patch(currentCommit)
		if err != nil {
			log.Fatal(err)
		}

		filesPatched := d.FilePatches()

		f := RetrieveFiles(filesPatched)

		for _, v := range f {
			filesChanged = append(filesChanged, v)
		}
	}

	return filesChanged
}
