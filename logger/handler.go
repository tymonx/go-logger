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

// Handler defines interface for log handlers.
type Handler interface {
	SetFormatter(formatter *Formatter) Handler

	GetFormatter() *Formatter

	SetMinimumLevel(level int) Handler

	GetMinimumLevel() int

	SetMaximumLevel(level int) Handler

	GetMaximumLevel() int

	SetLevelRange(min, max int) Handler

	GetLevelRange() (min, max int)

	Emit(record *Record) error

	Close() error
}

// Handlers defines map of log handlers.
type Handlers map[string]Handler
