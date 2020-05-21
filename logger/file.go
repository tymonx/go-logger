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
	"io"
	"os"
)

// These constants define default values for File log handler.
const (
	DefaultFileName  = "log"
	DefaultFileMode  = 0644
	DefaultFileFlags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
)

// A File represents a log handler object for logging messages to file.
type File struct {
	name   string
	stream *Stream
	flags  int
	mode   os.FileMode
}

// NewFile creates a new File log handler object.
func NewFile() *File {
	f := &File{
		name:   DefaultFileName,
		mode:   DefaultFileMode,
		flags:  DefaultFileFlags,
		stream: NewStream(),
	}

	f.stream.SetOpener(f)

	return f
}

// Open file.
func (f *File) Open() (io.WriteCloser, error) {
	return os.OpenFile(f.name, f.flags, f.mode)
}

// Enable enables log handler.
func (f *File) Enable() Handler {
	return f.stream.Enable()
}

// Disable disabled log handler.
func (f *File) Disable() Handler {
	return f.stream.Disable()
}

// IsEnabled returns if log handler is enabled.
func (f *File) IsEnabled() bool {
	return f.stream.IsEnabled()
}

// SetFormatter sets Formatter.
func (f *File) SetFormatter(formatter *Formatter) Handler {
	return f.stream.SetFormatter(formatter)
}

// GetFormatter returns Formatter.
func (f *File) GetFormatter() *Formatter {
	return f.stream.GetFormatter()
}

// SetLevel sets log level.
func (f *File) SetLevel(level int) Handler {
	return f.stream.SetLevel(level)
}

// SetMinimumLevel sets minimum log level.
func (f *File) SetMinimumLevel(level int) Handler {
	return f.stream.SetMinimumLevel(level)
}

// GetMinimumLevel returns minimum log level.
func (f *File) GetMinimumLevel() int {
	return f.stream.GetMinimumLevel()
}

// SetMaximumLevel sets maximum log level.
func (f *File) SetMaximumLevel(level int) Handler {
	return f.stream.SetMaximumLevel(level)
}

// GetMaximumLevel returns maximum log level.
func (f *File) GetMaximumLevel() int {
	return f.stream.GetMaximumLevel()
}

// SetLevelRange sets minimum and maximum log level values.
func (f *File) SetLevelRange(min, max int) Handler {
	return f.stream.SetLevelRange(min, max)
}

// GetLevelRange returns minimum and maximum log level values.
func (f *File) GetLevelRange() (min, max int) {
	return f.stream.GetLevelRange()
}

// SetName sets file name used for log messages.
func (f *File) SetName(name string) *File {
	f.stream.Lock()
	defer f.stream.Unlock()

	if f.name != name {
		f.name = name
		f.stream.Reopen()
	}

	return f
}

// GetName sets file name used for log messages.
func (f *File) GetName() string {
	f.stream.RLock()
	defer f.stream.RUnlock()

	return f.name
}

// SetFlags sets file flags from os package.
func (f *File) SetFlags(flags int) *File {
	f.stream.Lock()
	defer f.stream.Unlock()

	if f.flags != flags {
		f.flags = flags
		f.stream.Reopen()
	}

	return f
}

// GetFlags returns file flags.
func (f *File) GetFlags() int {
	f.stream.RLock()
	defer f.stream.RUnlock()

	return f.flags
}

// SetMode sets file mode/permissions.
func (f *File) SetMode(mode os.FileMode) *File {
	f.stream.Lock()
	defer f.stream.Unlock()

	if f.mode != mode {
		f.mode = mode
		f.stream.Reopen()
	}

	return f
}

// GetMode returns file mode/permissions.
func (f *File) GetMode() os.FileMode {
	f.stream.RLock()
	defer f.stream.RUnlock()

	return f.mode
}

// Emit logs messages from Logger to file.
func (f *File) Emit(record *Record) error {
	return f.stream.Emit(record)
}

// Close closes opened file.
func (f *File) Close() error {
	return f.stream.Close()
}
