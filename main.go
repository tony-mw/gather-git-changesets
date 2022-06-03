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

func main() {

	cmd.Execute()

	log.Println("Starting program to check git status")

	switch os.Getenv("GIT_EVENT") {
	case "commit_main":
		ce = gitActions.CommitEvent{
			TerraformRepo: gitActions.Repo{
				Url:       os.Getenv("REPO_URL"),
				LocalPath: os.Getenv("LOCAL_REPO_PATH"),
				Branch:    os.Getenv("BRANCH"),
			},
			True: true,
		}
		tfDirs := gitActions.Do(ce)
		gitActions.FileWriter(tfDirs)
	case "pr_main":
		pre = gitActions.PREvent{
			TerraformRepo: gitActions.Repo{
				Url:       os.Getenv("REPO_URL"),
				LocalPath: os.Getenv("LOCAL_REPO_PATH"),
				Branch:    os.Getenv("BRANCH"),
			},
			BaseBranch: os.Getenv("BASE_BRANCH"),
			True:       true,
		}
		tfDirs := gitActions.Do(pre)
		gitActions.FileWriter(tfDirs)
	}
}
