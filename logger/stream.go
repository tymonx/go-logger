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
	"io"
	"sync"
)

// A Stream represents a log handler object for logging messages using stream
// object
type Stream struct {
	writer       io.WriteCloser
	formatter    *Formatter
	mutex        sync.RWMutex
	minimumLevel int
	maximumLevel int
}

// NewStream creates a new Stream log handler object
func NewStream() *Stream {
	return &Stream{
		formatter:    NewFormatter(),
		minimumLevel: MinimumLevel,
		maximumLevel: MaximumLevel,
	}
}

// init registers Stream log handler
func init() {
	RegisterHandler("stream", func() Handler {
		return NewStream()
	})
}

// SetFormatter sets Formatter
func (stream *Stream) SetFormatter(formatter *Formatter) Handler {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	stream.formatter = formatter

	return stream
}

// GetFormatter returns Formatter
func (stream *Stream) GetFormatter() *Formatter {
	stream.mutex.RLock()
	defer stream.mutex.RUnlock()

	return stream.formatter
}

// SetMinimumLevel sets minimum log level
func (stream *Stream) SetMinimumLevel(level int) Handler {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	stream.minimumLevel = level

	return stream
}

// GetMinimumLevel returns minimum log level
func (stream *Stream) GetMinimumLevel() int {
	stream.mutex.RLock()
	defer stream.mutex.RUnlock()

	return stream.minimumLevel
}

// SetMaximumLevel sets maximum log level
func (stream *Stream) SetMaximumLevel(level int) Handler {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	stream.maximumLevel = level

	return stream
}

// GetMaximumLevel returns maximum log level
func (stream *Stream) GetMaximumLevel() int {
	stream.mutex.RLock()
	defer stream.mutex.RUnlock()

	return stream.maximumLevel
}

// SetLevelRange sets minimum and maximum log level values
func (stream *Stream) SetLevelRange(min int, max int) Handler {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	stream.minimumLevel = min
	stream.maximumLevel = max

	return stream
}

// GetLevelRange returns minimum and maximum log level values
func (stream *Stream) GetLevelRange() (min int, max int) {
	stream.mutex.RLock()
	defer stream.mutex.RUnlock()

	return stream.minimumLevel, stream.maximumLevel
}

// SetWriter sets I/O writer object used for writing log messages
func (stream *Stream) SetWriter(writer io.WriteCloser) *Stream {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	if stream.writer != writer {
		stream.close()
		stream.writer = writer
	}

	return stream
}

// Emit logs messages from logger using I/O stream
func (stream *Stream) Emit(record *Record) error {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	if stream.writer != nil {
		_, err := fmt.Fprintln(stream.writer, stream.formatter.Format(record))

		if err != nil {
			return NewRuntimeError("cannot write to stream", err)
		}
	}

	return nil
}

// Close closes I/O stream
func (stream *Stream) Close() error {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	return stream.close()
}

func (stream *Stream) close() error {
	if stream.writer != nil {
		err := stream.writer.Close()

		stream.writer = nil

		if err != nil {
			return NewRuntimeError("cannot close stream", err)
		}
	}

	return nil
}
