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

// Package logger implements logging package. It defines a type, Logger, with
// methods for formatting output. Each logging operations creates and sends
// lightweight not formatted log message to separate worker thread. It offloads
// main code from unnecessary resource consuming formatting and I/O operations.
// On default it supports many different log handlers like logging to standard
// output, error output, file, stream or syslog.
package logger
