package pulli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Oppodelldog/pulli/internal/git"
)

const CommandName = "pulli"
const ArgNameDir = "dir"
const ArgNameFilter = "filter"
const ArgNameFilterMode = "filtermode"
const SubCommandNameBuildCommand = "build-command"

// PullAll finds git repositories and pulls if the filters allow to.
func PullAll(searchRoot string, filters []string, filterMode string) {
	filter := newFilter(filters, filterMode)

	err := filepath.Walk(searchRoot, func(path string, info os.FileInfo, err error) error {
		if repoDir, ok := checkForGitRepoDir(path); ok {
			if filter.isAllowed(repoDir) {
				pullRepo(repoDir)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("error while walking filesystem: %v", err)
	}
}

func pullRepo(repoDir string) {
	log.Printf("pulling %v", repoDir)
	branchName, err := git.GetCurrentBranchName(repoDir)

	if err != nil {
		log.Printf("unable top read branch name: %v", err)
		return
	}

	result, err := git.Pull(repoDir)
	if err != nil {
		log.Printf("%s: error while pulling: %v :%v", branchName, err, truncateString(result, 50))
		log.Printf("%s: %v", branchName, result)

		return
	}

	stats := git.GetStats(result)
	statsString := fmt.Sprintf(" f: %v +%v -%v", stats.Files, stats.Insertions, stats.Deletions)

	log.Printf("%s: pulled (%s)", branchName, statsString)
}

func checkForGitRepoDir(path string) (string, bool) {
	const gitFolderName = ".git"

	if len(path) > len(gitFolderName) && path[len(path)-len(gitFolderName):] == gitFolderName {
		repoFolder, _ := filepath.Split(path)
		return repoFolder, true
	}

	return "", false
}
