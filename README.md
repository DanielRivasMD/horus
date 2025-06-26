# HORUS

[![License](https://img.shields.io/badge/license-GPL-blue.svg)](LICENSE)
[![go Version](https://img.shields.io/badge/go-VERSION-green.svg)](LURL)

## Overview
`horus` is a Go error‐handling toolkit that does more than “return an error”

`horus` captures the _what_ (operation), _why_ (user‐friendly message), _where_ (stack trace), and _how_ (underlying cause and arbitrary key/value details), then lets you:

- Wrap and re‐wrap errors with layered context
- Seamlessly propagate or bail out with `CheckErr` (colored/JSON + configurable exit codes)
- Plug into any CLI, HTTP handler, or logger with zero ceremony
- Whether you’re building a command‐line tool, a microservice, or a big data pipeline, `horus` ensures nothing gets lost in translation when things go sideways

## Features

### Context-Rich Errors

Create `Herror` instances via `NewHerror`, `NewCategorizedHerror`, `Wrap` or `WithDetail` that carry:

- Operation name (`Op`)
- Human-readable message (`Message`)
- Category tag (`Category`)
- Arbitrary details map (`Details`)
- Full stack trace

### Error Propagation & Inspection

- `PropagateErr` for idiomatic upstream wrapping
- `RootCause(err)` to peel back nested failures
- Helpers like `AsHerror`, `IsHerror`, `Operation`, `UserMessage`, `GetDetail`, `Category`, `StackTrace`

### Flexible Formatting

- `JSONFormatter` for structured logs
- `PseudoJSONFormatter` for aligned, colorized tables in your terminal
- `PlainFormatter` or `SimpleColoredFormatter` for minimal output

### Check & Exit

- `CheckErr(err, opts...)` writes formatted error to your choice of `io.Writer` and exits with a customizable code
- Built-in overrides: writer, formatter, exit code, operation, category, message, details

### Not-Found Hooks

- `LogNotFound` / `NullAction` implement `NotFoundAction` for pluggable “resource missing” behaviors
- Fully testable via `WithLogWriter`

### Test Utilities

- `CollectingError` (implements `io.Writer` + `error`) to capture and inspect output in tests
- Easy use of `WithWriter(buf)`  to drive deterministic output

### Panic Integration

- `Panic(op, msg)` logs a colored panic banner, captures a stack, then panics with a full `Herror` payload


## Installation

### **Language-Specific**
| Language | Command                                                            |
|----------|--------------------------------------------------------------------|
| **Go**   | `go get github.com/DanielRivasMD/horus@latest`                     |

---

## Usage
```go
import "github.com/DanielRivasMD/horus"
```

## Quickstart

```
package main

import (
  "errors"
  "fmt"

  "github.com/your/module/horus"
)

func loadConfig(path string) error {
  // pretend this fails
  return errors.New("file not found")
}

func main() {
  err := loadConfig("/etc/app.cfg")
  if err != nil {
    // wrap with context, category, and detail
    wrapped := horus.NewCategorizedHerror(
      "LoadConfig",
      "IO_ERROR",
      "unable to load configuration",
      err,
      map[string]any{"path": "/etc/app.cfg"},
    )
    // print a pretty, colored table to stderr and exit
    horus.CheckErr(wrapped)
  }

  fmt.Println("config loaded")
}
```

## Example




### Error Handling Integration

The horus error-handling library provides a set of powerful functions to wrap, propagate, log, and format errors across your application. Here’s how you can leverage these functions at various layers:

#### 1. Lower-Level Functions

Wrap errors as soon as they occur. For example, when reading a configuration file:

```go
package fileutils

import (
    "io/ioutil"
    "github.com/DanielRivasMD/horus"
)

func ReadConfig(path string) ([]byte, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        // Wrap the error immediately with context, including a stack trace.
        return nil, horus.Wrap(err, "read config", "failed to read config file")
    }
    return data, nil
}
```

This approach ensures that the error is enriched with context right from the moment it occurs.

---

#### 2. Business Logic Functions

When higher-level functions catch errors from lower-level routines, use `PropagateErr` or `WithDetail` to add extra context:

```go
package business

import (
    "github.com/DanielRivasMD/horus"
    "github.com/yourrepo/fileutils"
)

func LoadAndProcessConfig(configPath string) error {
    data, err := fileutils.ReadConfig(configPath)
    if err != nil {
        // Propagate the error with extra details for context.
        return horus.PropagateErr(
            "load config",
            "config_error",
            "unable to load configuration",
            err,
            map[string]any{"configPath": configPath},
        )
    }
    // Further processing with data...
    _ = data
    return nil
}
```

Propagating the error in this manner ensures that every layer receives the necessary context for effective debugging.

---

#### 3. Centralized Error Reporting and Logging

At the entry point of your application—typically in your `main` function—you can use `CheckErr` to log errors in a colored, structured format and exit gracefully if necessary:

```go
package main

import (
    "github.com/DanielRivasMD/horus"
    "github.com/yourrepo/business"
)

func main() {
    err := business.LoadAndProcessConfig("config.json")
    if err != nil {
        // Logs the error with a colored, formatted output and registers error metrics.
        horus.CheckErr(err)
    }

    // Continue with the rest of your application...
}
```

`horus.CheckErr` acts as your centralized safety net for fatal errors.

---

#### 4. Integration with External Processes

For commands executed via external processes (e.g., running system commands), use functions like `ExecCmd` or `CaptureExecCmd`:

```go
package domovoi

import (
    "fmt"
    "github.com/DanielRivasMD/horus"
)

func SomeCommandRoutine() {
    stdout, stderr, err := horus.CaptureExecCmd("ls", "-la")
    if err != nil {
        // Handle critical error with CheckErr or log it appropriately.
        horus.CheckErr(err)
    }
    fmt.Println("Standard Output:", stdout)
    fmt.Println("Error Output:", stderr)
}
```

This function wraps and logs any error encountered during command execution, providing both stdout and stderr for complete diagnostics.

---

#### 5. OS Level Operations

For operations such as changing directories, use horus to capture detailed context if something goes wrong:

```go
package domovoi

import (
    "github.com/DanielRivasMD/horus"
)

func ChangeDirectory(path string) error {
    err := ChangeDir(path)
    if err != nil {
        // The error is wrapped with colored output and detailed context.
        return err
    }
    return nil
}
```

By utilizing functions like `horus.NewCategorizedHerror` within these operations, even system-level errors carry the necessary context for troubleshooting.

---

## Configuration

horus uses sensible defaults for error capturing and formatting but can be customized via environment variables or function parameters.

---
## Development

Build from source:
```bash
git clone https://github.com/DanielRivasMD/horus
cd horus
```

## Language-Specific Setup

| Language | Dev Dependencies | Hot Reload           |
|----------|------------------|----------------------|
| Go       | `go >= 1.22`     | `air` (live reload)  |

---
## FAQ
Q: How to resolve?
A: Use `horus.Wrap`, `horus.PropagateErr`, and `horus.CheckErr` for detailed context and troubleshooting.

Q: Cross-platform support?
A: Yes, horus is designed to work seamlessly on Windows, macOS, and Linux.

## License
GPL [2025] [Daniel Rivas]
