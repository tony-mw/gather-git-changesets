/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"gitActions/cmd"
	"gitActions/gitActions"
	"log"
	"os"
)

// Commit Main Example
// export GIT_EVENT=commit_main
// export REPO_URL=https://bitbucket.dentsplysirona.com/scm/atopoc/cirrus-poc-terraform.git
// export LOCAL_REPO_PATH="/Users/tonyprestifilippo/git/dentsply/cirru/cirrus-poc-terraform"
// export BRANCH=main
// export ROOT_DIR=terraform

// PR Main Example
// export GIT_EVENT=pr_main
// export REPO_URL=https://bitbucket.dentsplysirona.com/scm/atopoc/cirrus-poc-gitops.git
// export LOCAL_REPO_PATH="/Users/tonyprestifilippo/git/dentsply/cirru/cirrus-poc-gitops"
// export BRANCH=wip
// export BASE_BRANCH=main

var ce gitActions.CommitEvent
var pre gitActions.PREvent

var terraformRepoDirectories = []string{"terraform"}
var applicationRepoDirectories = []string{"services", "pkg", "k8s"}

func checkRepoType() gitActions.RepoType {
	switch os.Getenv("REPO_TYPE") {
	case "terraform":
		return gitActions.RepoType{
			Kind:        "terraform",
			DirsToCheck: terraformRepoDirectories,
		}
	case "app":
		return gitActions.RepoType{
			Kind:        "app",
			DirsToCheck: applicationRepoDirectories,
		}
	}
	return gitActions.RepoType{}
}

func main() {

	cmd.Execute()

	log.Println("Starting program to check git status")

	repoType := checkRepoType()

	switch os.Getenv("GIT_EVENT") {
	case "commit_main":
		ce = gitActions.CommitEvent{
			Repo: gitActions.Repo{
				Url:       os.Getenv("REPO_URL"),
				LocalPath: os.Getenv("LOCAL_REPO_PATH"),
				Branch:    os.Getenv("BRANCH"),
				RepoType:  repoType,
			},
			True: true,
		}
		tfDirs := gitActions.Do(ce)
		gitActions.FileWriter(tfDirs)
	case "pr_main":
		pre = gitActions.PREvent{
			Repo: gitActions.Repo{
				Url:       os.Getenv("REPO_URL"),
				LocalPath: os.Getenv("LOCAL_REPO_PATH"),
				Branch:    os.Getenv("BRANCH"),
				RepoType:  repoType,
			},
			BaseBranch: os.Getenv("BASE_BRANCH"),
			True:       true,
		}
		Dirs := gitActions.Do(pre)
		gitActions.FileWriter(Dirs)
	}
}
