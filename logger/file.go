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

// These constants define default values for File log handler
const (
	DefaultFileName  = "log"
	DefaultFileMode  = 0644
	DefaultFileFlags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
)

// A File represents a log handler object for logging messages to file
type File struct {
	name      string
	mode      os.FileMode
	flags     int
	handler   *os.File
	formatter *Formatter
}

// NewFile creates a new File log handler object
func NewFile() *File {
	return &File{
		name:      DefaultFileName,
		mode:      DefaultFileMode,
		flags:     DefaultFileFlags,
		formatter: NewFormatter(),
	}
}

// init registers File log handler
func init() {
	RegisterHandler("file", func() Handler {
		return NewFile()
	})
}

// GetLevelRange returns minimum and maximum log level values
func (file *File) GetLevelRange() (min int, max int) {
	return TraceLevel, PanicLevel
}

// SetName sets file name used for log messages
func (file *File) SetName(name string) *File {
	file.name = name
	return file
}

// GetName sets file name used for log messages
func (file *File) GetName() string {
	return file.name
}

// SetFlags sets file flags from os package
func (file *File) SetFlags(flags int) *File {
	file.flags = flags
	return file
}

// GetFlags returns file flags
func (file *File) GetFlags() int {
	return file.flags
}

// SetMode sets file mode/permissions
func (file *File) SetMode(mode os.FileMode) *File {
	file.mode = mode
	return file
}

// GetMode returns file mode/permissions
func (file *File) GetMode() os.FileMode {
	return file.mode
}

// Emit logs messages from Logger to file
func (file *File) Emit(record *Record) error {
	if file.handler == nil {
		var err error

		file.handler, err = os.OpenFile(
			file.name,
			file.flags,
			file.mode,
		)

		if err != nil {
			file.handler = nil
			return NewRuntimeError("cannot open file", err)
		}
	}

	if file.handler != nil {
		_, err := fmt.Fprintln(file.handler, file.formatter.Format(record))

		if err != nil {
			return NewRuntimeError("cannot append file", err)
		}
	}

	return nil
}

// Close closes opened file
func (file *File) Close() error {
	if file.handler != nil {
		err := file.handler.Close()

		file.handler = nil

		if err != nil {
			return NewRuntimeError("cannot close file", err)
		}
	}

	return nil
}
