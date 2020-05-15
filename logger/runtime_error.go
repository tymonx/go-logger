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
)

const (
	runtimeErrorSkipCall = 2
)

// RuntimeError defines runtime error with returned error message, file name,
// file line number and function name.
type RuntimeError struct {
	line     int
	file     string
	message  string
	function string
	err      error
}

// NewRuntimeError creates new RuntimeError object.
func NewRuntimeError(message string, err error) *RuntimeError {
	path, line, function := getPathLineFunction(runtimeErrorSkipCall)

	return &RuntimeError{
		line:     line,
		file:     filepath.Base(path),
		message:  message,
		function: filepath.Base(function),
		err:      err,
	}
}

// Error returns formatted error string with message, file name, file line
// number and function name.
func (r *RuntimeError) Error() string {
	message := fmt.Sprintf("%s:%d:%s(): %s",
		r.file,
		r.line,
		r.function,
		r.message,
	)

	if r.err != nil {
		message = fmt.Sprintf("%s: %v", message, r.err)
	}

	return message
}

// Unwrap wrapped error.
func (r *RuntimeError) Unwrap() error {
	return r.err
}
