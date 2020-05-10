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
	"strconv"
)

// These constants define default values for syslog
const (
	DefaultSyslogPort     = 514
	DefaultSyslogVersion  = 1
	DefaultSyslogNetwork  = "tcp"
	DefaultSyslogAddress  = "localhost"
	DefaultSyslogFormat   = "{syslogPriority}{syslogVersion} {iso8601} {address} {name} {pid} {id} - {file}:{line}:{function}(): {message}"
	DefaultSyslogFacility = 1
)

// A Syslog represents a log handler object for logging messages to running
// Syslog server
type Syslog struct {
	port       int
	version    int
	network    string
	address    string
	facility   int
	formatter  *Formatter
	connection net.Conn
}

// NewSyslog creates a new Syslog log handler object
func NewSyslog() *Syslog {
	syslog := &Syslog{
		port:      DefaultSyslogPort,
		version:   DefaultSyslogVersion,
		network:   DefaultSyslogNetwork,
		address:   DefaultSyslogAddress,
		facility:  DefaultSyslogFacility,
		formatter: NewFormatter().SetFormat(DefaultSyslogFormat),
	}

	syslog.setFormatterFuncs()

	return syslog
}

// init registers Syslog log handler
func init() {
	RegisterHandler("syslog", func() Handler {
		return NewSyslog()
	})
}

// GetLevelRange returns minimum and maximum log level values
func (syslog *Syslog) GetLevelRange() (min int, max int) {
	return TraceLevel, PanicLevel
}

// SetPort sets port number that is used to communicate with Syslog server
func (syslog *Syslog) SetPort(port int) *Syslog {
	if port > 0 {
		syslog.port = port
	} else {
		syslog.port = DefaultSyslogPort
	}

	return syslog
}

// GetPort returns port number that is used to communicate with Syslog server
func (syslog *Syslog) GetPort() int {
	return syslog.port
}

// SetNetwork sets network type like "udp" or "tcp" that is used to communicate
// with Syslog server
func (syslog *Syslog) SetNetwork(network string) *Syslog {
	if network != "" {
		syslog.network = network
	} else {
		syslog.network = DefaultSyslogNetwork
	}

	return syslog
}

// GetNetwork returns network type like "udp" or "tcp" that is used to
// communicate with Syslog server
func (syslog *Syslog) GetNetwork() string {
	return syslog.network
}

// SetAddress sets IP address or hostname that is used to communicate with
// Syslog server
func (syslog *Syslog) SetAddress(address string) *Syslog {
	if address != "" {
		syslog.address = address
	} else {
		syslog.address = DefaultSyslogAddress
	}

	return syslog
}

// GetAddress returns IP address or hostname that is used to communicate with
// Syslog server
func (syslog *Syslog) GetAddress() string {
	return syslog.network
}

// Emit logs messages from Logger to Syslog server
func (syslog *Syslog) Emit(record *Record) error {
	if syslog.connection == nil {
		var err error

		syslog.connection, err = net.Dial(
			syslog.network,
			syslog.address+":"+strconv.Itoa(syslog.port),
		)

		if err != nil {
			syslog.connection = nil
			return NewRuntimeError("cannot connect to syslog", err)
		}
	}

	if syslog.connection != nil {
		_, err := fmt.Fprintln(
			syslog.connection,
			syslog.formatter.Format(record),
		)

		if err != nil {
			return NewRuntimeError("cannot write to syslog", err)
		}
	}

	return nil
}

// Close closes communication to Syslog server
func (syslog *Syslog) Close() error {
	var err error

	if syslog.connection != nil {
		err = syslog.connection.Close()

		syslog.connection = nil

		if err != nil {
			return NewRuntimeError("cannot close connection to syslog", err)
		}
	}

	return nil
}

// setFormatterFuncs sets template functions that are specific for Syslog log
// messages
func (syslog *Syslog) setFormatterFuncs() {
	syslog.formatter.AddFuncs(map[string]interface{}{
		"syslogVersion": func() int {
			return syslog.version
		},
		"syslogPriority": func() string {
			severities := [8]int{
				FatalLevel,
				AlertLevel,
				CriticalLevel,
				ErrorLevel,
				WarningLevel,
				NoticeLevel,
				InfoLevel,
				DebugLevel,
			}

			severity := 0

			for i, level := range severities {
				if level <= syslog.formatter.GetRecord().Level.Value {
					severity = i
					break
				}
			}

			priority := ((0x1F & syslog.facility) << 3) | (0x07 & severity)

			return fmt.Sprintf("<%d>", priority)
		},
	})
}
