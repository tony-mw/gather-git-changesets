package gitActions

import "strings"

func (r Loglevel) Check() bool {

	positiveKeywords := []string{"on", "true", "yes"}

	for _, v := range positiveKeywords {
		if strings.ToLower(r.Debug) == v {
			return true
		}
	}
	return false
}
