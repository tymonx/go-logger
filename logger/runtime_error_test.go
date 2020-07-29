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

func TestRuntimeErrorNew(test *testing.T) {
	err := logger.NewRuntimeError("test")

	if err == nil {
		test.Error("NewRuntimeError() returns nil")
	}

	if err := err.Unwrap(); err != nil {
		test.Error("Unwrap() returns an unexpected error", err)
	}
}

func TestRuntimeErrorNoArguments(test *testing.T) {
	err := logger.NewRuntimeError("test")

	want := "runtime_error_test.go:36:logger_test.TestRuntimeErrorNoArguments(): test"

	if err == nil {
		test.Error("NewRuntimeError() returns nil")
	}

	if err.Error() != want {
		test.Error("Error() =", err.Error(), "; want", want)
	}

	if err := err.Unwrap(); err != nil {
		test.Error("Unwrap() returns an unexpected error", err)
	}
}

func TestRuntimeErrorAutoPlacedArguments(test *testing.T) {
	err := logger.NewRuntimeError("test", 3, "hello", "world", nil, 0)

	want := "runtime_error_test.go:54:logger_test.TestRuntimeErrorAutoPlacedArguments(): test 3 hello world <nil> 0"

	if err == nil {
		test.Error("NewRuntimeError() returns nil")
	}

	if err.Error() != want {
		test.Error("Error() =", err.Error(), "; want", want)
	}

	if err := err.Unwrap(); err != nil {
		test.Error("Unwrap() returns an unexpected error", err)
	}
}

func TestRuntimeErrorError(test *testing.T) {
	err := logger.NewRuntimeError("test", testError)

	want := "runtime_error_test.go:72:logger_test.TestRuntimeErrorError(): test My test error"

	if err == nil {
		test.Error("NewRuntimeError() returns nil")
	}

	if err.Error() != want {
		test.Error("Error() =", err.Error(), "; want", want)
	}

	if err.Unwrap() == nil {
		test.Error("Unwrap() returns nil")
	}
}

func TestRuntimeErrorErrors(test *testing.T) {
	errs := []interface{}{
		testError,
		testError,
	}

	err := logger.NewRuntimeError("test", errs...)

	want := "runtime_error_test.go:95:logger_test.TestRuntimeErrorErrors(): test My test error My test error"

	if err == nil {
		test.Error("NewRuntimeError() returns nil")
	}

	if err.Error() != want {
		test.Error("Error() =", err.Error(), "; want", want)
	}

	if err.Unwrap() == nil {
		test.Error("Unwrap() returns nil")
	}
}
