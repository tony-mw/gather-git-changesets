package gitActions

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
	"strings"
	"time"
)

func OpenRepoCommon(RepoPath string) (*git.Repository, error) {

	repo, err := git.PlainOpen(RepoPath)

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (c CommitEvent) OpenRepo() (*git.Repository, error) {

	log.Println("Opening repo for Commit Event")
	repo, err := OpenRepoCommon(c.Repo.LocalPath)
	if err != nil {
		log.Fatal(err)
	}

	return repo, nil
}

func (p PREvent) OpenRepo() (*git.Repository, error) {

	log.Println("Opening Repo for PR Event")

	repo, err := OpenRepoCommon(p.Repo.LocalPath)
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
			log.Println("No error: ", from.Path(), to.Path())
		} else {
			log.Println("The file was renamed from ", from.Path(), "to ", to.Path())
			filesChanged = append(filesChanged, to.Path())
		}
	}

	return filesChanged
}

func (c CommitEvent) GatherChangeset(r *git.Repository) []string {

	var currentCommitHash *object.Commit
	var previousCommitHash *object.Commit
	var myBranchRef plumbing.Hash

	refs, _ := r.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if strings.Contains(string(ref.Name()), fmt.Sprintf("refs/heads/%s", os.Getenv("BRANCH"))) {
				fmt.Println("REF IS: ", ref.Name())
				if logger.Check() {
					log.Println(ref.Name(), ref.Hash())
				}
				myBranchRef = ref.Hash()
			}
		}

		return nil
	})
	if myBranchRef.String() == "0000000000000000000000000000000000000000" {
		log.Fatal("Ref doesn't exist")
	}
	ref, _ := r.Log(&git.LogOptions{
		From: myBranchRef,
	})
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
				if logger.Check() {
					log.Println(ref.Name(), ref.Hash())
				}
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

func (p PREvent) GatherChangeset(repo *git.Repository) []string {

	var filesChanged []string
	var currentCommit *object.Commit
	var myBranchRef plumbing.Hash
	var commits []*object.Commit
	var appendOn bool = true

	currentCommit, err := GetBaseCommit(repo, p.BaseBranch)
	log.Printf("Base Branch Ref Hash is %s\n", currentCommit.Hash)
	if err != nil {
		log.Fatal(err)
	}

	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		//if ref.Type() == plumbing.HashReference {
		//	if strings.Contains(string(ref.Name()),"refs/heads") {
		//		if strings.Contains(string(ref.Name()), r.Repo.Branch) {
		//			log.Printf("Current Branch Ref name is: %s\n Current Branch Ref Ref Hash is %s\n", ref.Name(), ref.Hash())
		//			myBranchRef = ref.Hash()
		//		}
		//	}
		//}
		if ref.Type() == plumbing.HashReference {
			if strings.Contains(string(ref.Name()), p.Repo.Branch) {
				log.Printf("Current Branch Ref name is: %s\n Current Branch Ref Ref Hash is %s\n", ref.Name(), ref.Hash())
				myBranchRef = ref.Hash()
			}
		}

		return nil
	})

	//LogFilter := fmt.Sprintf("%s/*", os.Getenv("ROOT_DIR"))
	ref, err := repo.Log(&git.LogOptions{
		From: myBranchRef,
		//FileName: &LogFilter,
	})

	if err != nil {
		log.Println("Error with log : ", err)
	}
	//object.NewCommitAllIter()
	start := time.Now()
	ref.ForEach(func(c *object.Commit) error {
		if appendOn == false {
			return nil
		}

		ancestor, err := c.IsAncestor(currentCommit)
		if err != nil {
			log.Println("History is not traversable: ", err)
		}

		if logger.Check() {
			log.Printf("Is ancestor: %t", ancestor)
		}
		if ancestor {
			//Check if this was a merge back into main
			if logger.Check() {
				log.Println(c.Message)
			}
			appendOn = false
			return nil
		}
		commits = append(commits, c)
		//if appendOn {
		//	if logger.Check() {
		//		log.Println(c.Message)
		//	}
		//	commits = append(commits, c)
		//}

		return nil
	})
	finish := time.Now()
	tDiff := finish.Sub(start)
	fmt.Println("Took: ", tDiff)
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
	if logger.Check() {
		log.Println("Files changed are: ", filesChanged)
	}
	return filesChanged
}
