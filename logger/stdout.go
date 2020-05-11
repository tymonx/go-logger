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
	"sync"
)

// A Stdout represents a log handler object for logging message to standard
// output
type Stdout struct {
	mutex        sync.RWMutex
	formatter    *Formatter
	minimumLevel int
	maximumLevel int
}

// NewStdout created a new Stdout log handler object
func NewStdout() *Stdout {
	return &Stdout{
		formatter:    NewFormatter(),
		minimumLevel: MinimumLevel,
		maximumLevel: WarningLevel,
	}
}

// init registers Stdout log handler
func init() {
	RegisterHandler("stdout", func() Handler {
		return NewStdout()
	})
}

// SetFormatter sets Formatter
func (stdout *Stdout) SetFormatter(formatter *Formatter) Handler {
	stdout.mutex.Lock()
	defer stdout.mutex.Unlock()

	stdout.formatter = formatter

	return stdout
}

// GetFormatter returns Formatter
func (stdout *Stdout) GetFormatter() *Formatter {
	stdout.mutex.RLock()
	defer stdout.mutex.RUnlock()

	return stdout.formatter
}

// SetMinimumLevel sets minimum log level
func (stdout *Stdout) SetMinimumLevel(level int) Handler {
	stdout.mutex.Lock()
	defer stdout.mutex.Unlock()

	stdout.minimumLevel = level

	return stdout
}

// GetMinimumLevel returns minimum log level
func (stdout *Stdout) GetMinimumLevel() int {
	stdout.mutex.RLock()
	defer stdout.mutex.RUnlock()

	return stdout.minimumLevel
}

// SetMaximumLevel sets maximum log level
func (stdout *Stdout) SetMaximumLevel(level int) Handler {
	stdout.mutex.Lock()
	defer stdout.mutex.Unlock()

	stdout.maximumLevel = level

	return stdout
}

// GetMaximumLevel returns maximum log level
func (stdout *Stdout) GetMaximumLevel() int {
	stdout.mutex.RLock()
	defer stdout.mutex.RUnlock()

	return stdout.maximumLevel
}

// SetLevelRange sets minimum and maximum log level values
func (stdout *Stdout) SetLevelRange(min int, max int) Handler {
	stdout.mutex.Lock()
	defer stdout.mutex.Unlock()

	stdout.minimumLevel = min
	stdout.maximumLevel = max

	return stdout
}

// GetLevelRange returns minimum and maximum log level values
func (stdout *Stdout) GetLevelRange() (min int, max int) {
	stdout.mutex.RLock()
	defer stdout.mutex.RUnlock()

	return stdout.minimumLevel, stdout.maximumLevel
}

// Emit logs messages from logger to standard output
func (stdout *Stdout) Emit(record *Record) error {
	message, err := stdout.formatter.Format(record)

	if err != nil {
		return NewRuntimeError("cannot format record", err)
	}

	_, err = fmt.Fprintln(os.Stdout, message)

	if err != nil {
		return NewRuntimeError("cannot print to stdout", err)
	}

	return nil
}

// Close does nothing
func (stdout *Stdout) Close() error {
	return nil
}
