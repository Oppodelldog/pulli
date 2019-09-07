package pulli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Oppodelldog/pulli/src/git"
	"github.com/sirupsen/logrus"
)

// PullAll finds git repositories and pulls if the filters allow to.
func PullAll(searchRoot string, filters []string, filterMode string) {

	filter := newFilter(filters, filterMode)

	err := filepath.Walk(searchRoot, func(path string, info os.FileInfo, err error) error {

		if repoDir, ok := checkForGitRepoDir(path); ok {
			if filter.isAllowed(repoDir) {
				setRepoNameForLogging(repoDir, searchRoot)
				pullRepo(repoDir)
			}
		}

		return nil
	})

	if err != nil {
		logrus.Fatalf("error while walking filesystem: %v", err)
	}
}

func pullRepo(repoDir string) {

	logrus.Debugf("pulling %v", repoDir)
	branchName, err := git.GetCurrentBranchName(repoDir)
	if err != nil {
		log.newEntry().Errorf("unable top read branch name: %v", err)
		return
	}

	result, err := git.Pull(repoDir)
	if err != nil {
		log.newEntry().WithField("branch", branchName).Errorf("error while pulling: %v :%v", err, truncateString(result, 50))
		log.newEntry().WithField("branch", branchName).Debug(result)
		return
	}

	log.newEntry().WithField("branch", branchName).Infof("pulled")
}

func checkForGitRepoDir(path string) (string, bool) {
	const gitFolderName = ".git"

	if len(path) > len(gitFolderName) && path[len(path)-len(gitFolderName):] == gitFolderName {
		repoFolder, _ := filepath.Split(path)
		return repoFolder, true
	}

	return "", false
}

func setRepoNameForLogging(repoDir string, searchRoot string) {
	log.currentRepoDisplayName = strings.Replace(repoDir, searchRoot, "", -1)
}
