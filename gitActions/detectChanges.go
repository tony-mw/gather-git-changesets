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
	Dirs := SetWorkingDirectories(filesChanged)

	return Dirs
}

func FormatDir(d string) string {
	var finalDir string
	if logger.Check() {
		log.Println("Root dir is services")
	}
	serviceName := strings.Split(filepath.Dir(d), "/")[1]
	if logger.Check() {
		log.Println("Service name is: ", serviceName)
	}
	finalDir = strings.Join([]string{fmt.Sprintf("%s", os.Getenv("ROOT_DIR")), serviceName}, "/")
	if logger.Check() {
		log.Println("Final dir is: ", finalDir)
	}
	return finalDir
}

func SetWorkingDirectories(f []string) []string {

	var WorkingDirs []string
	var dup bool = false
	for _, v := range f {
		var finalDir string

		if strings.Split(v, "/")[0] != os.Getenv("ROOT_DIR") {
			continue
		}
		for _, directory := range WorkingDirs {
			if os.Getenv("ROOT_DIR") == "services" {
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
				if os.Getenv("ROOT_DIR") == "services" {
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
