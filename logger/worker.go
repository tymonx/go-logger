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
	"os"
	"path/filepath"
	"sync"
	"time"
)

// These constants define default values for Worker.
const (
	DefaultQueueLength = 4096
)

// A Worker represents an active logger worker thread. It handles formatting
// received log messages and I/O operations.
type Worker struct {
	flush   chan *sync.WaitGroup
	records chan *Record
	mutex   sync.RWMutex
}

var gWorkerOnce sync.Once   // nolint:gochecknoglobals
var gWorkerInstance *Worker // nolint:gochecknoglobals

// NewWorker creates a new Worker object.
func NewWorker() *Worker {
	worker := &Worker{
		flush:   make(chan *sync.WaitGroup, 1),
		records: make(chan *Record, DefaultQueueLength),
	}

	go worker.run()

	return worker
}

// GetWorker returns logger worker instance. First call to it creates and
// starts logger worker thread.
func GetWorker() *Worker {
	gWorkerOnce.Do(func() {
		gWorkerInstance = NewWorker()
	})

	return gWorkerInstance
}

// SetQueueLength sets logger worker thread queue length for log messages.
func (w *Worker) SetQueueLength(length int) *Worker {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if length <= 0 {
		length = DefaultQueueLength
	}

	if cap(w.records) != length {
		w.records = make(chan *Record, length)
	}

	return w
}

// Flush flushes all log messages.
func (w *Worker) Flush() *Worker {
	flush := new(sync.WaitGroup)

	flush.Add(1)
	w.flush <- flush
	flush.Wait()

	return w
}

// Run processes all incoming log messages from loggers. It emits received log
// records to all added log handlers for specific logger.
func (w *Worker) run() {
	for {
		select {
		case flush := <-w.flush:
			for records := len(w.records); records > 0; records-- {
				record := <-w.records

				if record != nil {
					w.emit(record.logger, record)
				}
			}

			if flush != nil {
				flush.Done()
			}
		case record := <-w.records:
			if record != nil {
				w.emit(record.logger, record)
			}
		}
	}
}

// emit prepares provided log record and it dispatches to all added log
// handlers for further formatting and specific I/O implementation operations.
func (*Worker) emit(logger *Logger, record *Record) {
	var err error

	record.Type = DefaultTypeName
	record.File.Name = filepath.Base(record.File.Path)
	record.File.Function = filepath.Base(record.File.Function)
	record.Timestamp.Created = record.Time.Format(time.RFC3339)

	record.Address, err = getAddress()

	if err != nil {
		printError(NewRuntimeError("cannot get local IP address", err))
	}

	record.Hostname, err = getHostname()

	if err != nil {
		printError(NewRuntimeError("cannot get local hostname", err))
	}

	logger.mutex.RLock()
	defer logger.mutex.RUnlock()

	record.Name = logger.name
	record.ID, err = logger.idGenerator.Generate()

	if err != nil {
		printError(NewRuntimeError("cannot generate ID", err))
	}

	if record.Name == "" {
		record.Name = filepath.Base(os.Args[0])
	}

	for _, handler := range logger.handlers {
		min, max := handler.GetLevelRange()

		if (record.Level.Value >= min) && (record.Level.Value <= max) {
			err = handler.Emit(record)

			if err != nil {
				printError(NewRuntimeError("cannot emit record", err))
			}
		}
	}
}
