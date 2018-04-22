package logrushooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func initLogger(callerOptions *CallerOptions) (*logrus.Logger, *bytes.Buffer) {
	reader := new(bytes.Buffer)
	logger := &logrus.Logger{
		Out:       reader,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger.AddHook(NewCallerHook(callerOptions))
	return logger, reader
}

func TestCallerWithDefaultOptions(t *testing.T) {
	// initialization
	logger, reader := initLogger(DefaultCallerOptions)

	// expectations
	// /!\ watch out !! There is currently 5 lines difference between caller and test
	_, expectedCallerFile, expectedCallerLine, ok := runtime.Caller(0)
	require.True(t, ok)
	expectedCaller := fmt.Sprintf("%s:%d", expectedCallerFile, expectedCallerLine+5)

	// execution
	logger.Infoln("This is a test")
	var result struct {
		Caller string `json:"caller"`
	}
	err := json.Unmarshal(reader.Bytes(), &result)
	require.NoError(t, err)

	// assertions
	require.Equal(t, expectedCaller, result.Caller)
}

func TestCallerWithCustomOptions(t *testing.T) {
	// initialization
	appPackage := "logrushooks"
	logger, reader := initLogger(&CallerOptions{
		AppPackage: appPackage,
		CallerKey:  "test",
	})

	// expectations
	// /!\ watch out !! There is currently 7 lines difference between caller and test
	_, expectedCallerFile, expectedCallerLine, ok := runtime.Caller(0)
	require.True(t, ok)
	pos := strings.Index(expectedCallerFile, appPackage)
	expectedCallerFile = expectedCallerFile[pos+1+len(appPackage):]
	expectedCaller := fmt.Sprintf("%s:%d", expectedCallerFile, expectedCallerLine+7)

	// execution
	logger.Infoln("This is a test")
	var result struct {
		Caller string `json:"test"`
	}
	err := json.Unmarshal(reader.Bytes(), &result)
	require.NoError(t, err)

	// assertions
	require.Equal(t, expectedCaller, result.Caller)
}

func TestCallerWithBadCustomOptions(t *testing.T) {
	// initialization
	options := &CallerOptions{
		AppPackage:      "appPackage",
		CallerKey:       "",
		PackagesToSkip:  []string{"titi", "toto"},
		CallerSkipStart: 0,
	}

	// expectations
	expectedOptions := &CallerOptions{
		AppPackage:      "appPackage",
		CallerKey:       DefaultCallerOptions.CallerKey,
		PackagesToSkip:  append(DefaultCallerOptions.PackagesToSkip, "titi", "toto"),
		CallerSkipStart: DefaultCallerOptions.CallerSkipStart,
	}

	// executions
	_ = NewCallerHook(options)

	// assertions
	require.Equal(t, expectedOptions, options)
}

func TestCallerWithCallerError(t *testing.T) {
	// initialization
	logger, reader := initLogger(&CallerOptions{
		CallerSkipStart: 100, // impossible value for the runtime.Caller to resolve
	})

	// expectations
	expectedCaller := "unable to find caller: runtime.Caller failed"

	// execution
	logger.Infoln("This is a test")
	var result struct {
		Caller string `json:"caller"`
	}
	err := json.Unmarshal(reader.Bytes(), &result)
	require.NoError(t, err)

	// assertions
	require.Equal(t, expectedCaller, result.Caller)
}
