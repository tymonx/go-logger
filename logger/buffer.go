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
	"bytes"
)

// A Buffer represents a log handler object for logging messages using buffer
// object.
type Buffer struct {
	buffer bytes.Buffer
	stream *Stream
}

// NewBuffer creates a new buffer log handler object.
func NewBuffer() *Buffer {
	b := &Buffer{
		stream: NewStream(),
	}

	b.stream.writer = &b.buffer

	return b
}

// SetStreamHandler sets custom stream handler.
func (b *Buffer) SetStreamHandler(handler StreamHandler) *Buffer {
	b.stream.SetStreamHandler(handler)
	return b
}

// GetBuffer returns internal buffer object.
func (b *Buffer) GetBuffer() *bytes.Buffer {
	b.stream.RLock()
	defer b.stream.RUnlock()

	return &b.buffer
}

// Enable enables log handler.
func (b *Buffer) Enable() Handler {
	return b.stream.Enable()
}

// Disable disabled log handler.
func (b *Buffer) Disable() Handler {
	return b.stream.Disable()
}

// IsEnabled returns if log handler is enabled.
func (b *Buffer) IsEnabled() bool {
	return b.stream.IsEnabled()
}

// SetFormatter sets log formatter.
func (b *Buffer) SetFormatter(formatter *Formatter) Handler {
	return b.stream.SetFormatter(formatter)
}

// GetFormatter returns log formatter.
func (b *Buffer) GetFormatter() *Formatter {
	return b.stream.GetFormatter()
}

// SetLevel sets log level.
func (b *Buffer) SetLevel(level int) Handler {
	return b.stream.SetLevel(level)
}

// SetMinimumLevel sets minimum log level.
func (b *Buffer) SetMinimumLevel(level int) Handler {
	return b.stream.SetMinimumLevel(level)
}

// GetMinimumLevel returns minimum log level.
func (b *Buffer) GetMinimumLevel() int {
	return b.stream.GetMinimumLevel()
}

// SetMaximumLevel sets maximum log level.
func (b *Buffer) SetMaximumLevel(level int) Handler {
	return b.stream.SetMaximumLevel(level)
}

// GetMaximumLevel returns maximum log level.
func (b *Buffer) GetMaximumLevel() int {
	return b.stream.GetMaximumLevel()
}

// SetLevelRange sets minimum and maximum log level values.
func (b *Buffer) SetLevelRange(min, max int) Handler {
	return b.stream.SetLevelRange(min, max)
}

// GetLevelRange returns minimum and maximum log level values.
func (b *Buffer) GetLevelRange() (min, max int) {
	return b.stream.GetLevelRange()
}

// Emit logs messages from logger using buffer.
func (b *Buffer) Emit(record *Record) error {
	return b.stream.Emit(record)
}

// Close closes buffer.
func (b *Buffer) Close() error {
	return b.stream.Close()
}
