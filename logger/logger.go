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

// Package logger implements logging package. It defines a type, Logger, with
// methods for formatting output. Each logging operations creates and sends
// lightweight not formatted log message to separate worker thread. It offloads
// main code from unnecessary resource consuming formatting and I/O operations.
// On default it supports many different log handlers like logging to standard
// output, error output, file, stream or syslog.
package logger

import (
	"os"
	"sync"
	"time"
)

// These constants define log level values and names used by various logger
// functions like for example Debug or Info. It defines also default logger
// values.
const (
	OffsetLevel = 10

	TraceLevel    = 0
	DebugLevel    = OffsetLevel + TraceLevel
	InfoLevel     = OffsetLevel + DebugLevel
	NoticeLevel   = OffsetLevel + InfoLevel
	WarningLevel  = OffsetLevel + NoticeLevel
	ErrorLevel    = OffsetLevel + WarningLevel
	CriticalLevel = OffsetLevel + ErrorLevel
	AlertLevel    = OffsetLevel + CriticalLevel
	FatalLevel    = OffsetLevel + AlertLevel
	PanicLevel    = OffsetLevel + FatalLevel

	MinimumLevel = TraceLevel
	MaximumLevel = PanicLevel

	TraceName    = "trace"
	DebugName    = "debug"
	InfoName     = "info"
	NoticeName   = "notice"
	WarningName  = "warning"
	ErrorName    = "error"
	CriticalName = "critical"
	AlertName    = "alert"
	FatalName    = "fatal"
	PanicName    = "panic"

	DefaultTypeName = "log"

	DefaultErrorCode = 1

	loggerSkipCall = 3
)

// A Logger represents an active logging object that generates log messages for
// different added log handlers. Each logging operations creates and sends
// lightweight not formatted log message to separate worker thread. It offloads
// main code from unnecessary resource consuming formatting and I/O operations.
type Logger struct {
	name        string
	handlers    Handlers
	idGenerator IDGenerator
	errorCode   int
	mutex       sync.RWMutex
}

// New creates new logger instance with default handlers.
func New() *Logger {
	return &Logger{
		handlers: Handlers{
			"stdout": NewStdout(),
			"stderr": NewStderr(),
		},
		errorCode:   DefaultErrorCode,
		idGenerator: NewUUID4(),
	}
}

// SetErrorCode sets error code that is returned during Fatal call.
// On default it is 1.
func (l *Logger) SetErrorCode(errorCode int) *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.errorCode = errorCode

	return l
}

// GetErrorCode returns error code.
func (l *Logger) GetErrorCode() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.errorCode
}

// SetName sets logger name.
func (l *Logger) SetName(name string) *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.name = name

	return l
}

// GetName returns logger name.
func (l *Logger) GetName() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.name
}

// AddHandler sets log handler under provided identifier name.
func (l *Logger) AddHandler(name string, handler Handler) *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.handlers[name] = handler

	return l
}

// SetHandlers sets log handlers for logger.
func (l *Logger) SetHandlers(handlers Handlers) *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.handlers = handlers

	return l
}

// GetHandler returns added log handler by provided name.
func (l *Logger) GetHandler(name string) (Handler, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	handler, ok := l.handlers[name]

	if !ok {
		return nil, NewRuntimeError("cannot get handler"+name, nil)
	}

	return handler, nil
}

// GetHandlers returns all added log handlers.
func (l *Logger) GetHandlers() Handlers {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.handlers
}

// RemoveHandler removes added log handler by provided name.
func (l *Logger) RemoveHandler(name string) *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	delete(l.handlers, name)

	return l
}

// RemoveHandlers removes all added log handlers.
func (l *Logger) RemoveHandlers() *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.handlers = make(Handlers)

	return l
}

// ResetHandlers sets logger default log handlers.
func (l *Logger) ResetHandlers() *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.handlers = Handlers{
		"stdout": NewStdout(),
		"stderr": NewStderr(),
	}

	return l
}

// Reset resets logger to default state and default log handlers.
func (l *Logger) Reset() *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.idGenerator = NewUUID4()
	l.errorCode = DefaultErrorCode
	l.handlers = Handlers{
		"stdout": NewStdout(),
		"stderr": NewStderr(),
	}

	return l
}

// SetIDGenerator sets ID generator function that is called by logger to
// generate ID for created log messages.
func (l *Logger) SetIDGenerator(idGenerator IDGenerator) *Logger {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.idGenerator = idGenerator

	return l
}

// GetIDGenerator returns ID generator function that is called by logger to
// generate ID for created log messages.
func (l *Logger) GetIDGenerator() IDGenerator {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.idGenerator
}

// Trace logs finer-grained informational messages than the Debug. It creates
// and sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers.
func (l *Logger) Trace(message string, arguments ...interface{}) {
	l.logMessage(TraceLevel, TraceName, message, arguments...)
}

// Debug logs debugging messages. It creates and sends lightweight not formatted
// log messages to separate running logger thread for further formatting and
// I/O handling from different added log handlers.
func (l *Logger) Debug(message string, arguments ...interface{}) {
	l.logMessage(DebugLevel, DebugName, message, arguments...)
}

// Info logs informational messages. It creates and sends lightweight not
// formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers.
func (l *Logger) Info(message string, arguments ...interface{}) {
	l.logMessage(InfoLevel, InfoName, message, arguments...)
}

// Notice logs messages for significant conditions. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers.
func (l *Logger) Notice(message string, arguments ...interface{}) {
	l.logMessage(NoticeLevel, NoticeName, message, arguments...)
}

// Warning logs messages for warning conditions that can be potentially harmful.
// It creates and sends lightweight not formatted log messages to separate
// running logger thread for further formatting and I/O handling from different
// added log handlers.
func (l *Logger) Warning(message string, arguments ...interface{}) {
	l.logMessage(WarningLevel, WarningName, message, arguments...)
}

// Error logs messages for error conditions. It creates and sends lightweight
// not formatted log messages to separate running logger thread for further
// formatting and I/O handling from different log handlers.
func (l *Logger) Error(message string, arguments ...interface{}) {
	l.logMessage(ErrorLevel, ErrorName, message, arguments...)
}

// Critical logs messages for critical conditions. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers.
func (l *Logger) Critical(message string, arguments ...interface{}) {
	l.logMessage(CriticalLevel, CriticalName, message, arguments...)
}

// Alert logs messages for alert conditions. It creates and sends lightweight
// not formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers.
func (l *Logger) Alert(message string, arguments ...interface{}) {
	l.logMessage(AlertLevel, AlertName, message, arguments...)
}

// Fatal logs messages for fatal conditions. It stops logger worker thread and
// it exists the application with an error code. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers.
func (l *Logger) Fatal(message string, arguments ...interface{}) {
	l.logMessage(FatalLevel, FatalName, message, arguments...)
	Close()
	os.Exit(l.errorCode) // revive:disable-line
}

// Panic logs messages for fatal conditions. It stops logger worker thread and
// it exists the application with a panic. It creates and sends lightweight not
// formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers.
func (l *Logger) Panic(message string, arguments ...interface{}) {
	l.logMessage(PanicLevel, PanicName, message, arguments...)
	Close()
	panic(NewRuntimeError("Panic error", nil))
}

// Log logs messages with user defined log level value and name. It creates and
// sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers.
func (l *Logger) Log(level int, levelName, message string, arguments ...interface{}) {
	l.logMessage(level, levelName, message, arguments...)
}

// Flush flushes all log messages.
func (l *Logger) Flush() *Logger {
	GetWorker().Flush()

	return l
}

// Close closes all added log handlers.
func (l *Logger) Close() error {
	GetWorker().Flush()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	var err error

	for _, handler := range l.handlers {
		handlerError := handler.Close()

		if handlerError != nil {
			err = NewRuntimeError("cannot close log handler", handlerError)
			printError(err)
		}
	}

	return err
}

// logMessage logs message with defined log level value and name. It creates and
// sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers.
func (l *Logger) logMessage(level int, levelName, message string, arguments ...interface{}) {
	now := time.Now()
	path, line, function := getPathLineFunction(loggerSkipCall)

	GetWorker().records <- &Record{
		Time:      now,
		Message:   message,
		Arguments: arguments,
		Level: Level{
			Name:  levelName,
			Value: level,
		},
		File: Source{
			Line:     line,
			Path:     path,
			Function: function,
		},
		logger: l,
	}
}

// Emit emits provided log record to logger worker thread for further
// formatting and I/O handling from different addded log handlers.
func (l *Logger) Emit(record *Record) *Logger {
	record.logger = l
	GetWorker().records <- record

	return l
}
