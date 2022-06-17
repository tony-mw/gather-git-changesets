package gitActions

import (
	"fmt"
	"log"
	"os"
)

func FileWriter(dirs []string) {
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
