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
	"math/rand"
	"os"
	"testing"
)

const (
	coverageLevel = 0.8
)

func compare(test *testing.T, captured interface{}, expected interface{}) {
	if captured != expected {
		test.Error("captured =", captured, "want", expected)
	}
}

func Example() {
	logger := New()

	for _, handler := range logger.GetHandlers() {
		handler.GetFormatter().SetDateFormat("2020")
	}

	logger.Info("Hello from logger!")
	logger.Info("Automatic placeholders {p} {p} {p}", 1, 2, 3)
	logger.Info("Positional placeholders {p2} {p1} {p0}", 1, 2, 3)

	logger.Info("Named placeholders {z} {y} {x}", Named{
		"x": 1,
		"y": 2,
		"z": 3,
	})

	logger.Info("Object placeholders {.Z} {.Y} {.X}", struct {
		X, Y, Z int
	}{
		X: 1,
		Y: 2,
		Z: 3,
	})

	logger.Flush()

	// Output:
	// 2020 - Info     - logger_test.go:41:logger.Example(): Hello from logger!
	// 2020 - Info     - logger_test.go:42:logger.Example(): Automatic placeholders 1 2 3
	// 2020 - Info     - logger_test.go:43:logger.Example(): Positional placeholders 3 2 1
	// 2020 - Info     - logger_test.go:45:logger.Example(): Named placeholders 3 2 1
	// 2020 - Info     - logger_test.go:51:logger.Example(): Object placeholders 3 2 1
}

func TestMain(main *testing.M) {
	defer Close()

	returnCode := main.Run()

	if (returnCode == 0) && (testing.CoverMode() != "") {
		coverage := testing.Coverage()

		if coverage < coverageLevel {
			fmt.Fprintf(
				os.Stderr,
				"Tests passed but coverage failed at %.1f%%\n", coverage*100,
			)
			returnCode = -2
		}
	}

	os.Exit(returnCode)
}

func TestNew(test *testing.T) {
	var object interface{} = New()

	if object == nil {
		test.Fatal("invalid pointer value")
	}

	logger, ok := object.(*Logger)

	if !ok {
		test.Fatal("invalid pointer type")
	}

	if logger == nil {
		test.Fatal("invalid pointer value")
	}
}

func TestGetHandlers(test *testing.T) {
	handlers := New().GetHandlers()

	compare(test, len(handlers), 2)

	if _, ok := handlers["stdout"]; !ok {
		test.Error("handlers[\"stdout\"] doesn't exist")
	}

	if _, ok := handlers["stderr"]; !ok {
		test.Error("handlers[\"stderr\"] doesn't exist")
	}
}

func TestGetIDGenerator(test *testing.T) {
	idGenerator := New().GetIDGenerator()

	if idGenerator == nil {
		test.Error("idGenerator is nil")
	}
}

func TestSetName(test *testing.T) {
	logger := New()

	for _, expected := range []string{"logger", "log-2", "logx", "ll"} {
		captured := logger.SetName(expected).GetName()

		compare(test, captured, expected)
	}
}

func TestSetErrorCode(test *testing.T) {
	logger := New()

	for count := 0; count < 10; count++ {
		expected := rand.Int()
		captured := logger.SetErrorCode(expected).GetErrorCode()

		compare(test, captured, expected)
	}
}

func TestGetErrorCode(test *testing.T) {
	captured := New().GetErrorCode()

	if captured != DefaultErrorCode {
		compare(test, captured, DefaultErrorCode)
	}
}
