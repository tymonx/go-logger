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
	"io"
	"net"
	"strconv"
)

// These constants define default values for syslog.
const (
	DefaultSyslogPort    = 514
	DefaultSyslogVersion = 1
	DefaultSyslogNetwork = "tcp"
	DefaultSyslogAddress = "localhost"
	DefaultSyslogFormat  = "<{syslogPriority}>{syslogVersion} {iso8601} {address} {name} {pid} {id} - " +
		"{file}:{line}:{function}(): {message}"
	DefaultSyslogFacility = 1
)

// A Syslog represents a log handler object for logging messages to running
// Syslog server.
type Syslog struct {
	port     int
	version  int
	network  string
	address  string
	facility int
	stream   *Stream
}

// NewSyslog creates a new Syslog log handler object.
func NewSyslog() *Syslog {
	s := &Syslog{
		port:     DefaultSyslogPort,
		version:  DefaultSyslogVersion,
		network:  DefaultSyslogNetwork,
		address:  DefaultSyslogAddress,
		facility: DefaultSyslogFacility,
		stream:   NewStream(),
	}

	s.stream.GetFormatter().SetFormat(DefaultSyslogFormat)
	s.stream.SetOpener(s)

	return s
}

// Open opens new connection.
func (s *Syslog) Open() (io.WriteCloser, error) {
	return net.Dial(s.network, s.address+":"+strconv.Itoa(s.port))
}

// Enable enables log handler.
func (s *Syslog) Enable() Handler {
	return s.stream.Enable()
}

// Disable disabled log handler.
func (s *Syslog) Disable() Handler {
	return s.stream.Disable()
}

// IsEnabled returns if log handler is enabled.
func (s *Syslog) IsEnabled() bool {
	return s.stream.IsEnabled()
}

// SetFormatter sets Formatter.
func (s *Syslog) SetFormatter(formatter *Formatter) Handler {
	return s.stream.SetFormatter(formatter)
}

// GetFormatter returns Formatter.
func (s *Syslog) GetFormatter() *Formatter {
	return s.stream.GetFormatter()
}

// SetLevel sets log level.
func (s *Syslog) SetLevel(level int) Handler {
	return s.stream.SetLevel(level)
}

// SetMinimumLevel sets minimum log level.
func (s *Syslog) SetMinimumLevel(level int) Handler {
	return s.stream.SetMinimumLevel(level)
}

// GetMinimumLevel returns minimum log level.
func (s *Syslog) GetMinimumLevel() int {
	return s.stream.GetMinimumLevel()
}

// SetMaximumLevel sets maximum log level.
func (s *Syslog) SetMaximumLevel(level int) Handler {
	return s.stream.SetMaximumLevel(level)
}

// GetMaximumLevel returns maximum log level.
func (s *Syslog) GetMaximumLevel() int {
	return s.stream.GetMaximumLevel()
}

// SetLevelRange sets minimum and maximum log level values.
func (s *Syslog) SetLevelRange(min, max int) Handler {
	return s.stream.SetLevelRange(min, max)
}

// GetLevelRange returns minimum and maximum log level values.
func (s *Syslog) GetLevelRange() (min, max int) {
	return s.stream.GetLevelRange()
}

// SetPort sets port number that is used to communicate with Syslog server.
func (s *Syslog) SetPort(port int) *Syslog {
	s.stream.Lock()
	defer s.stream.Unlock()

	if port <= 0 {
		port = DefaultSyslogPort
	}

	if s.port != port {
		s.port = port
		s.stream.Reopen()
	}

	return s
}

// GetPort returns port number that is used to communicate with Syslog server.
func (s *Syslog) GetPort() int {
	s.stream.RLock()
	defer s.stream.RUnlock()

	return s.port
}

// SetNetwork sets network type like "udp" or "tcp" that is used to communicate
// with Syslog server.
func (s *Syslog) SetNetwork(network string) *Syslog {
	s.stream.Lock()
	defer s.stream.Unlock()

	if network == "" {
		network = DefaultSyslogNetwork
	}

	if s.network != network {
		s.network = network
		s.stream.Reopen()
	}

	return s
}

// GetNetwork returns network type like "udp" or "tcp" that is used to
// communicate with Syslog server.
func (s *Syslog) GetNetwork() string {
	s.stream.RLock()
	defer s.stream.RUnlock()

	return s.network
}

// SetAddress sets IP address or hostname that is used to communicate with
// Syslog server.
func (s *Syslog) SetAddress(address string) *Syslog {
	s.stream.Lock()
	defer s.stream.Unlock()

	if address == "" {
		address = DefaultSyslogAddress
	}

	if s.address != address {
		s.address = address
		s.stream.Reopen()
	}

	return s
}

// GetAddress returns IP address or hostname that is used to communicate with
// Syslog server.
func (s *Syslog) GetAddress() string {
	s.stream.RLock()
	defer s.stream.RUnlock()

	return s.network
}

// Emit logs messages from Logger to Syslog server.
func (s *Syslog) Emit(record *Record) error {
	s.stream.GetFormatter().AddFuncs(s.getRecordFuncs(record))

	return s.stream.Emit(record)
}

// Close closes communication to Syslog server.
func (s *Syslog) Close() error {
	return s.stream.Close()
}

// setFormatterFuncs sets template functions that are specific for Syslog log
// messages.
func (s *Syslog) getRecordFuncs(record *Record) FormatterFuncs {
	return FormatterFuncs{
		"syslogVersion": func() int {
			return s.version
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
				if level <= record.Level.Value {
					severity = i
					break
				}
			}

			return ((0x1F & s.facility) << 3) | (0x07 & severity)
		},
	}
}
