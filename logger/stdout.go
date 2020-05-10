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

// A Stdout represents a log handler object for logging message to standard
// output
type Stdout struct {
	formatter *Formatter
}

// NewStdout created a new Stdout log handler object
func NewStdout() *Stdout {
	return &Stdout{
		formatter: NewFormatter(),
	}
}

// init registers Stdout log handler
func init() {
	RegisterHandler("stdout", func() Handler {
		return NewStdout()
	})
}

// GetLevelRange returns minimum and maximum log level values
func (stdout *Stdout) GetLevelRange() (min int, max int) {
	return TraceLevel, WarningLevel
}

// Emit logs messages from logger to standard output
func (stdout *Stdout) Emit(record *Record) error {
	_, err := fmt.Fprintln(os.Stdout, stdout.formatter.Format(record))

	if err != nil {
		return NewRuntimeError("cannot print to stdout", err)
	}

	return nil
}

// Close does nothing
func (stdout *Stdout) Close() error {
	return nil
}
