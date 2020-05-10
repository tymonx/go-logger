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
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// These constants define log level values and names used by various logger
// functions like for example Debug or Info. It defines also default logger
// values
const (
	TraceLevel    = 0
	DebugLevel    = 10
	InfoLevel     = 20
	NoticeLevel   = 30
	WarningLevel  = 40
	ErrorLevel    = 50
	CriticalLevel = 60
	AlertLevel    = 70
	FatalLevel    = 80
	PanicLevel    = 90

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
)

// A Logger represents an active logging object that generates log messages for
// different added log handlers. Each logging operations creates and sends
// lightweight not formatted log message to separate worker thread. It offloads
// main code from unnecessary resource consuming formatting and I/O operations
type Logger struct {
	tag         string
	name        string
	handlers    Handlers
	idGenerator IDGenerator
	errorCode   int
}

// New creates new logger instance with default handlers
func New() *Logger {
	return &Logger{
		handlers: Handlers{
			"stdout": NewStdout(),
			"stderr": NewStderr(),
		},
		errorCode:   DefaultErrorCode,
		idGenerator: uuid4,
	}
}

// SetErrorCode sets error code that is returned during Fatal call.
// On default it is 1
func (logger *Logger) SetErrorCode(errorCode int) *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.errorCode = errorCode
	return logger
}

// GetErrorCode returns error code
func (logger *Logger) GetErrorCode() int {
	gMutex.RLock()
	defer gMutex.RUnlock()

	return logger.errorCode
}

// SetName sets logger name
func (logger *Logger) SetName(name string) *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.name = name
	return logger
}

// GetName returns logger name
func (logger *Logger) GetName() string {
	gMutex.RLock()
	defer gMutex.RUnlock()

	return logger.name
}

// AddHandler sets log handler under provided identifier name
func (logger *Logger) AddHandler(name string, handler Handler) *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.handlers[name] = handler
	return logger
}

// CreateAddHandler it creates registered log handler by provided name and it
// sets for logger
func (logger *Logger) CreateAddHandler(name string) *Logger {
	return logger.AddHandler(name, CreateHandler(name))
}

// SetHandlers sets log handlers for logger
func (logger *Logger) SetHandlers(handlers Handlers) *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.handlers = handlers
	return logger
}

// GetHandler returns added log handler by provided name
func (logger *Logger) GetHandler(name string) (handler Handler, ok bool) {
	gMutex.RLock()
	defer gMutex.RUnlock()

	handler, ok = logger.handlers[name]
	return
}

// GetHandlers returns all added log handlers
func (logger *Logger) GetHandlers() Handlers {
	gMutex.RLock()
	defer gMutex.RUnlock()

	return logger.handlers
}

// RemoveHandler removes added log handler by provided name
func (logger *Logger) RemoveHandler(name string) *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	delete(logger.handlers, name)
	return logger
}

// RemoveHandlers removes all added log handlers
func (logger *Logger) RemoveHandlers() *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.handlers = make(Handlers)
	return logger
}

// ResetHandlers sets logger default log handlers
func (logger *Logger) ResetHandlers() *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.handlers = Handlers{
		"stdout": NewStdout(),
		"stderr": NewStderr(),
	}

	return logger
}

// Reset resets logger to default state and default log handlers
func (logger *Logger) Reset() *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.idGenerator = uuid4
	logger.errorCode = DefaultErrorCode
	logger.handlers = Handlers{
		"stdout": NewStdout(),
		"stderr": NewStderr(),
	}

	return logger
}

// SetIDGenerator sets ID generator function that is called by logger to
// generate ID for created log messages
func (logger *Logger) SetIDGenerator(idGenerator IDGenerator) *Logger {
	gMutex.Lock()
	defer gMutex.Unlock()

	logger.idGenerator = idGenerator
	return logger
}

// GetIDGenerator returns ID generator function that is called by logger to
// generate ID for created log messages
func (logger *Logger) GetIDGenerator() IDGenerator {
	gMutex.RLock()
	defer gMutex.RUnlock()

	return logger.idGenerator
}

// CreateSetIDGenerator it creates registered ID generator function by provided
// name and it sets for logger
func (logger *Logger) CreateSetIDGenerator(name string) *Logger {
	logger.SetIDGenerator(CreateIDGenerator(name))
	return logger
}

// Trace logs finer-grained informational messages than the Debug. It creates
// and sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers
func (logger *Logger) Trace(message string, arguments ...interface{}) {
	logger.logMessage(TraceLevel, TraceName, message, arguments...)
}

// Debug logs debugging messages. It creates and sends lightweight not formatted
// log messages to separate running logger thread for further formatting and
// I/O handling from different added log handlers
func (logger *Logger) Debug(message string, arguments ...interface{}) {
	logger.logMessage(DebugLevel, DebugName, message, arguments...)
}

// Info logs informational messages. It creates and sends lightweight not
// formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers
func (logger *Logger) Info(message string, arguments ...interface{}) {
	logger.logMessage(InfoLevel, InfoName, message, arguments...)
}

// Notice logs messages for significant conditions. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers
func (logger *Logger) Notice(message string, arguments ...interface{}) {
	logger.logMessage(NoticeLevel, NoticeName, message, arguments...)
}

// Warning logs messages for warning conditions that can be potentially harmful.
// It creates and sends lightweight not formatted log messages to separate
// running logger thread for further formatting and I/O handling from different
// added log handlers
func (logger *Logger) Warning(message string, arguments ...interface{}) {
	logger.logMessage(WarningLevel, WarningName, message, arguments...)
}

// Error logs messages for error conditions. It creates and sends lightweight
// not formatted log messages to separate running logger thread for further
// formatting and I/O handling from different log handlers
func (logger *Logger) Error(message string, arguments ...interface{}) {
	logger.logMessage(ErrorLevel, ErrorName, message, arguments...)
}

// Critical logs messages for critical conditions. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers
func (logger *Logger) Critical(message string, arguments ...interface{}) {
	logger.logMessage(CriticalLevel, CriticalName, message, arguments...)
}

// Alert logs messages for alert conditions. It creates and sends lightweight
// not formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers
func (logger *Logger) Alert(message string, arguments ...interface{}) {
	logger.logMessage(AlertLevel, AlertName, message, arguments...)
}

// Fatal logs messages for fatal conditions. It stops logger worker thread and
// it exists the application with an error code. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers
func (logger *Logger) Fatal(message string, arguments ...interface{}) {
	logger.logMessage(FatalLevel, FatalName, message, arguments...)
	Close()
	os.Exit(logger.errorCode)
}

// Panic logs messages for fatal conditions. It stops logger worker thread and
// it exists the application with a panic. It creates and sends lightweight not
// formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers
func (logger *Logger) Panic(message string, arguments ...interface{}) {
	logger.logMessage(PanicLevel, PanicName, message, arguments...)
	Close()
	panic(message)
}

// Log logs messages with user defined log level value and name. It creates and
// sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers
func (logger *Logger) Log(level int, levelName string, message string, arguments ...interface{}) {
	logger.logMessage(level, levelName, message, arguments...)
}

// logMessage logs message with defined log level value and name. It creates and
// sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers
func (logger *Logger) logMessage(level int, levelName string, message string, arguments ...interface{}) {
	now := time.Now()
	path, line, function := getPathLineFunction(3)

	getWorker().records <- &Record{
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
		logger: logger,
	}
}

// emit prepares provided log record and it dispatches to all added log
// handlers for further formatting and specific I/O implementation operations
func (logger *Logger) emit(record *Record) {
	record.Type = DefaultTypeName
	record.File.Name = filepath.Base(record.File.Path)
	record.File.Function = filepath.Base(record.File.Function)
	record.Timestamp.Created = record.Time.Format(time.RFC3339)
	record.Address = getAddress()
	record.Hostname = getHostname()

	gMutex.RLock()
	defer gMutex.RUnlock()

	record.ID = logger.idGenerator()
	record.Name = logger.name

	if record.Name == "" {
		record.Name = filepath.Base(os.Args[0])
	}

	for _, handler := range logger.handlers {
		min, max := handler.GetLevelRange()

		if (record.Level.Value >= min) && (record.Level.Value <= max) {
			err := handler.Emit(record)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// Close closes all added log handlers
func (logger *Logger) Close() {
	gMutex.Lock()
	defer gMutex.Unlock()

	for _, handler := range logger.handlers {
		err := handler.Close()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
