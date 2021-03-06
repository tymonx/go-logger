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
	"crypto/rand"
	"math/big"
	"testing"

	"gitlab.com/tymonx/go-logger/logger"
)

func Example() {
	// To make testing this example more consistent, date must be constant
	for _, handler := range logger.GetHandlers() {
		handler.GetFormatter().SetDateFormat("2020")
	}

	logger.Info("Hello from logger!")
	logger.Info("Automatic placeholders {p} {p} {p}", 1, 2, 3)
	logger.Info("Positional placeholders {p2} {p1} {p0}", 1, 2, 3)

	logger.Info("Named placeholders {z} {y} {x}", logger.Named{
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
	// 2020 - Info     - logger_test.go:31:logger_test.Example(): Hello from logger!
	// 2020 - Info     - logger_test.go:32:logger_test.Example(): Automatic placeholders 1 2 3
	// 2020 - Info     - logger_test.go:33:logger_test.Example(): Positional placeholders 3 2 1
	// 2020 - Info     - logger_test.go:35:logger_test.Example(): Named placeholders 3 2 1
	// 2020 - Info     - logger_test.go:41:logger_test.Example(): Object placeholders 3 2 1
}

func TestNew(test *testing.T) {
	var object interface{} = logger.New()

	if object == nil {
		test.Fatal("invalid pointer value")
	}

	log, ok := object.(*logger.Logger)

	if !ok {
		test.Fatal("invalid pointer type")
	}

	if log == nil {
		test.Fatal("invalid pointer value")
	}
}

func TestGetHandlers(test *testing.T) {
	handlers := logger.New().GetHandlers()

	if len(handlers) != 2 {
		test.Errorf("len(handlers) = 2; want %d", len(handlers))
	}

	if _, ok := handlers["stdout"]; !ok {
		test.Error("handlers[\"stdout\"] doesn't exist")
	}

	if _, ok := handlers["stderr"]; !ok {
		test.Error("handlers[\"stderr\"] doesn't exist")
	}
}

func TestGetIDGenerator(test *testing.T) {
	idGenerator := logger.New().GetIDGenerator()

	if idGenerator == nil {
		test.Error("idGenerator is nil")
	}
}

func TestSetName(test *testing.T) {
	log := logger.New()

	for _, expected := range []string{"logger", "log-2", "logx", "ll"} {
		name := log.SetName(expected).GetName()

		if name != expected {
			test.Errorf("logger.SetName(%s); got %s", name, expected)
		}
	}
}

func TestSetErrorCode(test *testing.T) {
	log := logger.New()

	for count := 0; count < 10; count++ {
		n, err := rand.Int(rand.Reader, big.NewInt(1000))

		if err != nil {
			test.Fatal(err)
		}

		expected := int(n.Int64())
		errorCode := log.SetErrorCode(expected).GetErrorCode()

		if errorCode != expected {
			test.Errorf("logger.SetErrorCode(%d); got %d", errorCode, expected)
		}
	}
}

func TestGetErrorCode(test *testing.T) {
	errorCode := logger.New().GetErrorCode()

	if errorCode != logger.DefaultErrorCode {
		test.Errorf("logger.GetErrorCode() = %d; want %d", errorCode, logger.DefaultErrorCode)
	}
}
