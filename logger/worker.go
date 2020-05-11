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

// These constants define default values for Worker
const (
	DefaultQueueLength = 4096
)

// A Worker represents an active logger worker thread. It handles formatting
// received log messages and I/O operations
type Worker struct {
	flush   chan *sync.WaitGroup
	records chan *Record
	mutex   sync.RWMutex
}

var gWorkerOnce sync.Once
var gWorkerInstance *Worker

// NewWorker creates a new Worker object
func NewWorker() *Worker {
	worker := &Worker{
		flush:   make(chan *sync.WaitGroup, 1),
		records: make(chan *Record, DefaultQueueLength),
	}

	go worker.run()

	return worker
}

// GetWorker returns logger worker instance. First call to it creates and
// starts logger worker thread
func GetWorker() *Worker {
	gWorkerOnce.Do(func() {
		gWorkerInstance = NewWorker()
	})

	return gWorkerInstance
}

// SetQueueLength sets logger worker thread queue length for log messages
func (worker *Worker) SetQueueLength(length int) *Worker {
	worker.mutex.Lock()
	defer worker.mutex.Unlock()

	if length <= 0 {
		length = DefaultQueueLength
	}

	if cap(worker.records) != length {
		worker.records = make(chan *Record, length)
	}

	return worker
}

// Flush flushes all log messages
func (worker *Worker) Flush() *Worker {
	flush := new(sync.WaitGroup)

	flush.Add(1)
	worker.flush <- flush
	flush.Wait()

	return worker
}

// Run processes all incoming log messages from loggers. It emits received log
// records to all added log handlers for specific logger
func (worker *Worker) run() {
	for {
		select {
		case flush := <-worker.flush:
			for records := len(worker.records); records > 0; records-- {
				record := <-worker.records

				if record != nil {
					record.logger.emit(record)
				}
			}

			if flush != nil {
				flush.Done()
			}
		case record := <-worker.records:
			if record != nil {
				record.logger.emit(record)
			}
		}
	}
}
