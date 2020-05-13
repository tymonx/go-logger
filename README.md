# Go Logger

Lightweight, fast and powerful logger implementation in Go.

## Features

* All log formatting and I/O operations are offloaded to separate worker thread
* All calls to log methods are lightweight and consumes very little CPU resources
* It can simultaneously log message to different log handlers
* Various customizable built-in log handlers `Stdout`, `Stderr`, `File`, `Stream` and `Syslog`
* Various log methods `Trace`, `Debug`, `Info`, `Notice`, `Warning`, `Error`, `Critical`, `Alert`, `Fatal` and `Panic`
* Flexible log message formatter with some predefined named placeholders
* Use new created logger instance or use a global one
* Supporting automatic placeholders in log arguments with `{p}`
* Supporting positional placeholders in log arguments with `{pN}`
* Supporting named placeholders in log arguments with `{name}`, `{p.name}` or `{pN.name}`
* Supporting object placeholders in log arguments with `{.field}`, `{p.field}` or `{pN.field}`
* Supporting custom placeholder identification (default is `p`)
* Supporting custom log handlers
* Supporting custom log formatters
* Supporting custom log date formats
* Supporting custom log message formats
* Supporting custom log ID generators
* Supporting exporting log records to JSON output

## Install

```
go get -u gitlab.com/tymonx/go-logger
```

## Example

```go
package main

import (
	"gitlab.com/tymonx/go-logger/logger"
)

func main() {
	// The close method is needed because all log methods are offloaded to
	// separate worker thread. The Close() function guarantees that all log
	// messages will be flushed out and all log handlers will be properly closed
	defer logger.Close()

	logger.Info("Hello from logger!")
	logger.Info("Automatic placeholders {p} {p} {p}", 1, 2, 3)
	logger.Info("Positional placeholders {p2} {p1} {p0}", 1, 2, 3)

	logger.Info("Named placeholders {z} {y} {x}", logger.Named{
		"x": 1,
		"y": 2,
		"z": 3,
	})

	logger.Info("Object placeholders {.Z} {.Y} {.X}", struct {
		X, Y, Z int
	}{
		X: 1,
		Y: 2,
		Z: 3,
	})
}
```

## Documentation

Go logger [documentation](https://tymonx.gitlab.io/go-logger/doc/pkg/gitlab.com/tymonx/go-logger/logger/).
