package logrushooks

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func initLoggerWithFieldsHook(fieldsOptions *FieldsOptions) (*logrus.Logger, *bytes.Buffer) {
	reader := new(bytes.Buffer)
	logger := &logrus.Logger{
		Out:       reader,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger.AddHook(NewFieldsHook(fieldsOptions))
	return logger, reader
}

func TestFieldsWithDefaultOptions(t *testing.T) {
	// initialization
	logger, reader := initLoggerWithFieldsHook(DefaultFieldsOptions)

	// expectations
	type logEntry struct {
		Level   string    `json:"level"`
		Message string    `json:"msg"`
		Time    time.Time `json:"time"`
	}
	expectedResult := &logEntry{
		Level:   "info",
		Message: "This is a test",
		Time:    time.Now(),
	}

	// execution
	logger.Infoln("This is a test")
	result := &logEntry{}
	err := json.Unmarshal(reader.Bytes(), result)
	result.Time = expectedResult.Time
	require.NoError(t, err)

	// assertions
	require.Equal(t, expectedResult, result)
}

func TestFieldsWithCustomOptions(t *testing.T) {
	// initialization
	logger, reader := initLoggerWithFieldsHook(&FieldsOptions{
		Fields: map[string]interface{}{
			"field1": "valuefield1",
			"field2": "valuefield2",
		},
	})

	// expectations
	type logEntry struct {
		Level   string      `json:"level"`
		Message string      `json:"msg"`
		Time    time.Time   `json:"time"`
		Field1  interface{} `json:"field1"`
		Field2  interface{} `json:"field2"`
	}
	expectedResult := &logEntry{
		Level:   "info",
		Message: "This is a test",
		Time:    time.Now(),
		Field1:  "valuefield1",
		Field2:  "valuefield2",
	}

	// execution
	logger.Infoln("This is a test")
	result := &logEntry{}
	err := json.Unmarshal(reader.Bytes(), result)
	result.Time = expectedResult.Time
	require.NoError(t, err)

	// assertions
	require.Equal(t, expectedResult, result)
}
