package gitActions

import (
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

func SetWorkingDirectories(f []string) []string {

	var WorkingDirs []string
	var dup bool = false
	for _, v := range f {
		if strings.Split(v, "/")[0] != os.Getenv("ROOT_DIR") {
			continue
		}
		for _, directory := range WorkingDirs {
			if directory == filepath.Dir(v) {
				dup = true
				break
			} else {
				dup = false
			}
		}

		if dup != true {
			WorkingDirs = append(WorkingDirs, filepath.Dir(v))
		}
	}

	if len(WorkingDirs) == 0 {
		WorkingDirs = append(WorkingDirs, "no changes")
	}

	log.Println("Working Dirs are - ", WorkingDirs)

	return WorkingDirs
}
