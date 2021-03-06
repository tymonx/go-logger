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
	"os"
	"sync"
)

var gOnce sync.Once   // nolint:gochecknoglobals
var gInstance *Logger // nolint:gochecknoglobals

// Get returns global logger instance.
func Get() *Logger {
	gOnce.Do(func() {
		gInstance = New()
	})

	return gInstance
}

// Enable enables all added log handlers.
func Enable() *Logger {
	return Get().Enable()
}

// Disable disabled all added log handlers.
func Disable() *Logger {
	return Get().Disable()
}

// IsEnabled returns true if at least one of added log handlers is enabled.
func IsEnabled() bool {
	return Get().IsEnabled()
}

// SetLevel sets log level to all added log handlers.
func SetLevel(level int) *Logger {
	return Get().SetLevel(level)
}

// SetMinimumLevel sets minimum log level to all added log handlers.
func SetMinimumLevel(level int) *Logger {
	return Get().SetMinimumLevel(level)
}

// SetMaximumLevel sets maximum log level to all added log handlers.
func SetMaximumLevel(level int) *Logger {
	return Get().SetMaximumLevel(level)
}

// SetLevelRange sets minimum and maximum log level values to all added log handlers.
func SetLevelRange(min, max int) *Logger {
	return Get().SetLevelRange(min, max)
}

// SetFormatter sets provided formatter to all added log handlers.
func SetFormatter(formatter *Formatter) *Logger {
	return Get().SetFormatter(formatter)
}

// SetFormat sets provided format string to all added log handlers.
func SetFormat(format string) *Logger {
	return Get().SetFormat(format)
}

// SetDateFormat sets provided date format string to all added log handlers.
func SetDateFormat(format string) *Logger {
	return Get().SetDateFormat(format)
}

// SetPlaceholder sets provided placeholder string to all added log handlers.
func SetPlaceholder(placeholder string) *Logger {
	return Get().SetPlaceholder(placeholder)
}

// AddFuncs adds template functions to format log message to all added log handlers.
func AddFuncs(funcs FormatterFuncs) *Logger {
	return Get().AddFuncs(funcs)
}

// ResetFormatters resets all formatters from added log handlers.
func ResetFormatters() *Logger {
	return Get().ResetFormatters()
}

// SetErrorCode sets error code that is returned during Fatal call.
// On default it is 1.
func SetErrorCode(errorCode int) *Logger {
	return Get().SetErrorCode(errorCode)
}

// GetErrorCode returns error code.
func GetErrorCode() int {
	return Get().GetErrorCode()
}

// SetName sets logger name.
func SetName(name string) *Logger {
	return Get().SetName(name)
}

// GetName returns logger name.
func GetName() string {
	return Get().GetName()
}

// AddHandler sets log handler under provided identifier name.
func AddHandler(name string, handler Handler) *Logger {
	return Get().AddHandler(name, handler)
}

// SetHandler sets a single log handler for logger. It is equivalent to
// logger.RemoveHandlers().SetHandlers(logger.Handlers{name: handler}).
func SetHandler(name string, handler Handler) *Logger {
	return Get().SetHandler(name, handler)
}

// SetHandlers sets log handlers for logger.
func SetHandlers(handlers Handlers) *Logger {
	return Get().SetHandlers(handlers)
}

// GetHandler returns added log handler by provided name.
func GetHandler(name string) (Handler, error) {
	return Get().GetHandler(name)
}

// GetHandlers returns all added log handlers.
func GetHandlers() Handlers {
	return Get().GetHandlers()
}

// RemoveHandler removes added log handler by provided name.
func RemoveHandler(name string) *Logger {
	return Get().RemoveHandler(name)
}

// RemoveHandlers removes all added log handlers.
func RemoveHandlers() *Logger {
	return Get().RemoveHandlers()
}

// ResetHandlers sets logger default log handlers.
func ResetHandlers() *Logger {
	return Get().ResetHandlers()
}

// Reset resets logger to default state and default log handlers.
func Reset() *Logger {
	return Get().ResetHandlers()
}

// SetIDGenerator sets ID generator function that is called by logger to
// generate ID for created log messages.
func SetIDGenerator(idGenerator IDGenerator) *Logger {
	return Get().SetIDGenerator(idGenerator)
}

// GetIDGenerator returns ID generator function that is called by logger to
// generate ID for created log messages.
func GetIDGenerator() IDGenerator {
	return Get().GetIDGenerator()
}

// Trace logs finer-grained informational messages than the Debug. It creates
// and sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers.
func Trace(message string, arguments ...interface{}) {
	Get().LogMessage(TraceLevel, TraceName, message, arguments...)
}

// Debug logs debugging messages. It creates and sends lightweight not formatted
// log messages to separate running logger thread for further formatting and
// I/O handling from different added log handlers.
func Debug(message string, arguments ...interface{}) {
	Get().LogMessage(DebugLevel, DebugName, message, arguments...)
}

// Info logs informational messages. It creates and sends lightweight not
// formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers.
func Info(message string, arguments ...interface{}) {
	Get().LogMessage(InfoLevel, InfoName, message, arguments...)
}

// Notice logs messages for significant conditions. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers.
func Notice(message string, arguments ...interface{}) {
	Get().LogMessage(NoticeLevel, NoticeName, message, arguments...)
}

// Warning logs messages for warning conditions that can be potentially harmful.
// It creates and sends lightweight not formatted log messages to separate
// running logger thread for further formatting and I/O handling from different
// added log handlers.
func Warning(message string, arguments ...interface{}) {
	Get().LogMessage(WarningLevel, WarningName, message, arguments...)
}

// Error logs messages for error conditions. It creates and sends lightweight
// not formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers.
func Error(message string, arguments ...interface{}) {
	Get().LogMessage(ErrorLevel, ErrorName, message, arguments...)
}

// Critical logs messages for critical conditions. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers.
func Critical(message string, arguments ...interface{}) {
	Get().LogMessage(CriticalLevel, CriticalName, message, arguments...)
}

// Alert logs messages for alert conditions. It creates and sends lightweight
// not formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers.
func Alert(message string, arguments ...interface{}) {
	Get().LogMessage(AlertLevel, AlertName, message, arguments...)
}

// Fatal logs messages for fatal conditions. It stops logger worker thread and
// it exists the application with an error code. It creates and sends
// lightweight not formatted log messages to separate running logger thread for
// further formatting and I/O handling from different added log handlers.
func Fatal(message string, arguments ...interface{}) {
	Get().LogMessage(FatalLevel, FatalName, message, arguments...)
	Close()
	os.Exit(Get().GetErrorCode()) // revive:disable-line
}

// Panic logs messages for fatal conditions. It stops logger worker thread and
// it exists the application with a panic. It creates and sends lightweight not
// formatted log messages to separate running logger thread for further
// formatting and I/O handling from different added log handlers.
func Panic(message string, arguments ...interface{}) {
	Get().LogMessage(PanicLevel, PanicName, message, arguments...)
	Close()
	panic(NewRuntimeError("Panic error"))
}

// Log logs messages with user defined log level value and name. It creates and
// sends lightweight not formatted log messages to separate running logger
// thread for further formatting and I/O handling from different added log
// handlers.
func Log(level int, levelName, message string, arguments ...interface{}) {
	Get().LogMessage(level, levelName, message, arguments...)
}

// Emit emits provided log record to logger worker thread for further
// formatting and I/O handling from different addded log handlers.
func Emit(record *Record) *Logger {
	return Get().Emit(record)
}

// Flush flushes all log messages.
func Flush() *Logger {
	return Get().Flush()
}

// Close closes all added log handlers.
func Close() {
	err := Get().Close()

	if err != nil {
		printError(NewRuntimeError("cannot close logger"))
	}
}
