package errors

import (
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"
)

// TraceableType type alias for string
type TraceableType string

const (
	NotFound     TraceableType = "NotFound"
	CreateError  TraceableType = "CreateError"
	UpdateError  TraceableType = "UpdateError"
	DeleteError  TraceableType = "DeleteError"
	BadRequest   TraceableType = "BadRequest"
	Unauthorized TraceableType = "Unauthorized"
	TypeMismatch TraceableType = "TypeMismatch"
)

// Traceable contains all the details and error stack
type Traceable struct {
	ErrorType TraceableType `json:"error_type"`
	Message   string        `json:"message"`
	Cause     error         `json:"cause"`
}

// NewTraceable constructor for Traceable error
func NewTraceable(errorType TraceableType, message string, originError error) *Traceable {
	err := &Traceable{
		ErrorType: errorType,
		Message:   message,
		Cause:     originError,
	}
	log.Error(err.Error())
	return err
}

// Error interface implementation
func (e Traceable) Error() string {
	out, _ := json.Marshal(e)
	return string(out)
}

// StackTrace print the stacktrace for the current error
func (e Traceable) StackTrace() string {
	if e.Cause != nil {
		if reflect.TypeOf(e.Cause) == reflect.TypeOf(&Traceable{}) {
			(e.Cause.(Traceable)).StackTrace()
		}
		log.Error(e.Error())
	}
	out, _ := json.Marshal(e)
	return string(out)
}
