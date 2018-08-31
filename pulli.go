package pulli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Oppodelldog/pulli/git"
	"github.com/sirupsen/logrus"
)

var repoDisplayName string

// PullAll finds git repositories and pulls if the filters allow to.
func PullAll(searchRoot string, filters []string, filterMode string) {

	filter := newFilter(filters, filterMode)

	err := filepath.Walk(searchRoot, func(path string, info os.FileInfo, err error) error {

		if repoDir, ok := checkForGitRepoDir(path); ok {
			if filter.isAllowed(repoDir) {
				setCurrentRepoDisplayName(repoDir, searchRoot)

				pullRepo(repoDir)
			}
		}

		return nil
	})

	if err != nil {
		logrus.Fatalf("error while walking filesystem: %v", err)
	}
}

// ValidateFlags validate program flags.
// If some flag is invalid a error message is written and the program will exit(1)
func ValidateFlags(searchRoot string, filterMode string) {
	dirInfo, err := os.Stat(searchRoot)
	if err != nil {
		logrus.Fatalf("error investigating dir '%s': %v", searchRoot, err)
	}
	if !dirInfo.IsDir() {
		logrus.Fatalf("'dir' '%s': is not a directory", searchRoot)
	}
	if filterMode != filterModeWhiteList && filterMode != filterModeBlackList {
		logrus.Fatalf("filtermode must be either '%s' or '%s'", filterModeWhiteList, filterModeBlackList)
	}
}

func setCurrentRepoDisplayName(repoDir string, searchRoot string) {
	repoDisplayName = strings.Replace(repoDir, searchRoot, "", -1)
}

func pullRepo(repoDir string) {

	logrus.Debugf("pulling %v", repoDir)
	branchName, err := git.GetCurrentBranchName(repoDir)
	if err != nil {
		logEntry().Errorf("unable top read branch name: %v", err)
		return
	}

	result, err := git.Pull(repoDir)
	if err != nil {
		logEntry().WithField("branch", branchName).Errorf("error while pulling: %v :%v", err, truncateString(result, 50))
		logEntry().WithField("branch", branchName).Debug(result)
		return
	}

	logEntry().WithField("branch", branchName).Infof("pulled")
}

func truncateString(s string, limit int) string {
	if len(s) < limit {
		return s
	}

	return s[:limit]
}

func checkForGitRepoDir(path string) (string, bool) {
	const gitFolderName = ".git"

	if len(path) > len(gitFolderName) && path[len(path)-len(gitFolderName):] == gitFolderName {
		repoFolder, _ := filepath.Split(path)
		return repoFolder, true
	}

	return "", false
}

func logEntry() *logrus.Entry {
	return logrus.WithField("repository", repoDisplayName)
}
