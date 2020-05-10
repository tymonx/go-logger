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

// Handler defines interface for log handlers
type Handler interface {
	Emit(record *Record) error
	GetLevelRange() (min int, max int)
	Close() error
}

// Handlers defines map of log handlers
type Handlers map[string]Handler

// HandlerConstructor creates specific log handler
type HandlerConstructor func() Handler

var gHandlerConstructors = make(map[string]HandlerConstructor)

// RegisterHandler registers log handler under provided name
func RegisterHandler(name string, constructor HandlerConstructor) {
	gMutex.Lock()
	defer gMutex.Unlock()

	gHandlerConstructors[name] = constructor
}

// CreateHandler creates registered log handler by provided name
func CreateHandler(name string) Handler {
	gMutex.Lock()
	defer gMutex.Unlock()

	return gHandlerConstructors[name]()
}
