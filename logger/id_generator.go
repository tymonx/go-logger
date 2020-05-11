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
	"sync"
)

// IDGenerator type that returns generated ID used in log messages
type IDGenerator interface {
	Generate() (id interface{}, err error)
}

// IDGeneratorConstructor defines ID generator constructor
type IDGeneratorConstructor func() IDGenerator

// gIDGenerators contains all registered ID generators
var gIDGeneratorConstructors = make(map[string]IDGeneratorConstructor)
var gIDGeneratorMutex sync.RWMutex

// RegisterIDGenerator registers ID generator under provided identifier name
func RegisterIDGenerator(name string, constructor IDGeneratorConstructor) {
	gIDGeneratorMutex.Lock()
	defer gIDGeneratorMutex.Unlock()

	gIDGeneratorConstructors[name] = constructor
}

// CreateIDGenerator returns registered ID generator by provided identifier name
func CreateIDGenerator(name string) (idGenerator IDGenerator, err error) {
	gIDGeneratorMutex.RLock()
	defer gIDGeneratorMutex.RUnlock()

	constructor, ok := gIDGeneratorConstructors[name]

	if ok {
		idGenerator = constructor()
	} else {
		err = NewRuntimeError("cannot create ID generator "+name, nil)
	}

	return
}
