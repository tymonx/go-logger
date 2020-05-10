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
	"encoding/json"
	"time"
)

// Arguments defines log arguments
type Arguments []interface{}

// Level defines log level information fields
type Level struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

// Timestamp defines log timestamp information fields
type Timestamp struct {
	Created string `json:"created"`
}

// Source defines log file information fields
type Source struct {
	Function string `json:"function"`
	Name     string `json:"name"`
	Path     string `json:"-"`
	Line     int    `json:"line"`
}

// Record defines log record fields created by Logger and it is used by
// Formatter to format log message based on these fields
type Record struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Time      time.Time `json:"-"`
	Level     Level     `json:"level"`
	Address   string    `json:"address"`
	Hostname  string    `json:"hostname"`
	Message   string    `json:"message"`
	File      Source    `json:"file"`
	Arguments Arguments `json:"arguments"`
	Timestamp Timestamp `json:"timestamp"`
	logger    *Logger
}

// ToJSON packs data to JSON
func (record *Record) ToJSON() ([]byte, error) {
	return json.Marshal(record)
}

// FromJSON unpacks data from JSON
func (record *Record) FromJSON(data []byte) error {
	return json.Unmarshal(data, record)
}

// GetMessage returns formatted message
func (record *Record) GetMessage() string {
	return NewFormatter().FormatMessage(record)
}
