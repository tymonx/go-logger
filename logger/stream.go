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
)

// A Stream represents a log handler object for logging messages using stream
// object
type Stream struct {
	writer    io.WriteCloser
	formatter *Formatter
}

// NewStream creates a new Stream log handler object
func NewStream() *Stream {
	return &Stream{
		formatter: NewFormatter(),
	}
}

// init registers Stream log handler
func init() {
	RegisterHandler("stream", func() Handler {
		return NewStream()
	})
}

// GetLevelRange returns minimum and maximum log level values
func (stream *Stream) GetLevelRange() (min int, max int) {
	return TraceLevel, PanicLevel
}

// SetWriter sets I/O writer object used for writing log messages
func (stream *Stream) SetWriter(writer io.WriteCloser) *Stream {
	stream.writer = writer
	return stream
}

// Emit logs messages from logger using I/O stream
func (stream *Stream) Emit(record *Record) error {
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
	if stream.writer != nil {
		err := stream.writer.Close()

		stream.writer = nil

		if err != nil {
			return NewRuntimeError("cannot close stream", err)
		}
	}

	return nil
}
