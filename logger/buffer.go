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
	"fmt"
	"sync"
)

// A Buffer represents a log handler object for logging messages using buffer
// object.
type Buffer struct {
	buffer       bytes.Buffer
	formatter    *Formatter
	mutex        sync.RWMutex
	minimumLevel int
	maximumLevel int
	isDisabled   bool
}

// NewBuffer creates a new buffer log handler object.
func NewBuffer() *Buffer {
	return &Buffer{
		formatter:    NewFormatter(),
		minimumLevel: MinimumLevel,
		maximumLevel: MaximumLevel,
	}
}

// GetBuffer returns internal buffer object.
func (b *Buffer) GetBuffer() *bytes.Buffer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return &b.buffer
}

// Enable enables log handler.
func (b *Buffer) Enable() Handler {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.isDisabled = false

	return b
}

// Disable disabled log handler.
func (b *Buffer) Disable() Handler {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.isDisabled = true

	return b
}

// IsEnabled returns if log handler is enabled.
func (b *Buffer) IsEnabled() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return !b.isDisabled
}

// SetFormatter sets log formatter.
func (b *Buffer) SetFormatter(formatter *Formatter) Handler {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.formatter = formatter

	return b
}

// GetFormatter returns log formatter.
func (b *Buffer) GetFormatter() *Formatter {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.formatter
}

// SetLevel sets log level.
func (b *Buffer) SetLevel(level int) Handler {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.minimumLevel = level
	b.maximumLevel = level

	return b
}

// SetMinimumLevel sets minimum log level.
func (b *Buffer) SetMinimumLevel(level int) Handler {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.minimumLevel = level

	return b
}

// GetMinimumLevel returns minimum log level.
func (b *Buffer) GetMinimumLevel() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.minimumLevel
}

// SetMaximumLevel sets maximum log level.
func (b *Buffer) SetMaximumLevel(level int) Handler {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.maximumLevel = level

	return b
}

// GetMaximumLevel returns maximum log level.
func (b *Buffer) GetMaximumLevel() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.maximumLevel
}

// SetLevelRange sets minimum and maximum log level values.
func (b *Buffer) SetLevelRange(min, max int) Handler {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.minimumLevel = min
	b.maximumLevel = max

	return b
}

// GetLevelRange returns minimum and maximum log level values.
func (b *Buffer) GetLevelRange() (min, max int) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.minimumLevel, b.maximumLevel
}

// Emit logs messages from logger using buffer.
func (b *Buffer) Emit(record *Record) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	message, err := b.formatter.Format(record)

	if err != nil {
		return NewRuntimeError("cannot format record", err)
	}

	_, err = fmt.Fprintln(&b.buffer, message)

	if err != nil {
		return NewRuntimeError("cannot write to buffer", err)
	}

	return nil
}

// Close closes buffer.
func (b *Buffer) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.buffer.Reset()

	return nil
}
