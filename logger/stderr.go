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

// A Stderr represents a log handler object for logging message to error
// output
type Stderr struct {
	mutex        sync.RWMutex
	formatter    *Formatter
	minimumLevel int
	maximumLevel int
}

// NewStderr created a new Stderr log handler object
func NewStderr() *Stderr {
	return &Stderr{
		formatter:    NewFormatter(),
		minimumLevel: ErrorLevel,
		maximumLevel: MaximumLevel,
	}
}

// init registers Stderr log handler
func init() {
	RegisterHandler("stderr", func() Handler {
		return NewStderr()
	})
}

// SetFormatter sets Formatter
func (stderr *Stderr) SetFormatter(formatter *Formatter) Handler {
	stderr.mutex.Lock()
	defer stderr.mutex.Unlock()

	stderr.formatter = formatter

	return stderr
}

// GetFormatter returns Formatter
func (stderr *Stderr) GetFormatter() *Formatter {
	stderr.mutex.RLock()
	defer stderr.mutex.RUnlock()

	return stderr.formatter
}

// SetMinimumLevel sets minimum log level
func (stderr *Stderr) SetMinimumLevel(level int) Handler {
	stderr.mutex.Lock()
	defer stderr.mutex.Unlock()

	stderr.minimumLevel = level

	return stderr
}

// GetMinimumLevel returns minimum log level
func (stderr *Stderr) GetMinimumLevel() int {
	stderr.mutex.RLock()
	defer stderr.mutex.RUnlock()

	return stderr.minimumLevel
}

// SetMaximumLevel sets maximum log level
func (stderr *Stderr) SetMaximumLevel(level int) Handler {
	stderr.mutex.Lock()
	defer stderr.mutex.Unlock()

	stderr.maximumLevel = level

	return stderr
}

// GetMaximumLevel returns maximum log level
func (stderr *Stderr) GetMaximumLevel() int {
	stderr.mutex.RLock()
	defer stderr.mutex.RUnlock()

	return stderr.maximumLevel
}

// SetLevelRange sets minimum and maximum log level values
func (stderr *Stderr) SetLevelRange(min int, max int) Handler {
	stderr.mutex.Lock()
	defer stderr.mutex.Unlock()

	stderr.minimumLevel = min
	stderr.maximumLevel = max

	return stderr
}

// GetLevelRange returns minimum and maximum log level values
func (stderr *Stderr) GetLevelRange() (min int, max int) {
	stderr.mutex.RLock()
	defer stderr.mutex.RUnlock()

	return stderr.minimumLevel, stderr.maximumLevel
}

// Emit logs messages from logger to error output
func (stderr *Stderr) Emit(record *Record) error {
	message, err := stderr.formatter.Format(record)

	if err != nil {
		return NewRuntimeError("cannot format record", err)
	}

	_, err = fmt.Fprintln(os.Stderr, message)

	if err != nil {
		return NewRuntimeError("cannot print to stderr", err)
	}

	return nil
}

// Close does nothing
func (stderr *Stderr) Close() error {
	return nil
}
