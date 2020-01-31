package pulli

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/Oppodelldog/pulli/log"
)

func TestValidateFlags(t *testing.T) {
	const testDir = "/tmp/pulli/flagvalidation"
	testFile := path.Join(testDir, "testFile")
	prepareTestDir(t, testDir, testFile)
	defer cleanupTestDir(t, testDir)

	testCases := map[string]struct {
		searchRoot       string
		filters          []string
		filterMode       string
		expectedResult   bool
		expectedLogInput loggerInput
	}{
		"searchRoot does not exist": {
			searchRoot:       "/DOES-NOT_EXIST",
			expectedLogInput: loggerInput{format: "error investigating -dir '%s': %v", v: []interface{}{"/DOES-NOT_EXIST", &os.PathError{Op: "stat", Path: "/DOES-NOT_EXIST", Err: errors.New("no such file or directory")}}},
			expectedResult:   false,
		},
		"searchRoot is no directory": {
			searchRoot:       testFile,
			expectedLogInput: loggerInput{format: "-dir '%s': is not a directory"},
			expectedResult:   false,
		},
		"searchRoot is valid, filters and filemode is empty, ok": {
			searchRoot:     testDir,
			expectedResult: true,
		},
		"searchRoot is valid, with filter set, empty filemode is not allowed": {
			searchRoot:       testDir,
			filters:          []string{"some-filter"},
			filterMode:       "",
			expectedLogInput: loggerInput{format: "filtermode must be either '%s' or '%s'", v: []interface{}{"/tmp/pulli/flagvalidation/testFile"}},
			expectedResult:   false,
		},
		"searchRoot is valid and filemode (whitelist) is allowed": {
			searchRoot:     testDir,
			filters:        []string{"some-filter"},
			filterMode:     filterModeWhiteList,
			expectedResult: true,
		},
		"searchRoot is valid and filemode (blacklist) is allowed": {
			searchRoot:     testDir,
			filters:        []string{"some-filter"},
			filterMode:     filterModeBlackList,
			expectedResult: true,
		},
	}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			logMock := new(loggerMock)
			log.SetPrintf(logMock.Printf)
			log.SetFatalf(logMock.Fatalf)

			result := ValidateFlags(testData.searchRoot, testData.filterMode, testData.filters)

			if testData.expectedResult != result {
				t.Fatalf("result should have been %v, but was: %v", testData.expectedResult, result)
			}

			latestLogInput := logMock.GetLatestPrintf()
			if testData.expectedLogInput.format != latestLogInput.format {
				t.Fatalf("latest logged input format should have been %v, but was: %v",
					testData.expectedLogInput.format,
					latestLogInput.format,
				)
			}
			for key, inputValue := range testData.expectedLogInput.v {
				switch expected := inputValue.(type) {
				case *os.PathError:
					loggedPathError := latestLogInput.v[key].(*os.PathError)
					if expected.Path != loggedPathError.Path {
						t.Fatalf("PathError.Path should have been %v, but was: %v", expected.Path, loggedPathError.Path)
					}
					if expected.Op != loggedPathError.Op {
						t.Fatalf("PathError.Op should have been %v, but was: %v", expected.Op, loggedPathError.Op)
					}
					if expected.Err.Error() != loggedPathError.Err.Error() {
						t.Fatalf("PathError.Err.Error() should have been %v, but was: %v", expected.Err.Error(), loggedPathError.Err.Error())
					}
				}
			}
		})
	}
}

type loggerInput struct {
	format string
	v      []interface{}
}

type loggerMock struct {
	printfInputs []loggerInput
	fatalfInputs []loggerInput
}

func (l *loggerMock) Printf(format string, v ...interface{}) {
	l.printfInputs = append(l.printfInputs, loggerInput{format: format, v: v})
}

func (l *loggerMock) Fatalf(format string, v ...interface{}) {
	l.fatalfInputs = append(l.fatalfInputs, loggerInput{format: format, v: v})
}

func (l *loggerMock) GetLatestPrintf() loggerInput {
	if len(l.printfInputs) == 0 {
		return loggerInput{}
	}

	return l.printfInputs[len(l.printfInputs)-1]
}

func cleanupTestDir(t *testing.T, testDir string) {
	err := os.RemoveAll(testDir)
	if err != nil {
		t.Fatal(err)
	}
}

func prepareTestDir(t *testing.T, testDir string, testFile string) {
	err := os.MkdirAll(testDir, 0777)
	if err != nil {
		t.Fatalf("Error preparing test dir: %v", err)
	}

	_, err = os.Create(testFile)
	if err != nil {
		t.Fatalf("Error preparing testfile: %v", err)
	}
}
