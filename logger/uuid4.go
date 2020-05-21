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
	"crypto/rand"
	"encoding/hex"
	"io"
)

// An UUID4 represents uui4 generator.
type UUID4 struct{}

// NewUUID4 create a new UUID4 object.
func NewUUID4() *UUID4 {
	return &UUID4{}
}

// Generate generates new UUID4.
func (u *UUID4) Generate() (interface{}, error) {
	var uuid [16]byte

	var buffer [36]byte

	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return "", NewRuntimeError("cannot generate UUID4", err)
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10

	u.encodeHex(uuid, buffer[:])

	return string(buffer[:]), nil
}

// encodeHex encodes uuid to UUID format.
func (*UUID4) encodeHex(uuid [16]byte, buffer []byte) {
	hex.Encode(buffer, uuid[:4])
	buffer[8] = '-'

	hex.Encode(buffer[9:13], uuid[4:6])
	buffer[13] = '-'

	hex.Encode(buffer[14:18], uuid[6:8])
	buffer[18] = '-'

	hex.Encode(buffer[19:23], uuid[8:10])
	buffer[23] = '-'

	hex.Encode(buffer[24:], uuid[10:])
}
