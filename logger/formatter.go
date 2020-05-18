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
	record        *Record
	template      *template.Template
	placeholder   string
	timeBuffer    *bytes.Buffer
	formatBuffer  *bytes.Buffer
	messageBuffer *bytes.Buffer
	mutex         sync.RWMutex
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

	f.setFuncs()

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

// GetRecord returns assigned log record object to formatter.
func (f *Formatter) GetRecord() *Record {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.record
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

	message, err := f.formatRecord(record)

	if err != nil {
		return "", NewRuntimeError("cannot format record", err)
	}

	return message, nil
}

// FormatTime returns formatted date string based on provided log record object.
func (f *Formatter) FormatTime(record *Record) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	message, err := f.formatTimeUnsafe(record)

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

	message, err := f.formatMessageUnsafe(record)

	if err != nil {
		return "", NewRuntimeError("cannot format message", err)
	}

	return message, nil
}

// Format returns formatted log message string based on provided log record
// object.
func (f *Formatter) formatRecord(record *Record) (string, error) {
	f.record = record

	return f.formatString(
		f.template,
		f.formatBuffer,
		f.format,
		nil,
	)
}

// formatTimeUnsafe returns formatted date string based on provided log record object.
func (f *Formatter) formatTimeUnsafe(record *Record) (string, error) {
	f.record = record

	return f.formatString(
		f.template,
		f.timeBuffer,
		f.dateFormat,
		nil,
	)
}

// formatMessageUnsafe returns formatted user message string based on provided log
// record object.
func (f *Formatter) formatMessageUnsafe(record *Record) (string, error) {
	var err error

	f.record = record
	message := record.Message

	if len(record.Arguments) > 0 {
		var object interface{}

		funcMap := make(template.FuncMap)

		funcMap[f.placeholder] = f.argumentAutomatic()

		for position, argument := range record.Arguments {
			placeholder := f.placeholder + strconv.Itoa(position)

			funcMap[placeholder] = f.argumentValue(argument)

			valueOf := reflect.ValueOf(argument)

			switch valueOf.Kind() {
			case reflect.Map:
				if reflect.TypeOf(argument).Key().Kind() == reflect.String {
					for _, key := range valueOf.MapKeys() {
						funcMap[key.String()] = f.argumentValue(valueOf.MapIndex(key).Interface())
					}
				}
			case reflect.Struct:
				object = argument
			}
		}

		message, err = f.formatString(
			template.New("").Delims("{", "}").Funcs(funcMap),
			f.messageBuffer,
			message,
			object,
		)
	}

	return message, err
}

// argumentValue returns closure that returns log argument used in log message.
func (*Formatter) argumentValue(argument interface{}) func() interface{} {
	return func() interface{} {
		return argument
	}
}

// argumentAutomatic returns closure that returns log argument from automatic
// placeholder used in log message.
func (f *Formatter) argumentAutomatic() func() interface{} {
	position := 0
	arguments := len(f.record.Arguments)

	return func() interface{} {
		var argument interface{}

		if position < arguments {
			argument = f.record.Arguments[position]
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

// setFuncs sets default template functions used to formatting log message.
func (f *Formatter) setFuncs() {
	f.template.Funcs(template.FuncMap{
		"gid":       os.Getgid,
		"pid":       os.Getpid,
		"ppid":      os.Getppid,
		"getEnv":    os.Getenv,
		"expandEnv": os.ExpandEnv,
		"executable": func() string {
			return filepath.Base(os.Args[0])
		},
		"date": func() string {
			date, err := f.formatTimeUnsafe(f.record)

			if err != nil {
				printError(NewRuntimeError("cannot format date", err))
			}

			return date
		},
		"message": func() string {
			message, err := f.formatMessageUnsafe(f.record)

			if err != nil {
				printError(NewRuntimeError("cannot format message", err))
			}

			return message
		},
		"levelValue": func() int {
			return f.record.Level.Value
		},
		"level": func() string {
			return strings.ToLower(f.record.Level.Name)
		},
		"Level": func() string {
			return strings.Title(strings.ToLower(f.record.Level.Name))
		},
		"LEVEL": func() string {
			return strings.ToUpper(f.record.Level.Name)
		},
		"iso8601": func() string {
			return f.record.Time.Format(time.RFC3339)
		},
		"id": func() interface{} {
			return f.record.ID
		},
		"name": func() string {
			return f.record.Name
		},
		"host": func() string {
			return f.record.Address
		},
		"hostname": func() string {
			return f.record.Hostname
		},
		"address": func() string {
			return f.record.Address
		},
		"nanosecond": func() string {
			return fmt.Sprintf("%09d", f.record.Time.Nanosecond())
		},
		"microsecond": func() string {
			return fmt.Sprintf("%06d", f.record.Time.Nanosecond()/kilo)
		},
		"millisecond": func() string {
			return fmt.Sprintf("%03d", f.record.Time.Nanosecond()/mega)
		},
		"second": func() string {
			return fmt.Sprintf("%02d", f.record.Time.Second())
		},
		"minute": func() string {
			return fmt.Sprintf("%02d", f.record.Time.Minute())
		},
		"hour": func() string {
			return fmt.Sprintf("%02d", f.record.Time.Hour())
		},
		"day": func() string {
			return fmt.Sprintf("%02d", f.record.Time.Day())
		},
		"month": func() string {
			return fmt.Sprintf("%02d", f.record.Time.Month())
		},
		"YEAR": func() string {
			return fmt.Sprintf("%02d", f.record.Time.Year()%percentage)
		},
		"year": func() int {
			return f.record.Time.Year()
		},
		"file": func() string {
			return f.record.File.Name
		},
		"line": func() int {
			return f.record.File.Line
		},
		"function": func() string {
			return f.record.File.Function
		},
	})
}
