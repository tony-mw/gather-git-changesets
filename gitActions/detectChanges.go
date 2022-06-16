package gitActions

import (
	"fmt"
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
	Dirs := g.SetWorkingDirectories(filesChanged)

	return Dirs
}

func FormatDir(d string) string {
	fmt.Println("Formatting: ", d)
	var finalDir string
	if logger.Check() {
		log.Println("Formatting")
	}
	tld :=  strings.Split(filepath.Dir(d), "/")[0]
	serviceName := strings.Split(filepath.Dir(d), "/")[1]
	if logger.Check() {
		log.Println("Name is: ", serviceName)
	}
	finalDir = strings.Join([]string{fmt.Sprintf("%s", tld), serviceName}, "/")
	if logger.Check() {
		log.Println("Final dir is: ", finalDir)
	}
	return finalDir
}

func (p PREvent) SetWorkingDirectories(f []string) []string {
	return SetWorkingDirectoriesCommon(f, p.Repo.RepoType)
}

func (c CommitEvent) SetWorkingDirectories(f []string) []string {
	return SetWorkingDirectoriesCommon(f, c.Repo.RepoType)
}

func contains(s string, sl []string) bool {
	for _, v := range sl {
		if s == v {
			return true
		}
	}
	return false
}

func SetWorkingDirectoriesCommon(f []string, p RepoType) []string {
	fmt.Println(f)
	var WorkingDirs []string
	var dup bool = false
	var format bool
	if p.Kind == "app" {
		fmt.Println("app")
		format = true
	} else if p.Kind == "terraform" {
		format = false
	}

	for _, v := range f {
		var finalDir string
		if contains(strings.Split(v, "/")[0], p.DirsToCheck) == false {
			fmt.Println("skipping")
			continue
		}
		for _, directory := range WorkingDirs {
			if format {
				finalDir = FormatDir(v)
			} else {
				finalDir = filepath.Dir(v)
			}
			if directory == finalDir {
				dup = true
				if logger.Check() {
					log.Println("There was a dup ", directory, finalDir)
				}
				break
			} else {
				dup = false
			}
		}

		if dup != true {
			if len(finalDir) == 0 {
				fmt.Println("first")
				if format {
					finalDir = FormatDir(v)
				} else {
					finalDir = filepath.Dir(v)
				}
			}
			WorkingDirs = append(WorkingDirs, finalDir)
		}
	}

	if len(WorkingDirs) == 0 {
		WorkingDirs = append(WorkingDirs, "no changes")
	}

	log.Println("Working Dirs are - ", WorkingDirs)

	return WorkingDirs
}
