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
	"net"
	"os"
	"runtime"
	"sync"
)

// Named is used as named string placeholders for logger functions
type Named map[string]interface{}

// gMutex global read/write mutex
var gMutex = new(sync.RWMutex)

// getHostname returns local hostname
func getHostname() string {
	hostname, err := os.Hostname()

	if err != nil {
		hostname = "localhost"
		NewRuntimeError("cannot get hostname", err).Print()
	}

	return hostname
}

// getAddress returns local IP address
func getAddress() string {
	var address string

	connection, err := net.Dial("udp", "8.8.8.8:80")

	if err == nil {
		defer connection.Close()

		address = connection.LocalAddr().(*net.UDPAddr).IP.String()
	} else {
		address = "127.0.0.1"
		NewRuntimeError("cannot connect to primary Google DNS", err).Print()
	}

	return address
}

// getPathLineFunction returns absolute file path, file line number and function
// name
func getPathLineFunction(skip int) (path string, line int, function string) {
	var pc uintptr
	var ok bool

	pc, path, line, ok = runtime.Caller(skip)
	function = runtime.FuncForPC(pc).Name()

	if !ok {
		NewRuntimeError("cannot recover runtime information", nil).Print()
	}

	return
}
