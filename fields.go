package logrushooks

import (
	"github.com/sirupsen/logrus"
)

type (
	// FieldsOptions stores options to configure the hook Fields.
	FieldsOptions struct {
		Fields map[string]interface{}
	}

	// Fields adds all Fields to the logrus fields.
	Fields struct {
		fields map[string]interface{}
	}
)

var (
	// DefaultFieldsOptions stores the default options for the Fields hook.
	DefaultFieldsOptions = &FieldsOptions{
		Fields: map[string]interface{}{},
	}
)

// NewFieldsHook creates a new Fields hook based on the options.
func NewFieldsHook(opts *FieldsOptions) *Fields {
	return &Fields{
		fields: opts.Fields,
	}
}

// Levels returns the list of available log levels
func (hook Fields) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is the code executed when all the hooks are fired
func (hook Fields) Fire(entry *logrus.Entry) error {
	for field, value := range hook.fields {
		entry.Data[field] = value
	}
	return nil
}
