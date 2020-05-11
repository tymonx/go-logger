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
	"sync"
)

// These constants define default values for syslog
const (
	DefaultSyslogPort     = 514
	DefaultSyslogVersion  = 1
	DefaultSyslogNetwork  = "tcp"
	DefaultSyslogAddress  = "localhost"
	DefaultSyslogFormat   = "<{syslogPriority}>{syslogVersion} {iso8601} {address} {name} {pid} {id} - {file}:{line}:{function}(): {message}"
	DefaultSyslogFacility = 1
)

// A Syslog represents a log handler object for logging messages to running
// Syslog server
type Syslog struct {
	port         int
	mutex        sync.RWMutex
	version      int
	network      string
	address      string
	facility     int
	formatter    *Formatter
	connection   net.Conn
	minimumLevel int
	maximumLevel int
	reconnect    bool
}

// NewSyslog creates a new Syslog log handler object
func NewSyslog() *Syslog {
	syslog := &Syslog{
		port:         DefaultSyslogPort,
		version:      DefaultSyslogVersion,
		network:      DefaultSyslogNetwork,
		address:      DefaultSyslogAddress,
		facility:     DefaultSyslogFacility,
		formatter:    NewFormatter().SetFormat(DefaultSyslogFormat),
		minimumLevel: MinimumLevel,
		maximumLevel: MaximumLevel,
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

// SetFormatter sets Formatter
func (syslog *Syslog) SetFormatter(formatter *Formatter) Handler {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	syslog.formatter = formatter

	return syslog
}

// GetFormatter returns Formatter
func (syslog *Syslog) GetFormatter() *Formatter {
	syslog.mutex.RLock()
	defer syslog.mutex.RUnlock()

	return syslog.formatter
}

// SetMinimumLevel sets minimum log level
func (syslog *Syslog) SetMinimumLevel(level int) Handler {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	syslog.minimumLevel = level

	return syslog
}

// GetMinimumLevel returns minimum log level
func (syslog *Syslog) GetMinimumLevel() int {
	syslog.mutex.RLock()
	defer syslog.mutex.RUnlock()

	return syslog.minimumLevel
}

// SetMaximumLevel sets maximum log level
func (syslog *Syslog) SetMaximumLevel(level int) Handler {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	syslog.maximumLevel = level

	return syslog
}

// GetMaximumLevel returns maximum log level
func (syslog *Syslog) GetMaximumLevel() int {
	syslog.mutex.RLock()
	defer syslog.mutex.RUnlock()

	return syslog.maximumLevel
}

// SetLevelRange sets minimum and maximum log level values
func (syslog *Syslog) SetLevelRange(min int, max int) Handler {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	syslog.minimumLevel = min
	syslog.maximumLevel = max

	return syslog
}

// GetLevelRange returns minimum and maximum log level values
func (syslog *Syslog) GetLevelRange() (min int, max int) {
	syslog.mutex.RLock()
	defer syslog.mutex.RUnlock()

	return syslog.minimumLevel, syslog.maximumLevel
}

// SetPort sets port number that is used to communicate with Syslog server
func (syslog *Syslog) SetPort(port int) *Syslog {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	if port <= 0 {
		port = DefaultSyslogPort
	}

	if syslog.port != port {
		syslog.port = port
		syslog.reconnect = true
	}

	return syslog
}

// GetPort returns port number that is used to communicate with Syslog server
func (syslog *Syslog) GetPort() int {
	syslog.mutex.RLock()
	defer syslog.mutex.RUnlock()

	return syslog.port
}

// SetNetwork sets network type like "udp" or "tcp" that is used to communicate
// with Syslog server
func (syslog *Syslog) SetNetwork(network string) *Syslog {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	if network == "" {
		network = DefaultSyslogNetwork
	}

	if syslog.network != network {
		syslog.network = network
		syslog.reconnect = true
	}

	return syslog
}

// GetNetwork returns network type like "udp" or "tcp" that is used to
// communicate with Syslog server
func (syslog *Syslog) GetNetwork() string {
	syslog.mutex.RLock()
	defer syslog.mutex.RUnlock()

	return syslog.network
}

// SetAddress sets IP address or hostname that is used to communicate with
// Syslog server
func (syslog *Syslog) SetAddress(address string) *Syslog {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	if address == "" {
		address = DefaultSyslogAddress
	}

	if syslog.address != address {
		syslog.address = address
		syslog.reconnect = true
	}

	return syslog
}

// GetAddress returns IP address or hostname that is used to communicate with
// Syslog server
func (syslog *Syslog) GetAddress() string {
	syslog.mutex.RLock()
	defer syslog.mutex.RUnlock()

	return syslog.network
}

// Emit logs messages from Logger to Syslog server
func (syslog *Syslog) Emit(record *Record) error {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	if syslog.reconnect {
		err := syslog.close()

		if err != nil {
			return NewRuntimeError("cannot close connection to syslog", err)
		}

		syslog.reconnect = false
	}

	if syslog.connection == nil {
		err := syslog.connect()

		if err != nil {
			return NewRuntimeError("cannot open connection to syslog", err)
		}
	}

	if syslog.connection != nil {
		message, err := syslog.formatter.Format(record)

		if err != nil {
			return NewRuntimeError("cannot fomat record", err)
		}

		_, err = fmt.Fprintln(syslog.connection, message)

		if err != nil {
			return NewRuntimeError("cannot write to syslog", err)
		}
	}

	return nil
}

// Close closes communication to Syslog server
func (syslog *Syslog) Close() error {
	syslog.mutex.Lock()
	defer syslog.mutex.Unlock()

	err := syslog.close()

	if err != nil {
		return NewRuntimeError("cannot close connection to syslog", err)
	}

	return nil
}

func (syslog *Syslog) connect() error {
	var err error

	syslog.connection, err = net.Dial(
		syslog.network,
		syslog.address+":"+strconv.Itoa(syslog.port),
	)

	if err != nil {
		syslog.connection = nil
	}

	return err
}

func (syslog *Syslog) close() error {
	var err error

	if syslog.connection != nil {
		err = syslog.connection.Close()
		syslog.connection = nil
	}

	return err
}

// setFormatterFuncs sets template functions that are specific for Syslog log
// messages
func (syslog *Syslog) setFormatterFuncs() {
	syslog.formatter.AddFuncs(map[string]interface{}{
		"syslogVersion": func() int {
			return syslog.version
		},
		"syslogPriority": func() int {
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

			return ((0x1F & syslog.facility) << 3) | (0x07 & severity)
		},
	})
}
