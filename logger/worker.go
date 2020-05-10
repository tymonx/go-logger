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
	"os/signal"
	"sync"
	"syscall"
)

// These constants define default values for Worker
const (
	DefaultWorkerQueueLength = 4096
)

// A Worker represents an active logger worker thread. It handles formatting
// received log messages and I/O operations
type Worker struct {
	records   chan *Record
	signals   chan os.Signal
	waitGroup sync.WaitGroup
}

var gWorkerOnce sync.Once
var gWorkerInstance *Worker

// NewWorker creates a new Worker object
func NewWorker() *Worker {
	worker := &Worker{
		records: make(chan *Record, DefaultWorkerQueueLength),
		signals: make(chan os.Signal, 1),
	}

	worker.waitGroup.Add(1)

	signal.Notify(
		worker.signals,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

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
	gMutex.Lock()
	defer gMutex.Unlock()

	if length <= 0 {
		length = DefaultWorkerQueueLength
	}

	if cap(worker.records) != length {
		worker.records = make(chan *Record, length)
	}

	return worker
}

// SetSignals sets signals to stop logger worker thread
func (worker *Worker) SetSignals(signals ...os.Signal) *Worker {
	gMutex.Lock()
	defer gMutex.Unlock()

	signal.Stop(worker.signals)

	if len(signals) > 0 {
		signal.Notify(worker.signals, signals...)
	}

	return worker
}

// Close stops logger worker thread
func (worker *Worker) Close() {
	worker.signals <- os.Interrupt
	worker.waitGroup.Wait()
	signal.Stop(worker.signals)
}

// Run processes all incoming log messages from loggers. It emits received log
// records to all added log handlers for specific logger
func (worker *Worker) run() {
	var record *Record

	running := true

	for running {
		select {
		case <-worker.signals:
			running = false
		case record, running = <-worker.records:
			if record != nil {
				record.logger.emit(record)
			}
		}
	}

	for records := len(worker.records); records > 0; records-- {
		record = <-worker.records

		if record != nil {
			record.logger.emit(record)
		}
	}

	worker.waitGroup.Done()
}
