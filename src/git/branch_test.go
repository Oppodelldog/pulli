package git

import (
	"os/exec"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
)

var originals = struct {
	execFunc execWithTimeoutFuncDef
}{
	execFunc: execFunc,
}

func restoreOriginals() {
	execFunc = originals.execFunc
}

func TestGetCurrentBranchName_FindsBranchInCommandOutput(t *testing.T) {
	defer restoreOriginals()

	testCases := []struct {
		gitBranchOutput    string
		gitLogOutput       string
		expectedBranchName string
	}{
		{gitBranchOutput: `develop`, expectedBranchName: ""},
		{gitBranchOutput: `* develop`, expectedBranchName: "develop"},
		{gitBranchOutput: "* develop\nmaster", expectedBranchName: "develop"},
		{gitBranchOutput: "release/v1.0.1\n* develop\nmaster", expectedBranchName: "develop"},
		{gitBranchOutput: "feature/xyz\ndevelop\n* release/v1.0.1\nmaster", expectedBranchName: "release/v1.0.1"},
		{gitBranchOutput: ``, gitLogOutput: `fatal: your current branch 'master' does not have any commits yet`, expectedBranchName: "master"},
	}

	for i, testCase := range testCases {
		t.Run(string(i), func(t *testing.T) {
			execFunc = createExecCommandMock(t, testCase.gitBranchOutput, testCase.gitLogOutput)
			path := "/tmp"
			branchName, err := GetCurrentBranchName(path)
			assert.NoError(t, err)
			assert.Exactly(t, testCase.expectedBranchName, branchName)
		})
	}
}

func createExecCommandMock(t *testing.T, gitBranchOutput, gitLogOutput string) execWithTimeoutFuncDef {
	var callCounter int
	return func(c context.Context, s1 string, s2 ...string) *exec.Cmd {
		callCounter++
		if s1 != "git" {
			t.Fatalf("expected 'git' to be executed, but got %v", s1)
		}
		switch callCounter {
		case 1:
			if s2[0] != "branch" {
				t.Fatalf("expected 'branch' to be parameter of git command, but got %v", s2[0])
			}
			return exec.Command("echo", gitBranchOutput)
		case 2:
			if s2[0] != "log" {
				t.Fatalf("expected 'branch' to be parameter of git command, but got %v", s2[0])
			}
			return exec.Command("echo", gitLogOutput)
		}
		t.Fatalf("not expected to be called that often! (%v times)", callCounter)
		return nil
	}
}

func TestGetCurrentBranchName_ExpectGitCommandIsCalledProperly(t *testing.T) {
	t.SkipNow()
	defer restoreOriginals()

	execFunc = func(c context.Context, s1 string, s2 ...string) *exec.Cmd {
		assert.Exactly(t, "git", s1)
		assert.Equal(t, []string{"branch", "log"}, s2)

		return exec.Command("", "")
	}

	GetCurrentBranchName("test-folder")
}

func TestGetCurrentBranchName_ReturnsErrorFromOneOfBothPossibleCommandExecutions(t *testing.T) {
	defer restoreOriginals()

	for i := 1; i <= 2; i++ {
		t.Run(string(i), func(t *testing.T) {
			execFunc = createExecFuncErrorStub(i)

			_, err := GetCurrentBranchName("test-folder")
			assert.Error(t, err)
		})
	}
}

func createExecFuncErrorStub(errorCommandAtCall int) execWithTimeoutFuncDef {
	var counter int
	return func(c context.Context, s1 string, s2 ...string) *exec.Cmd {
		counter++
		if counter == errorCommandAtCall {
			return exec.Command("thiscommandwillnotbefound", "")
		}
		//noinspection SpellCheckingInspection
		return exec.Command("echo", "works")
	}
}

func TestDefaultExecFuncIsExecCommand(t *testing.T) {
	assert.IsType(t, execWithTimeoutFuncDef(exec.CommandContext), execFunc)
}