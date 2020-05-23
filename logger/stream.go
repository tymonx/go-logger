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

// StreamHandler defines a custom stream handler for writing log records with writer.
type StreamHandler func(writer io.Writer, record *Record, formatter *Formatter) error

// Opener implements Open method.
type Opener interface {
	Open() (io.WriteCloser, error)
}

// A Stream represents a log handler object for logging messages using stream
// object.
type Stream struct {
	writer       io.Writer
	closer       io.Closer
	formatter    *Formatter
	mutex        sync.RWMutex
	opener       Opener
	minimumLevel int
	maximumLevel int
	reopen       bool
	isDisabled   bool
	handler      StreamHandler
}

// NewStream creates a new Stream log handler object.
func NewStream() *Stream {
	return &Stream{
		formatter:    NewFormatter(),
		minimumLevel: MinimumLevel,
		maximumLevel: MaximumLevel,
		handler:      StreamHandlerDefault,
	}
}

// Lock locks stream.
func (s *Stream) Lock() {
	s.mutex.Lock()
}

// Unlock locks stream.
func (s *Stream) Unlock() {
	s.mutex.Unlock()
}

// RLock locks stream.
func (s *Stream) RLock() {
	s.mutex.RLock()
}

// RUnlock locks stream.
func (s *Stream) RUnlock() {
	s.mutex.RUnlock()
}

// SetStreamHandler sets custom stream handler.
func (s *Stream) SetStreamHandler(handler StreamHandler) *Stream {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if handler != nil {
		s.handler = handler
	}

	return s
}

// SetWriter sets new writer to stream.
func (s *Stream) SetWriter(writer io.Writer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.writer == writer {
		return nil
	}

	if s.closer != nil {
		err := s.closer.Close()

		if err != nil {
			return NewRuntimeError("cannot close stream", err)
		}
	}

	s.writer = writer
	s.closer = nil

	return nil
}

// SetWriteCloser sets new writer and closer to stream.
func (s *Stream) SetWriteCloser(writeCloser io.WriteCloser) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if (s.writer == writeCloser) && (s.closer == writeCloser) {
		return nil
	}

	if s.closer != nil {
		err := s.closer.Close()

		if err != nil {
			return NewRuntimeError("cannot close stream", err)
		}
	}

	s.writer = writeCloser
	s.closer = writeCloser

	return nil
}

// SetOpener sets opener interface.
func (s *Stream) SetOpener(opener Opener) *Stream {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.opener = opener

	return s
}

// Reopen reopens stream.
func (s *Stream) Reopen() *Stream {
	s.reopen = true

	return s
}

// Enable enables log handler.
func (s *Stream) Enable() Handler {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.isDisabled = false

	return s
}

// Disable disabled log handler.
func (s *Stream) Disable() Handler {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.isDisabled = true

	return s
}

// IsEnabled returns if log handler is enabled.
func (s *Stream) IsEnabled() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return !s.isDisabled
}

// SetFormatter sets Formatter.
func (s *Stream) SetFormatter(formatter *Formatter) Handler {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.formatter = formatter

	return s
}

// GetFormatter returns Formatter.
func (s *Stream) GetFormatter() *Formatter {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.formatter
}

// SetLevel sets log level.
func (s *Stream) SetLevel(level int) Handler {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.minimumLevel = level
	s.maximumLevel = level

	return s
}

// SetMinimumLevel sets minimum log level.
func (s *Stream) SetMinimumLevel(level int) Handler {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.minimumLevel = level

	return s
}

// GetMinimumLevel returns minimum log level.
func (s *Stream) GetMinimumLevel() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.minimumLevel
}

// SetMaximumLevel sets maximum log level.
func (s *Stream) SetMaximumLevel(level int) Handler {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.maximumLevel = level

	return s
}

// GetMaximumLevel returns maximum log level.
func (s *Stream) GetMaximumLevel() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.maximumLevel
}

// SetLevelRange sets minimum and maximum log level values.
func (s *Stream) SetLevelRange(min, max int) Handler {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.minimumLevel = min
	s.maximumLevel = max

	return s
}

// GetLevelRange returns minimum and maximum log level values.
func (s *Stream) GetLevelRange() (min, max int) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.minimumLevel, s.maximumLevel
}

// Emit logs messages from logger using I/O stream.
func (s *Stream) Emit(record *Record) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.reopen {
		if s.closer != nil {
			err := s.closer.Close()

			if err != nil {
				return NewRuntimeError("cannot close stream", err)
			}

			s.writer = nil
			s.closer = nil
		}

		s.reopen = false
	}

	if (s.writer == nil) && (s.closer == nil) && (s.opener != nil) {
		writer, err := s.opener.Open()

		if err != nil {
			return NewRuntimeError("cannot open stream", err)
		}

		s.writer = writer
		s.closer = writer
	}

	if s.writer != nil {
		if err := s.handler(s.writer, record, s.formatter); err != nil {
			return NewRuntimeError("cannot write to stream", err)
		}
	}

	return nil
}

// Close closes I/O stream.
func (s *Stream) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closer != nil {
		err := s.closer.Close()

		if err != nil {
			return NewRuntimeError("cannot close stream", err)
		}

		s.writer = nil
		s.closer = nil
	}

	return nil
}

// StreamHandlerDefault is a default stream handler for writing log records to stream.
func StreamHandlerDefault(writer io.Writer, record *Record, formatter *Formatter) error {
	message, err := formatter.Format(record)

	if err != nil {
		return NewRuntimeError("cannot format record", err)
	}

	if _, err := fmt.Fprintln(writer, message); err != nil {
		return NewRuntimeError("cannot write to stream", err)
	}

	return nil
}

// StreamHandlerNDJSON handles writing log records in the NDJSON format.
func StreamHandlerNDJSON(writer io.Writer, record *Record, _ *Formatter) error {
	bytes, err := record.ToJSON()

	if err != nil {
		return NewRuntimeError("cannot format record", err)
	}

	if _, err := fmt.Fprintln(writer, string(bytes)); err != nil {
		return NewRuntimeError("cannot write to stream", err)
	}

	return nil
}
