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

const (
	// DefaultWorkerQueueLength defines default queue length for log messages
	// between logger and logger worker thread
	DefaultWorkerQueueLength = 4096
)

// A worker represents an active logger worker thread. It handles formatting
// received log messages and I/O operations
type worker struct {
	records       chan *Record
	signals       chan os.Signal
	waitGroup     sync.WaitGroup
	isQueueClosed bool
}

var gWorkerOnce sync.Once
var gWorkerInstance *worker
var gWorkerQueueLength = DefaultWorkerQueueLength

// gWorkerSignals default signals used to stop logger worker thread
var gWorkerSignals = []os.Signal{
	os.Kill,
	os.Interrupt,
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
}

// SetWorkerQueueLength sets logger worker thread queue length for log messages.
// This function must be called before using logger
func SetWorkerQueueLength(length int) {
	gMutex.Lock()
	defer gMutex.Unlock()

	if length > 0 {
		gWorkerQueueLength = length
	} else {
		gWorkerQueueLength = DefaultWorkerQueueLength
	}
}

// SetWorkerSignals sets signals to stop logger worker thread. This function
// must be called before using logger
func SetWorkerSignals(signals ...os.Signal) {
	gMutex.Lock()
	defer gMutex.Unlock()

	gWorkerSignals = signals
}


// getWorker returns global logger worker instance. First call creates and
// starts logger worker thread
func getWorker() *worker {
	gWorkerOnce.Do(func() {
		gWorkerInstance = &worker{
			records: make(chan *Record, gWorkerQueueLength),
			signals: make(chan os.Signal, 1),
		}

		gWorkerInstance.waitGroup.Add(1)

		if len(gWorkerSignals) > 0 {
			signal.Notify(gWorkerInstance.signals, gWorkerSignals...)
		}

		go gWorkerInstance.run()
	})

	return gWorkerInstance
}

// run processes all incoming log messages from loggers. It emits received log
// records to all added log handlers for specific logger
func (worker *worker) run() {
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

// closeQueue closes log message queue
func (worker *worker) closeQueue() {
	gMutex.Lock()
	defer gMutex.Unlock()

	if !worker.isQueueClosed {
		close(worker.records)
		worker.isQueueClosed = true
	}
}

// close stops logger worker thread
func (worker *worker) close() {
	worker.closeQueue()
	worker.waitGroup.Wait()
}
