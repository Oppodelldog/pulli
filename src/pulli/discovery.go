package pulli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func BuildCommand(searchRoot string) {
	var whitelist []string
	var blacklist []string

	err := filepath.Walk(searchRoot, func(path string, info os.FileInfo, err error) error {
		if repoDir, ok := checkForGitRepoDir(path); ok {
			repoDir = strings.Replace(repoDir, searchRoot, "", -1)
			if shallBeIncluded(repoDir) {
				whitelist = append(whitelist, repoDir)
			} else {
				blacklist = append(blacklist, repoDir)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error discovering repositories: %v\n", err)
		return
	}

	absoluteProjectsFolder, err := filepath.Abs(searchRoot)
	if err != nil {
		fmt.Printf("error getting absolute path for '%s': %v\n", searchRoot, err)
		return
	}

	whitelistFilters := mapRepoDirsToFilterArguments(whitelist)
	blacklistFilters := mapRepoDirsToFilterArguments(blacklist)

	fmt.Printf("\n\nfrom the given information I can suggest the following commands:\n\n")
	if len(whitelistFilters) > 0 {
		fmt.Printf("%s:\n\n", FilterModeWhiteList)
		printCommand(absoluteProjectsFolder, FilterModeWhiteList, whitelistFilters)
	}
	if len(blacklistFilters) > 0 {
		fmt.Printf("%s:\n\n", FilterModeBlackList)
		printCommand(absoluteProjectsFolder, FilterModeBlackList, blacklistFilters)
	}
}

func printCommand(absoluteProjectsFolder, mode string, filters []string) {
	fmt.Printf("%s -%s=\"%s\" -%s=%s %s\n\n\n",
		CommandName,
		ArgNameDir,
		absoluteProjectsFolder,
		ArgNameFilterMode,
		mode,
		strings.Join(filters, " "),
	)
}

func shallBeIncluded(repoPath string) bool {
	fmt.Printf("Include %s [Y/n]: ", repoPath)
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Println(err)
		return false
	}

	switch char {
	case 'n':
		fallthrough
	case 'N':
		return false
	}

	return true
}

func mapRepoDirsToFilterArguments(a []string) []string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = fmt.Sprintf(`-filter="%s"`, v)
	}

	return b
}
