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
)

// A Stderr represents a log handler object for logging message to error
// output
type Stderr struct {
	formatter *Formatter
}

// NewStderr created a new Stderr log handler object
func NewStderr() *Stderr {
	return &Stderr{
		formatter: NewFormatter(),
	}
}

// init registers Stderr log handler
func init() {
	RegisterHandler("stderr", func() Handler {
		return NewStderr()
	})
}

// GetLevelRange returns minimum and maximum log level values
func (stderr *Stderr) GetLevelRange() (min int, max int) {
	return ErrorLevel, PanicLevel
}

// Emit logs messages from logger to error output
func (stderr *Stderr) Emit(record *Record) error {
	_, err := fmt.Fprintln(os.Stderr, stderr.formatter.Format(record))

	if err != nil {
		return NewRuntimeError("cannot print to stderr", err)
	}

	return nil
}

// Close does nothing
func (stderr *Stderr) Close() error {
	return nil
}
