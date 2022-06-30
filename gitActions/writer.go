package gitActions

import (
	"fmt"
	"log"
	"os"
)

func FileWriter(dirs []string) {
	err := os.Remove("changeSet")
	if err != nil {
		log.Println("File doesn't exist - creating and appending dirs")
	}
	for _, d := range dirs {

		f, err := os.OpenFile("changeSet", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}

		defer f.Close()

		_, err = f.WriteString(fmt.Sprintf("%s\n", d))
		if err != nil {
			log.Println(err)
		}
	}
}
