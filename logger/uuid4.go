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
	"github.com/google/uuid"
)

// An UUID4 represents uui4 generator
type UUID4 struct{}

// NewUUID4 create a new UUID4 object
func NewUUID4() *UUID4 {
	return &UUID4{}
}

func init() {
	RegisterIDGenerator("uuid4", func() IDGenerator {
		return NewUUID4()
	})
}

// Generate generates new UUID4
func (uuid4 *UUID4) Generate() (interface{}, error) {
	var id string

	result, err := uuid.NewRandom()

	if err == nil {
		id = result.String()
	} else {
		return "", NewRuntimeError("cannot generate UUID4", err)
	}

	return id, nil
}
