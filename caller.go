package logrushooks

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	// CallerOptions stores options to configure the hook caller.
	CallerOptions struct {
		// AppPackage removes, if not empty and applicable, the AppPackage from the file path.
		AppPackage string
		// CallerKey is the name of the logrus'field key.
		CallerKey string
		// PackagesToSkip stores a list of packages that will be skipped from the runtime.Caller
		// you would usually add something unless you have a log pre-processor that would not
		// reflect the true caller.
		PackagesToSkip []string
		// CallerSkipStart is an optimization variable if you have a lot of stack you want to skip
		// to avoid calling runtime.Caller that much times.
		CallerSkipStart int
	}

	// Caller adds a "caller" field (field key can be configured, see CallerOptions)
	// to the logrus fields, containing the file and line of the logrus caller.
	Caller struct {
		appPackage      string
		callerKey       string
		packagesToSkip  []string
		callerSkipStart int
	}
)

var (
	// DefaultCallerOptions stores the default options for the caller hook.
	DefaultCallerOptions = &CallerOptions{
		AppPackage:      "",
		CallerKey:       "caller",
		PackagesToSkip:  []string{"github.com/sirupsen/logrus"},
		CallerSkipStart: 7,
	}

	// ErrCallerError defines a runtime.Caller error
	ErrCallerError = errors.New("runtime.Caller failed")
)

// NewCallerHook creates a new caller hook based on the options.
func NewCallerHook(opts *CallerOptions) *Caller {
	if opts.CallerKey == "" {
		opts.CallerKey = DefaultCallerOptions.CallerKey
	}
	if opts.CallerSkipStart == 0 {
		opts.CallerSkipStart = DefaultCallerOptions.CallerSkipStart
	}
	opts.PackagesToSkip = append(DefaultCallerOptions.PackagesToSkip, opts.PackagesToSkip...)

	return &Caller{
		appPackage:      opts.AppPackage,
		callerKey:       opts.CallerKey,
		callerSkipStart: opts.CallerSkipStart,
		packagesToSkip:  opts.PackagesToSkip,
	}
}

// Levels returns the list of available log levels
func (hook Caller) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is the code executed when all the hooks are fired
func (hook Caller) Fire(entry *logrus.Entry) error {
	callerFile, callerLine, err := hook.findInterestingCaller()
	if err != nil {
		entry.Data[hook.callerKey] = errors.Wrap(err, "unable to find caller")
		return nil
	}

	// ƒrom useless-path/appPackage/file-path we only want to keep file-path
	callerFile = hook.stripPath(callerFile)

	entry.Data[hook.callerKey] = fmt.Sprintf("%s:%d", callerFile, callerLine)
	return nil
}

func (hook *Caller) findInterestingCaller() (file string, line int, err error) {
	var done = false
	for callerSkip := hook.callerSkipStart + 1; !done; callerSkip++ {
		var ok bool
		if _, file, line, ok = runtime.Caller(callerSkip); !ok {
			return "", 0, ErrCallerError
		}

		done = true

		// we may want to do another loop to skip logger (or other) packages
		for _, packageToSkip := range hook.packagesToSkip {
			if strings.Contains(file, packageToSkip) {
				done = false
			}
		}
	}

	// if we're still here, we got what we wanted
	return file, line, nil
}

func (hook *Caller) stripPath(path string) (stripped string) {
	if pos := strings.Index(path, hook.appPackage); pos > 0 {
		stripped = path[pos+1+len(hook.appPackage):]
	} else {
		stripped = path
	}

	return stripped
}
