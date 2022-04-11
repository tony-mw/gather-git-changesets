package gitActions

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"path/filepath"
	"strings"
)

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

func (t CommitEvent) GatherChangeset(r *git.Repository) []string {
	var filesChanged []string
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

	fmt.Println(currentCommitHash.Hash, previousCommitHash.Hash)
	diff, err := currentCommitHash.Patch(previousCommitHash)
	if err != nil {
		log.Fatal(err)
	}
	filesPatched := diff.FilePatches()
	for _, fileList := range filesPatched {
		from, to := fileList.Files()
		if fileList.IsBinary() {
			fmt.Println("Ignoring binary files")
			continue
		}
		if from.Path() == to.Path() {
			filesChanged = append(filesChanged, from.Path())
		} else {
			fmt.Println("err")
		}
	}

	return filesChanged
}

func GetBaseCommit(repo *git.Repository, baseBranch string) (*object.Commit, error) {
	var currentCommitHash *object.Commit
	var myBranchRef plumbing.Hash

	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if strings.Contains(string(ref.Name()), baseBranch) {
				fmt.Println(ref.Name(), ref.Hash())
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
		fmt.Println("Error checking out branch")
		fmt.Println(err)
	}

	ref, _ := repo.Log(&git.LogOptions{
		From:  myBranchRef,
		Order: git.LogOrderCommitterTime,
	})
	counter := 0

	ref.ForEach(func(c *object.Commit) error {
		if counter == 0 {
			fmt.Println("From Main\n", c)
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
	//var InitialCommitHash *object.Commit
	var myBranchRef plumbing.Hash

	currentCommit, err := GetBaseCommit(repo, r.BaseBranch)
	if err != nil {
		log.Fatal(err)
	}
	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if strings.Contains(string(ref.Name()), r.TerraformRepo.Branch) {
				fmt.Println(ref.Name(), ref.Hash())
				myBranchRef = ref.Hash()
			}
		}

		return nil
	})

	ref, _ := repo.Log(&git.LogOptions{
		From:  myBranchRef,
		Order: git.LogOrderDefault,
		All:   false,
	})
	var commits []*object.Commit
	var appendOn bool = true
	ref.ForEach(func(c *object.Commit) error {

		if c.Hash == currentCommit.Hash {
			appendOn = false
			return nil
		}
		if appendOn {
			commits = append(commits, c)
		}
		return nil
	})
	fmt.Println(commits)

	for _, v := range commits {
		diff, err := v.Patch(currentCommit)
		if err != nil {
			log.Fatal(err)
		}
		filesPatched := diff.FilePatches()
		for _, fileList := range filesPatched {
			if fileList.IsBinary() {
				fmt.Println("Ignoring binary files")
				continue
			}
			from, to := fileList.Files()
			if from.Path() == to.Path() {
				filesChanged = append(filesChanged, from.Path())
			} else {
				fmt.Println("err")
			}
		}
	}

	return filesChanged
}
