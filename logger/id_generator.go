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

// IDGenerator function type that returns generated ID used in log messages
type IDGenerator func() string

// gIDGenerators contains all registered ID generators
var gIDGenerators = make(map[string]IDGenerator)

// RegisterIDGenerator registers ID generator function under provided
// identifier name
func RegisterIDGenerator(name string, idGenerator IDGenerator) {
	gMutex.Lock()
	defer gMutex.Unlock()

	gIDGenerators[name] = idGenerator
}

// CreateIDGenerator returns registered ID generator function by provided
// identifier name
func CreateIDGenerator(name string) IDGenerator {
	gMutex.Lock()
	defer gMutex.Unlock()

	return gIDGenerators[name]
}
