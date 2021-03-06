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
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

// These constants define default values for Formatter.
const (
	DefaultDateFormat  = "{year}-{month}-{day} {hour}:{minute}:{second},{millisecond}"
	DefaultFormat      = "{date} - {Level | printf \"%-8s\"} - {file}:{line}:{function}(): {message}"
	DefaultPlaceholder = "p"

	kilo       = 1e3
	mega       = 1e6
	percentage = 100
)

// FormatterFuncs defines map of template functions.
type FormatterFuncs map[string]interface{}

// A Formatter represents a formatter object used by log handler to format log
// message.
type Formatter struct {
	format        string
	dateFormat    string
	template      *template.Template
	placeholder   string
	timeBuffer    *bytes.Buffer
	formatBuffer  *bytes.Buffer
	messageBuffer *bytes.Buffer
	mutex         sync.RWMutex
	usedArguments map[int]bool
}

// NewFormatter creates a new Formatter object with default format settings.
func NewFormatter() *Formatter {
	f := &Formatter{
		format:        DefaultFormat,
		dateFormat:    DefaultDateFormat,
		template:      template.New("").Delims("{", "}"),
		placeholder:   DefaultPlaceholder,
		timeBuffer:    new(bytes.Buffer),
		formatBuffer:  new(bytes.Buffer),
		messageBuffer: new(bytes.Buffer),
	}

	return f
}

// Reset resets Formatter.
func (f *Formatter) Reset() *Formatter {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.format = DefaultFormat
	f.dateFormat = DefaultDateFormat
	f.placeholder = DefaultPlaceholder

	return f
}

// SetPlaceholder sets placeholder string prefix used for automatic and
// positional placeholders to format log message.
func (f *Formatter) SetPlaceholder(placeholder string) *Formatter {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.placeholder = placeholder

	return f
}

// GetPlaceholder returns placeholder string prefix used for automatic and
// positional placeholders to format log message.
func (f *Formatter) GetPlaceholder() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.placeholder
}

// AddFuncs adds template functions to format log message.
func (f *Formatter) AddFuncs(funcs FormatterFuncs) *Formatter {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.template.Funcs(template.FuncMap(funcs))

	return f
}

// SetFormat sets format string used for formatting log message.
func (f *Formatter) SetFormat(format string) *Formatter {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.format = format

	return f
}

// GetFormat returns format string used for formatting log message.
func (f *Formatter) GetFormat() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.format
}

// SetDateFormat sets format string used for formatting date in log message.
func (f *Formatter) SetDateFormat(dateFormat string) *Formatter {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.dateFormat = dateFormat

	return f
}

// GetDateFormat returns format string used for formatting date in log message.
func (f *Formatter) GetDateFormat() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.format
}

// Format returns formatted log message string based on provided log record
// object.
func (f *Formatter) Format(record *Record) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.template.Funcs(f.getRecordFuncs(record))

	message, err := f.formatString(f.template, f.formatBuffer, f.format, nil)

	if err != nil {
		return "", NewRuntimeError("cannot format record", err)
	}

	return message, nil
}

// FormatTime returns formatted date string based on provided log record object.
func (f *Formatter) FormatTime(record *Record) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.template.Funcs(f.getRecordFuncs(record))

	message, err := f.formatString(f.template, f.timeBuffer, f.dateFormat, nil)

	if err != nil {
		return "", NewRuntimeError("cannot format time", err)
	}

	return message, nil
}

// FormatMessage returns formatted user message string based on provided log
// record object.
func (f *Formatter) FormatMessage(record *Record) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	message, err := f.formatMessageRecord(record)

	if err != nil {
		return "", NewRuntimeError("cannot format message", err)
	}

	return message, nil
}

// formatMessageUnsafe returns formatted user message string based on provided log
// record object.
func (f *Formatter) formatMessageRecord(record *Record) (string, error) {
	if len(record.Arguments) == 0 {
		return record.Message, nil
	}

	var err error

	var object interface{}

	message := record.Message

	f.usedArguments = make(map[int]bool)

	funcMap := make(template.FuncMap)

	funcMap[f.placeholder] = f.argumentAutomatic(record)

	for position, argument := range record.Arguments {
		placeholder := f.placeholder + strconv.Itoa(position)

		funcMap[placeholder] = f.argumentValue(position, argument)

		valueOf := reflect.ValueOf(argument)

		switch valueOf.Kind() {
		case reflect.Map:
			if reflect.TypeOf(argument).Key().Kind() == reflect.String {
				for _, key := range valueOf.MapKeys() {
					funcMap[key.String()] = f.argumentValue(position, valueOf.MapIndex(key).Interface())
				}
			}
		case reflect.Struct:
			object = argument
		}
	}

	if message, err = f.formatString(
		template.New("").Delims("{", "}").Funcs(f.getRecordFuncs(record)).Funcs(funcMap),
		f.messageBuffer,
		message,
		object,
	); err != nil {
		return "", err
	}

	if len(f.usedArguments) >= len(record.Arguments) {
		return message, nil
	}

	for position, argument := range record.Arguments {
		if !f.isArgumentUsed(position, argument) {
			if message != "" {
				message += " "
			}

			message += fmt.Sprint(argument)
		}
	}

	return message, nil
}

func (f *Formatter) isArgumentUsed(position int, argument interface{}) bool {
	valueOf := reflect.ValueOf(argument)

	switch valueOf.Kind() {
	case reflect.Map:
		if reflect.TypeOf(argument).Key().Kind() == reflect.String {
			return true
		}
	case reflect.Struct:
		return true
	}

	return f.usedArguments[position]
}

// argumentValue returns closure that returns log argument used in log message.
func (f *Formatter) argumentValue(position int, argument interface{}) func() interface{} {
	return func() interface{} {
		f.usedArguments[position] = true
		return argument
	}
}

// argumentAutomatic returns closure that returns log argument from automatic
// placeholder used in log message.
func (f *Formatter) argumentAutomatic(record *Record) func() interface{} {
	position := 0
	arguments := len(record.Arguments)

	return func() interface{} {
		var argument interface{}

		if position < arguments {
			f.usedArguments[position] = true
			argument = record.Arguments[position]
			position++
		}

		return argument
	}
}

// formatString returns formatted string.
func (*Formatter) formatString(templ *template.Template, buffer *bytes.Buffer, format string, object interface{}) (string, error) {
	var message string

	if format != "" {
		var err error

		templ, err = templ.Parse(format)

		if err != nil {
			return "", NewRuntimeError("cannot parse text template", err)
		}

		buffer.Reset()

		err = templ.Execute(buffer, object)

		if err != nil {
			return "", NewRuntimeError("cannot execute text template", err)
		}

		message = buffer.String()
	}

	return message, nil
}

// getRecordFuncs sets default template functions used to formatting log message.
func (f *Formatter) getRecordFuncs(record *Record) template.FuncMap {
	return template.FuncMap{
		"uid":       os.Getuid,
		"gid":       os.Getgid,
		"pid":       os.Getpid,
		"egid":      os.Getegid,
		"euid":      os.Geteuid,
		"ppid":      os.Getppid,
		"getEnv":    os.Getenv,
		"expandEnv": os.ExpandEnv,
		"executable": func() string {
			return filepath.Base(os.Args[0])
		},
		"date": func() string {
			date, err := f.formatString(f.template, f.timeBuffer, f.dateFormat, nil)

			if err != nil {
				printError(NewRuntimeError("cannot format date", err))
			}

			return date
		},
		"message": func() string {
			message, err := f.formatMessageRecord(record)

			if err != nil {
				printError(NewRuntimeError("cannot format message", err))
			}

			return message
		},
		"levelValue": func() int {
			return record.Level.Value
		},
		"level": func() string {
			return strings.ToLower(record.Level.Name)
		},
		"Level": func() string {
			return strings.Title(strings.ToLower(record.Level.Name))
		},
		"LEVEL": func() string {
			return strings.ToUpper(record.Level.Name)
		},
		"iso8601": func() string {
			return record.Time.Format(time.RFC3339)
		},
		"id": func() interface{} {
			return record.ID
		},
		"name": func() string {
			return record.Name
		},
		"host": func() string {
			return record.Address
		},
		"hostname": func() string {
			return record.Hostname
		},
		"address": func() string {
			return record.Address
		},
		"nanosecond": func() string {
			return fmt.Sprintf("%09d", record.Time.Nanosecond())
		},
		"microsecond": func() string {
			return fmt.Sprintf("%06d", record.Time.Nanosecond()/kilo)
		},
		"millisecond": func() string {
			return fmt.Sprintf("%03d", record.Time.Nanosecond()/mega)
		},
		"second": func() string {
			return fmt.Sprintf("%02d", record.Time.Second())
		},
		"minute": func() string {
			return fmt.Sprintf("%02d", record.Time.Minute())
		},
		"hour": func() string {
			return fmt.Sprintf("%02d", record.Time.Hour())
		},
		"day": func() string {
			return fmt.Sprintf("%02d", record.Time.Day())
		},
		"month": func() string {
			return fmt.Sprintf("%02d", record.Time.Month())
		},
		"YEAR": func() string {
			return fmt.Sprintf("%02d", record.Time.Year()%percentage)
		},
		"year": func() int {
			return record.Time.Year()
		},
		"file": func() string {
			return record.File.Name
		},
		"line": func() int {
			return record.File.Line
		},
		"function": func() string {
			return record.File.Function
		},
	}
}
