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

package logger_test

import (
	"testing"

	"gitlab.com/tymonx/go-logger/logger"
)

func TestFormatterNew(test *testing.T) {
	formatter := logger.NewFormatter()

	if formatter == nil {
		test.Error("NewFormatter() returns nil")
	}
}

func TestFormatterFormatMessage(test *testing.T) {
	var err error

	var message string

	record := &logger.Record{}

	formatter := logger.NewFormatter()

	if message, err = formatter.FormatMessage(record); err != nil {
		test.Error("FormatMessage() returns an unexpected error", err)
	}

	if message != "" {
		test.Error("FormatMessage() returns an unexpected message", message)
	}
}

func TestFormatterFormatMessageNoArguments(test *testing.T) {
	var err error

	var message string

	want := "{x}"

	record := &logger.Record{
		Message: want,
	}

	formatter := logger.NewFormatter()

	if message, err = formatter.FormatMessage(record); err != nil {
		test.Error("FormatMessage() returns an unexpected error", err)
	}

	if message != want {
		test.Error("FormatMessage() =", message, "; want", want)
	}
}

func TestFormatterFormatMessageNamedArguments(test *testing.T) {
	var err error

	var message string

	record := &logger.Record{
		Message: "{name}",
		Arguments: []interface{}{
			logger.Named{
				"name": testMessage,
			},
		},
	}

	formatter := logger.NewFormatter()

	if message, err = formatter.FormatMessage(record); err != nil {
		test.Error("FormatMessage() returns an unexpected error", err)
	}

	if message != testMessage {
		test.Error("FormatMessage() =", message, "; want", testMessage)
	}
}

func TestFormatterFormatMessageObjectArguments(test *testing.T) {
	var err error

	var message string

	var object struct {
		Name string
	}

	object.Name = testMessage

	record := &logger.Record{
		Message: "{.Name}",
		Arguments: []interface{}{
			object,
		},
	}

	formatter := logger.NewFormatter()

	if message, err = formatter.FormatMessage(record); err != nil {
		test.Error("FormatMessage() returns an unexpected error", err)
	}

	if message != testMessage {
		test.Error("FormatMessage() =", message, "; want", testMessage)
	}
}

func TestFormatterFormatMessageAutoAppend(test *testing.T) {
	var err error

	var message string

	want := "Test message 5 4 hello world 0.3 <nil> " + testError.Error() + " [x y]"

	record := &logger.Record{
		Message: "Test message {named2}",
		Arguments: []interface{}{
			4,
			"hello world",
			0.3,
			nil,
			logger.Named{
				"named1": 3,
				"named2": 5,
			},
			testError,
			[]error{
				Error("x"),
				Error("y"),
			},
		},
	}

	formatter := logger.NewFormatter()

	if message, err = formatter.FormatMessage(record); err != nil {
		test.Error("FormatMessage() returns an unexpected error", err)
	}

	if message != want {
		test.Error("FormatMessage() =", message, "; want", want)
	}
}

func TestFormatterFormatMessageErrors(test *testing.T) {
	var err error

	var message string

	errs := []interface{}{
		testError,
		testError,
	}

	want := testError.Error() + " " + testError.Error()

	record := &logger.Record{
		Message:   "",
		Arguments: errs,
	}

	formatter := logger.NewFormatter()

	if message, err = formatter.FormatMessage(record); err != nil {
		test.Error("FormatMessage() returns an unexpected error", err)
	}

	if message != want {
		test.Error("FormatMessage() =", message, "; want", want)
	}
}
