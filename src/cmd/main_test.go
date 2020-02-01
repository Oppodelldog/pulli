package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/Oppodelldog/pulli/src/pulli"
)

func restoreOriginalFuncs() func() {
	originalOsArgsProvider := osArgsProvider
	originalBuildCommand := buildCommand
	originalPullAll := pullAll
	originalPulliBuildCommandFunc := pulliBuildCommandFunc
	originalPulliPullAllFunc := pulliPullAllFunc
	originalPulliValidateFlags := pulliValidateFlags
	originalExitProgram := exitProgram

	return func() {
		osArgsProvider = originalOsArgsProvider
		buildCommand = originalBuildCommand
		pullAll = originalPullAll
		pulliBuildCommandFunc = originalPulliBuildCommandFunc
		pulliPullAllFunc = originalPulliPullAllFunc
		pulliValidateFlags = originalPulliValidateFlags
		exitProgram = originalExitProgram
	}
}

func TestMainCallsPulli(t *testing.T) {
	defer restoreOriginalFuncs()()

	wasPullAllCalled := false

	osArgsProvider = func() []string {
		// if no is given, the main command will be called
		return []string{pulli.CommandName}
	}
	pullAll = func(args []string) {
		wasPullAllCalled = true
	}

	main()

	if !wasPullAllCalled {
		t.Fatalf("expected pullAll was called, but it was not")
	}
}

func TestMainCallsBuildCommand(t *testing.T) {
	defer restoreOriginalFuncs()()

	wasBuildCommandCalled := false

	osArgsProvider = func() []string {
		return []string{pulli.CommandName, pulli.SubCommandNameBuildCommand}
	}

	buildCommand = func(args []string) {
		wasBuildCommandCalled = true
	}

	main()

	if !wasBuildCommandCalled {
		t.Fatalf("expected buildCommand was called, but it was not")
	}
}

type pulliPullAllFuncMock struct {
	pulliFuncWasCalled bool
	parmSearchRoot     string
	parmFilters        []string
	parmFilterMode     string
}

func (m *pulliPullAllFuncMock) pullAll(searchRoot string, filters []string, filterMode string) {
	m.parmSearchRoot = searchRoot
	m.parmFilters = filters
	m.parmFilterMode = filterMode
	m.pulliFuncWasCalled = true
}

func TestPullAll(t *testing.T) {
	defer restoreOriginalFuncs()()

	testCases := map[string]struct {
		inputArgs          []string
		expectedSearchRoot string
		expectedFilters    []string
		expectedFilterMode string
	}{
		"no args": {inputArgs: []string{pulli.CommandName}, expectedSearchRoot: "."},
		"-dir": {
			inputArgs:          []string{pulli.CommandName, "-dir=some-dir"},
			expectedSearchRoot: "some-dir"},
		"-filterMode=whitelist": {
			inputArgs:          []string{pulli.CommandName, "-" + pulli.ArgNameFilterMode + "=" + pulli.FilterModeWhiteList},
			expectedSearchRoot: ".",
			expectedFilterMode: pulli.FilterModeWhiteList},
		"-filterMode=blacklist": {
			inputArgs:          []string{pulli.CommandName, "-" + pulli.ArgNameFilterMode + "=" + pulli.FilterModeBlackList},
			expectedSearchRoot: ".",
			expectedFilterMode: pulli.FilterModeBlackList},
		"-filter=test1": {
			inputArgs:          []string{pulli.CommandName, "-" + pulli.ArgNameFilter + "=test1"},
			expectedSearchRoot: ".",
			expectedFilters:    []string{"test1"}},
		"-filter=test1 -filter=test2": {
			inputArgs:          []string{pulli.CommandName, "-" + pulli.ArgNameFilter + "=test1", "-" + pulli.ArgNameFilter + "=test2"},
			expectedSearchRoot: ".",
			expectedFilters:    []string{"test1", "test2"}},
		"all parameters": {
			inputArgs: []string{
				pulli.CommandName,
				"-dir=some-dir",
				"-" + pulli.ArgNameFilterMode + "=" + pulli.FilterModeBlackList,
				"-" + pulli.ArgNameFilter + "=test1",
				"-" + pulli.ArgNameFilter + "=test2"},
			expectedSearchRoot: "some-dir",
			expectedFilterMode: pulli.FilterModeBlackList,
			expectedFilters:    []string{"test1", "test2"}},
	}

	pulliValidateFlags = func(searchRoot string, filterMode string, filters []string) bool {
		return true
	}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			mock := new(pulliPullAllFuncMock)
			pulliPullAllFunc = mock.pullAll
			pullAll(testData.inputArgs)

			if testData.expectedFilterMode != mock.parmFilterMode {
				t.Fatalf("expected pullAll was called with filterMode='%s', but was '%s'", testData.expectedFilterMode, mock.parmFilterMode)
			}
			if testData.expectedSearchRoot != mock.parmSearchRoot {
				t.Fatalf("expected pullAll was called with searchRoot='%s', but was '%s'", testData.expectedSearchRoot, mock.parmSearchRoot)
			}
			if !reflect.DeepEqual(testData.expectedFilters, mock.parmFilters) {
				t.Fatalf("expected pullAll was called with filters='%#v', but was '%#v'", testData.expectedFilters, mock.parmFilters)
			}
		})
	}
}

type pulliValidateMock struct {
	pulliFuncWasCalled bool
	parmSearchRoot     string
	parmFilters        []string
	parmFilterMode     string
	returnValue        bool
}

func (m *pulliValidateMock) validate(searchRoot string, filterMode string, filters []string) bool {
	m.parmSearchRoot = searchRoot
	m.parmFilters = filters
	m.parmFilterMode = filterMode
	m.pulliFuncWasCalled = true

	return m.returnValue
}

func TestPullAllCallsValidate(t *testing.T) {
	defer restoreOriginalFuncs()()

	testCases := map[string]struct {
		inputArgs          []string
		expectedSearchRoot string
		expectedFilters    []string
		expectedFilterMode string
	}{
		"no args": {inputArgs: []string{pulli.CommandName}, expectedSearchRoot: "."},
		"all parameters": {
			inputArgs: []string{
				pulli.CommandName,
				"-dir=some-dir",
				"-" + pulli.ArgNameFilterMode + "=" + pulli.FilterModeBlackList,
				"-" + pulli.ArgNameFilter + "=test1",
				"-" + pulli.ArgNameFilter + "=test2"},
			expectedSearchRoot: "some-dir",
			expectedFilterMode: pulli.FilterModeBlackList,
			expectedFilters:    []string{"test1", "test2"}},
	}

	pulliPullAllFunc = func(searchRoot string, filters []string, filterMode string) {}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			mock := &pulliValidateMock{returnValue: true}
			pulliValidateFlags = mock.validate
			pullAll(testData.inputArgs)

			if testData.expectedFilterMode != mock.parmFilterMode {
				t.Fatalf("expected pullAll was called with filterMode='%s', but was '%s'", testData.expectedFilterMode, mock.parmFilterMode)
			}
			if testData.expectedSearchRoot != mock.parmSearchRoot {
				t.Fatalf("expected pullAll was called with searchRoot='%s', but was '%s'", testData.expectedSearchRoot, mock.parmSearchRoot)
			}
			if !reflect.DeepEqual(testData.expectedFilters, mock.parmFilters) {
				t.Fatalf("expected pullAll was called with filters='%#v', but was '%#v'", testData.expectedFilters, mock.parmFilters)
			}
		})
	}
}

type programExitMock struct {
	wasCalled bool
	exitCode  int
}

func (m *programExitMock) exit(code int) {
	m.wasCalled = true
	m.exitCode = code
}

func TestPullAllValidationFailsProgrammWillExit(t *testing.T) {
	defer restoreOriginalFuncs()()

	pullAllFuncMock := new(pulliPullAllFuncMock)
	exitCodeMock := new(programExitMock)
	exitProgram = exitCodeMock.exit

	pulliValidateFlags = (&pulliValidateMock{returnValue: false}).validate
	pulliPullAllFunc = pullAllFuncMock.pullAll

	args := []string{pulli.CommandName}
	pullAll(args)

	if pullAllFuncMock.pulliFuncWasCalled {
		t.Fatalf("validation failed so pulliPullAllFunc should not have been called")
	}
}

func TestPullAllValidationSucceedsPulliPullAllIsCalled(t *testing.T) {
	defer restoreOriginalFuncs()()

	pullAllFuncMock := new(pulliPullAllFuncMock)
	exitCodeMock := new(programExitMock)
	exitProgram = exitCodeMock.exit

	pulliValidateFlags = (&pulliValidateMock{returnValue: true}).validate
	pulliPullAllFunc = pullAllFuncMock.pullAll

	args := []string{pulli.CommandName}
	pullAll(args)

	if !pullAllFuncMock.pulliFuncWasCalled {
		t.Fatalf("validation succceed so pulliPullAllFunc should have been called")
	}
}

func TestPullUnknownArgumentLeadsToProgramExit(t *testing.T) {
	defer restoreOriginalFuncs()()

	pullAllFuncMock := new(pulliPullAllFuncMock)
	exitCodeMock := new(programExitMock)
	exitProgram = exitCodeMock.exit

	pulliValidateFlags = (&pulliValidateMock{returnValue: false}).validate
	pulliPullAllFunc = pullAllFuncMock.pullAll

	args := []string{pulli.CommandName, "-unknown"}
	pullAll(args)

	if !exitCodeMock.wasCalled {
		t.Fatalf("unknown argument program should have quit")
	}
	if pullAllFuncMock.pulliFuncWasCalled {
		t.Fatalf("unknown argument, pulliPullAllFunc should not have been called")
	}
}

func TestBuildCommand(t *testing.T) {
	defer restoreOriginalFuncs()()

	pulliFuncWasCalled := false
	pulliFuncSearchRootParameter := ""

	pulliBuildCommandFunc = func(searchRoot string) {
		pulliFuncSearchRootParameter = searchRoot
		pulliFuncWasCalled = true
	}

	args := []string{pulli.CommandName, pulli.SubCommandNameBuildCommand}

	buildCommand(args)

	if !pulliFuncWasCalled {
		t.Fatalf("expected pulliFuncWasCalled was called, but it was not")
	}
	if pulliFuncSearchRootParameter != "." {
		t.Fatalf("expected pulliFuncWasCalled was called with searchRoot='.', but it was '%s'", pulliFuncSearchRootParameter)
	}
}

func TestBuildCommandWithDirParameter(t *testing.T) {
	defer restoreOriginalFuncs()()

	pulliFuncWasCalled := false
	pulliFuncSearchRootParameter := ""

	pulliBuildCommandFunc = func(searchRoot string) {
		pulliFuncSearchRootParameter = searchRoot
		pulliFuncWasCalled = true
	}

	stubbedDir := "dirInputParameter"

	args := []string{pulli.CommandName, pulli.SubCommandNameBuildCommand, "-dir", stubbedDir}

	buildCommand(args)

	if !pulliFuncWasCalled {
		t.Fatalf("expected pulliFuncWasCalled was called, but it was not")
	}
	if pulliFuncSearchRootParameter != stubbedDir {
		t.Fatalf("expected pulliFuncWasCalled was called with searchRoot='%s', but it was '%s'", stubbedDir, pulliFuncSearchRootParameter)
	}
}

func TestBuildCommandExitsOnFlagError(t *testing.T) {
	defer restoreOriginalFuncs()()

	exitMock := &programExitMock{}
	exitProgram = exitMock.exit

	args := []string{pulli.CommandName, pulli.SubCommandNameBuildCommand, "-unknown"}

	buildCommand(args)

	if !exitMock.wasCalled {
		t.Fatalf("eprogram should have exited because of argument error")
	}
}

func TestOsArgsProviderProvidesOsArgs(t *testing.T) {
	if !reflect.DeepEqual(os.Args, osArgsProvider()) {
		t.FailNow()
	}
}
