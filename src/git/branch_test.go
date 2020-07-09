package git

import (
	"os/exec"
	"reflect"
	"testing"

	"context"
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
		{gitBranchOutput: ``,
			gitLogOutput:       `fatal: your current branch 'master' does not have any commits yet`,
			expectedBranchName: "master"},
	}

	for i, testCase := range testCases {
		t.Run(string(i), func(t *testing.T) {
			execFunc = createExecCommandMock(t, testCase.gitBranchOutput, testCase.gitLogOutput)
			path := "/tmp"
			branchName, err := GetCurrentBranchName(path)
			if err != nil {
				t.Fatalf("expected no error, but got: %v", err)
			}

			if testCase.expectedBranchName != branchName {
				t.Fatalf("expected branch: %v, but got: %v", testCase.expectedBranchName, branchName)
			}
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
		const expectedS1 = "git"
		if s1 != expectedS1 {
			t.Fatalf("expected s1 to be %v, but was: %v", expectedS1, s1)
		}

		expectedS2 := []string{"branch", "log"}
		if reflect.DeepEqual(expectedS2, s2) {
			t.Fatalf("expected s2 to be %v, but was: %v", expectedS2, s2)
		}

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
			if err == nil {
				t.Fatal("expected error, but got nil")
			}
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
	var (
		expectedType = reflect.ValueOf(execWithTimeoutFuncDef(exec.CommandContext)).Type()
		execFuncType = reflect.ValueOf(execFunc).Type()
	)

	if expectedType != execFuncType {
		t.Fatalf("expected execFunc to be opf type %T, but was %T", expectedType, execFuncType)
	}
}
