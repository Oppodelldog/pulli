package git

import (
	"regexp"
)

// GetCurrentBranchName executes 'git branch' to get the current branch in the given directory.
// In a fresh branch it parses as a fallback fatal error message containing the branch name.
func GetCurrentBranchName(directory string) (string, error) {
	output, _ := git(directory, "branch")

	branchName := getBranchNameFromGitOutput(output)

	if branchName == "" {
		output, err := git(directory, "log")
		if err != nil {
			return "", err
		}

		branchName = getBranchNameFromGitLogOutput(output)
	}

	return branchName, nil
}

func getBranchNameFromGitOutput(gitOutput string) string {
	return extractFromString(gitOutput, `(?m)^\* (.*)$`)
}

func getBranchNameFromGitLogOutput(gitOutput string) string {
	// https://github.com/git/git/blob/ed843436dd4924c10669820cc73daf50f0b4dabd/revision.c#L2303
	pattern := `(?m)^fatal: your current branch '(.*)' does not have any commits yet$`

	return extractFromString(gitOutput, pattern)
}

func extractFromString(s, regexPattern string) string {
	var re = regexp.MustCompile(regexPattern)

	matches := re.FindAllStringSubmatch(s, 1)

	if len(matches) > 0 {
		if len(matches[0]) > 1 {
			return matches[0][1]
		}
	}

	return ""
}
