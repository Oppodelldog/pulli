package git

// Pull executes a "git pull" in the given directory.
func Pull(dir string) (string, error) {
	return git(dir, "pull")
}
