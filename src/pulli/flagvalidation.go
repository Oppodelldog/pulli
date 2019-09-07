package pulli

import (
	"github.com/sirupsen/logrus"
	"os"
)

// ValidateFlags validate program flags.
// If some flag is invalid a error message is written and false is returned
func ValidateFlags(searchRoot string, filterMode string, filters []string) bool {
	dirInfo, err := os.Stat(searchRoot)
	if err != nil {
		logrus.Errorf("error investigating -dir '%s': %v", searchRoot, err)
		return false
	}

	if !dirInfo.IsDir() {
		logrus.Errorf("-dir '%s': is not a directory", searchRoot)
		return false
	}

	if len(filters) > 0 && filterMode != filterModeWhiteList && filterMode != filterModeBlackList {
		logrus.Errorf("filtermode must be either '%s' or '%s'", filterModeWhiteList, filterModeBlackList)
		return false
	}

	return true
}
