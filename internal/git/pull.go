package git

import (
	"regexp"
	"strconv"
)

//regExpStats helps extracting stats from print_stat_summary_inserts_deletes.
//
// https://github.com/git/git/blob/101b3204f37606972b40fc17dec84560c22f69f6/diff.c#L2593-L2615
//
// " %d file changed"
// " %d files changed"
// ", %d insertion(+)"
// ", %d insertions(+)"
// ", %d deletion(-)"
// ", %d deletions(-)"
//
const p = `(?m) (?P<file>\d+) file(s)? changed(, (?P<ins>\d+) insertion(s)?\(\+\))?(, (?P<del>\d+) deletion(s)?\(-\))?`

var regExpStats = regexp.MustCompile(p)

type PullStats struct {
	Files      int
	Insertions int
	Deletions  int
}

// Pull executes a "git pull" in the given directory.
func Pull(dir string) (string, error) {
	return git(dir, "pull")
}

func GetStats(gitOutput string) PullStats {
	var res PullStats

	regExMatch := regExpStats.FindStringSubmatch(gitOutput)
	names := regExpStats.SubexpNames()

	for i := range regExMatch {
		if i != 0 {
			number, err := strconv.Atoi(regExMatch[i])
			if err != nil {
				continue
			}

			switch names[i] {
			case "file":
				res.Files = number
			case "ins":
				res.Insertions = number
			case "del":
				res.Deletions = number
			}
		}
	}

	return res
}
