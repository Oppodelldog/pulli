package pulli

import (
	"os"

	"github.com/Oppodelldog/pulli/log"
)

// ValidateFlags validate program flags.
// If some flag is invalid a error message is written and false is returned
func ValidateFlags(searchRoot string, filterMode string, filters []string) bool {
	dirInfo, err := os.Stat(searchRoot)
	if err != nil {
		log.Printf("error investigating -dir '%s': %v", searchRoot, err)
		return false
	}

	if !dirInfo.IsDir() {
		log.Printf("-dir '%s': is not a directory", searchRoot)
		return false
	}

	if len(filters) > 0 && filterMode != filterModeWhiteList && filterMode != filterModeBlackList {
		log.Printf("filtermode must be either '%s' or '%s'", filterModeWhiteList, filterModeBlackList)
		return false
	}

	return true
}
