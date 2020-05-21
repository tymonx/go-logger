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
	"net"
	"os"
)

// Named is used as named string placeholders for logger functions.
type Named map[string]interface{}

// getHostname returns local hostname.
func getHostname() (string, error) {
	hostname, err := os.Hostname()

	if err != nil {
		return "localhost", NewRuntimeError("cannot get hostname", err)
	}

	return hostname, nil
}

// getAddress returns local IP address.
func getAddress() (string, error) {
	connection, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return "127.0.0.1", NewRuntimeError("cannot connect to primary Google DNS", err)
	}

	defer func() {
		err := connection.Close()

		if err != nil {
			printError(NewRuntimeError("cannot close UDP connection", err))
		}
	}()

	return connection.LocalAddr().(*net.UDPAddr).IP.String(), nil
}

// printError prints error to error output.
func printError(err error) {
	fmt.Fprintln(os.Stderr, "Logger error:", err)
}
