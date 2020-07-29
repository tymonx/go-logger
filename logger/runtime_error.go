// Copyright 2020 Tymoteusz Blazejczyk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// These constants are used for the RuntimeError.
const (
	RuntimeErrorSkipCall = 1
)

// RuntimeError defines runtime error with returned error message, file name,
// file line number and function name.
type RuntimeError struct {
	line      int
	file      string
	message   string
	function  string
	arguments []interface{}
}

// NewRuntimeError creates new RuntimeError object.
func NewRuntimeError(message string, arguments ...interface{}) *RuntimeError {
	return NewRuntimeErrorBase(RuntimeErrorSkipCall, message, arguments...)
}

// NewRuntimeErrorBase creates new RuntimeError object using custom skip call value.
func NewRuntimeErrorBase(skipCall int, message string, arguments ...interface{}) *RuntimeError {
	pc, path, line, _ := runtime.Caller(skipCall + 1)

	return &RuntimeError{
		line:      line,
		file:      filepath.Base(path),
		message:   message,
		function:  filepath.Base(runtime.FuncForPC(pc).Name()),
		arguments: arguments,
	}
}

// Error returns formatted error string with message, file name, file line
// number and function name.
func (r *RuntimeError) Error() string {
	var formatted string

	var err error

	record := &Record{
		Message:   r.message,
		Arguments: r.arguments,
	}

	if formatted, err = NewFormatter().FormatMessage(record); err != nil {
		formatted = r.message
	}

	return fmt.Sprintf("%s:%d:%s(): %s",
		r.file,
		r.line,
		r.function,
		formatted,
	)
}

// Unwrap wrapped error.
func (r *RuntimeError) Unwrap() error {
	for _, argument := range r.arguments {
		if err, ok := argument.(error); ok {
			return err
		}
	}

	return nil
}
