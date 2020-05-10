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
	"text/template"
	"time"
)

// These constants define default values for Formatter
const (
	DefaultDateFormat = "{year}-{month}-{day} {hour}:{minute}:{second},{milisecond}"
	DefaultFormat     = "{date} - {Level | printf \"%-8s\"} - {file}:{line}:{function}(): {message}"
)

// FormatterFuncs defines map of template functions
type FormatterFuncs map[string]interface{}

// A Formatter represents a formatter object used by log handler to format log
// message
type Formatter struct {
	format        string
	dateFormat    string
	record        *Record
	template      *template.Template
	placeholder   string
	timeBuffer    *bytes.Buffer
	formatBuffer  *bytes.Buffer
	messageBuffer *bytes.Buffer
}

// NewFormatter creates a new Formatter object with default format settings
func NewFormatter() *Formatter {
	formatter := &Formatter{
		format:        DefaultFormat,
		dateFormat:    DefaultDateFormat,
		template:      template.New("").Delims("{", "}"),
		placeholder:   "p",
		timeBuffer:    new(bytes.Buffer),
		formatBuffer:  new(bytes.Buffer),
		messageBuffer: new(bytes.Buffer),
	}

	formatter.setFuncs()

	return formatter
}

// SetPlaceholder sets placeholder string prefix used for automatic and
// positional placeholders to format log message
func (formatter *Formatter) SetPlaceholder(placeholder string) *Formatter {
	formatter.placeholder = placeholder
	return formatter
}

// GetPlaceholder returns placeholder string prefix used for automatic and
// positional placeholders to format log message
func (formatter *Formatter) GetPlaceholder() string {
	return formatter.placeholder
}

// AddFuncs adds template functions to format log message
func (formatter *Formatter) AddFuncs(funcs FormatterFuncs) *Formatter {
	formatter.template.Funcs(template.FuncMap(funcs))
	return formatter
}

// GetRecord returns assigned log record object to formatter
func (formatter *Formatter) GetRecord() *Record {
	return formatter.record
}

// SetFormat sets format string used for formatting log message
func (formatter *Formatter) SetFormat(format string) *Formatter {
	formatter.format = format
	return formatter
}

// GetFormat returns format string used for formatting log message
func (formatter *Formatter) GetFormat() string {
	return formatter.format
}

// SetDateFormat sets format string used for formatting date in log message
func (formatter *Formatter) SetDateFormat(dateFormat string) *Formatter {
	formatter.dateFormat = dateFormat
	return formatter
}

// GetDateFormat returns format string used for formatting date in log message
func (formatter *Formatter) GetDateFormat() string {
	return formatter.format
}

// Format returns formatted log message string based on provided log record
// object
func (formatter *Formatter) Format(record *Record) string {
	formatter.record = record

	return formatter.formatString(
		formatter.template,
		formatter.formatBuffer,
		formatter.format,
		nil,
	)
}

// FormatTime returns formatted date string based on provided log record object
func (formatter *Formatter) FormatTime(record *Record) string {
	formatter.record = record

	return formatter.formatString(
		formatter.template,
		formatter.timeBuffer,
		formatter.dateFormat,
		nil,
	)
}

// FormatMessage returns formatted user message string based on provided log
// record object
func (formatter *Formatter) FormatMessage(record *Record) string {
	message := record.Message

	if len(record.Arguments) > 0 {
		var object interface{}

		funcMap := make(template.FuncMap)

		funcMap[formatter.placeholder] = formatter.argumentAutomatic()

		for position, argument := range formatter.record.Arguments {
			placeholder := formatter.placeholder + strconv.Itoa(position)

			funcMap[placeholder] = formatter.argumentValue(argument)

			switch argument.(type) {
			case Named:
				for key, value := range argument.(Named) {
					funcMap[key] = formatter.argumentValue(value)
				}
			default:
				if reflect.TypeOf(argument).Kind() == reflect.Struct {
					if object == nil {
						object = argument
					}
				}
			}
		}

		message = formatter.formatString(
			template.New("").Delims("{", "}").Funcs(funcMap),
			formatter.messageBuffer,
			message,
			object,
		)
	}

	return message
}

// argumentValue returns closure that returns log argument used in log message
func (formatter *Formatter) argumentValue(argument interface{}) func() interface{} {
	return func() interface{} {
		return argument
	}
}

// argumentAutomatic returns closure that returns log argument from automatic
// placeholder used in log message
func (formatter *Formatter) argumentAutomatic() func() interface{} {
	position := 0
	arguments := len(formatter.record.Arguments)

	return func() interface{} {
		var argument interface{}

		if position < arguments {
			argument = formatter.record.Arguments[position]
			position++
		}

		return argument
	}
}

// formatString returns formatted string
func (formatter *Formatter) formatString(template *template.Template, buffer *bytes.Buffer, format string, object interface{}) string {
	var err error

	template, err = template.Parse(format)

	if err != nil {
		NewRuntimeError("cannot parse text template", err).Print()
		return format
	}

	buffer.Reset()

	err = template.Execute(buffer, object)

	if err != nil {
		NewRuntimeError("cannot execute text template", err).Print()
		return format
	}

	return buffer.String()
}

// setFuncs sets default template functions used to formatting log message
func (formatter *Formatter) setFuncs() {
	formatter.template.Funcs(template.FuncMap{
		"executable": func() string {
			return filepath.Base(os.Args[0])
		},
		"getEnv": func(key string) string {
			return os.Getenv(key)
		},
		"expandEnv": func(name string) string {
			return os.ExpandEnv(name)
		},
		"date": func() string {
			return formatter.FormatTime(formatter.record)
		},
		"message": func() string {
			return formatter.FormatMessage(formatter.record)
		},
		"levelValue": func() int {
			return formatter.record.Level.Value
		},
		"level": func() string {
			return strings.ToLower(formatter.record.Level.Name)
		},
		"Level": func() string {
			return strings.Title(strings.ToLower(formatter.record.Level.Name))
		},
		"LEVEL": func() string {
			return strings.ToUpper(formatter.record.Level.Name)
		},
		"iso8601": func() string {
			return formatter.record.Time.Format(time.RFC3339)
		},
		"id": func() string {
			return formatter.record.ID
		},
		"gid": func() int {
			return os.Getgid()
		},
		"pid": func() int {
			return os.Getpid()
		},
		"ppid": func() int {
			return os.Getppid()
		},
		"name": func() string {
			return formatter.record.Name
		},
		"host": func() string {
			return formatter.record.Address
		},
		"hostname": func() string {
			return formatter.record.Hostname
		},
		"address": func() string {
			return formatter.record.Address
		},
		"nanosecond": func() string {
			return fmt.Sprintf("%09d", formatter.record.Time.Nanosecond())
		},
		"microsecond": func() string {
			return fmt.Sprintf("%06d", formatter.record.Time.Nanosecond()/1e3)
		},
		"milisecond": func() string {
			return fmt.Sprintf("%03d", formatter.record.Time.Nanosecond()/1e6)
		},
		"second": func() string {
			return fmt.Sprintf("%02d", formatter.record.Time.Second())
		},
		"minute": func() string {
			return fmt.Sprintf("%02d", formatter.record.Time.Minute())
		},
		"hour": func() string {
			return fmt.Sprintf("%02d", formatter.record.Time.Hour())
		},
		"day": func() string {
			return fmt.Sprintf("%02d", formatter.record.Time.Day())
		},
		"month": func() string {
			return fmt.Sprintf("%02d", formatter.record.Time.Month())
		},
		"YEAR": func() string {
			return fmt.Sprintf("%02d", formatter.record.Time.Year()%100)
		},
		"year": func() int {
			return formatter.record.Time.Year()
		},
		"file": func() string {
			return formatter.record.File.Name
		},
		"line": func() int {
			return formatter.record.File.Line
		},
		"function": func() string {
			return formatter.record.File.Function
		},
	})
}
