package logrushooks_test

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/krostar/logrushooks"
)

func BenchmarkLogWithoutCallerHook(b *testing.B) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("TEST")
	}
}

func BenchmarkLogWithFieldsHook(b *testing.B) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger.AddHook(logrushooks.NewFieldsHook(&logrushooks.FieldsOptions{
		Fields: map[string]interface{}{
			"answer": 42,
			"foo":    "var",
		},
	}))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("TEST")
	}
}

func BenchmarkLogWithCallerHook(b *testing.B) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger.AddHook(logrushooks.NewCallerHook(&logrushooks.CallerOptions{
		AppPackage: "github.com/krostar/domain-gen",
		CallerKey:  "caller",
	}))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("TEST")
	}
}

func BenchmarkLogWithAllHooks(b *testing.B) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	logger.AddHook(logrushooks.NewCallerHook(&logrushooks.CallerOptions{
		AppPackage: "github.com/krostar/domain-gen",
		CallerKey:  "caller",
	}))
	logger.AddHook(logrushooks.NewFieldsHook(&logrushooks.FieldsOptions{
		Fields: map[string]interface{}{
			"answer": 42,
			"foo":    "var",
		},
	}))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("TEST")
	}
}
