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
	"os"
	"path/filepath"
)

// RuntimeError defines runtime error with returned error message, file name,
// file line number and function name
type RuntimeError struct {
	line     int
	file     string
	message  string
	function string
	err      error
}

// NewRuntimeError creates new RuntimeError object
func NewRuntimeError(message string, err error) *RuntimeError {
	path, line, function := getPathLineFunction(2)

	return &RuntimeError{
		line:     line,
		file:     filepath.Base(path),
		message:  message,
		function: filepath.Base(function),
		err:      err,
	}
}

// Error returns formatted error string with message, file name, file line
// number and function name
func (runtimeError *RuntimeError) Error() string {
	message := fmt.Sprintf("Logger error %s:%d:%s(): %s",
		runtimeError.file,
		runtimeError.line,
		runtimeError.function,
		runtimeError.message,
	)

	if runtimeError.err != nil {
		message = fmt.Sprintf("%s: %v", message, runtimeError.err)
	}

	return message
}

// Unwrap wrapped error
func (runtimeError *RuntimeError) Unwrap() error {
	return runtimeError.err
}

// Print prints to error output formatted error with message, file name, file
// line number and function name
func (runtimeError *RuntimeError) Print() *RuntimeError {
	fmt.Fprintln(os.Stderr, runtimeError.Error())
	return runtimeError
}
