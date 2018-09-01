package pulli

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
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
		expectedLogEntry *logrus.Entry
	}{
		"searchRoot does not exist": {
			searchRoot: "/DOES-NOT_EXIST",
			expectedLogEntry: &logrus.Entry{
				Message: "error investigating -dir '/DOES-NOT_EXIST': stat /DOES-NOT_EXIST: no such file or directory",
				Level:   logrus.ErrorLevel,
			},
			expectedResult: false,
		},
		"searchRoot is no directory": {
			searchRoot: testFile,
			expectedLogEntry: &logrus.Entry{
				Message: "-dir '/tmp/pulli/flagvalidation/testFile': is not a directory",
				Level:   logrus.ErrorLevel,
			},
			expectedResult: false,
		},
		"searchRoot is valid, filters and filemode is empty, ok": {
			searchRoot:     testDir,
			expectedResult: true,
		},
		"searchRoot is valid, with filter set, empty filemode is not allowed": {
			searchRoot: testDir,
			filters:    []string{"some-filter"},
			filterMode: "",
			expectedLogEntry: &logrus.Entry{
				Message: "filtermode must be either 'whitelist' or 'blacklist'",
				Level:   logrus.ErrorLevel,
			},
			expectedResult: false,
		},
		"searchRoot is valid and filemode (whitelist) is allowed": {
			searchRoot:       testDir,
			filters:          []string{"some-filter"},
			filterMode:       filterModeWhiteList,
			expectedLogEntry: nil,
			expectedResult:   true,
		},
		"searchRoot is valid and filemode (blacklist) is allowed": {
			searchRoot:       testDir,
			filters:          []string{"some-filter"},
			filterMode:       filterModeBlackList,
			expectedLogEntry: nil,
			expectedResult:   true,
		},
	}

	logHook := test.NewLocal(logrus.StandardLogger())

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			result := ValidateFlags(testData.searchRoot, testData.filterMode, testData.filters)

			assert.Exactly(t, testData.expectedResult, result)
			if testData.expectedLogEntry != nil {
				assert.Exactly(t, testData.expectedLogEntry.Level, logHook.LastEntry().Level)
				assert.Exactly(t, testData.expectedLogEntry.Message, logHook.LastEntry().Message)
			}
		})
	}
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
