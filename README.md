# HORUS

[![License](https://img.shields.io/badge/license-GPL-blue.svg)](LICENSE)
[![Documentation](https://godoc.org/github.com/DanielRivasMD/horus?status.svg)](http://godoc.org/github.com/DanielRivasMD/horus)
[![Go Report Card](https://goreportcard.com/badge/github.com/DanielRivasMD/horus)](https://goreportcard.com/report/github.com/DanielRivasMD/horus)
[![Release](https://img.shields.io/github/release/DanielRivasMD/horus.svg?label=Release)](https://github.com/DanielRivasMD/horus/releases)
[![Coverage Status](https://coveralls.io/repos/github/DanielRivasMD/horus/badge.svg?branch=main)](https://coveralls.io/github/DanielRivasMD/horus?branch=main)

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

```go
err := doSomething()
// only if err != nil it gets wrapped
return horus.PropagateErr("DoSomething", "SERVICE", "failed", err, nil)
```

### Flexible Formatting

- `JSONFormatter` for structured logs
- `PseudoJSONFormatter` for aligned, colorized tables in your terminal
- `PlainFormatter` or `SimpleColoredFormatter` for minimal output

### Check & Exit

- `CheckErr(err, opts...)` writes formatted error to your choice of `io.Writer` and exits with a customizable code
- Built-in overrides: writer, formatter, exit code, operation, category, message, details

```go
horus.CheckErr(err)               // default: colored table + os.Exit(1)
horus.CheckErr(err, horus.WithWriter(os.Stdout), horus.WithExitCode(42))
```

### Not-Found Hooks

- `LogNotFound` / `NullAction` implement `NotFoundAction` for pluggable “resource missing” behaviors
- Fully testable via `WithLogWriter`

```go
act := horus.LogNotFound("cache miss")
resolved, err := act("user:123")
// resolved==false, err==nil
```

### Test Utilities

- `CollectingError` (implements `io.Writer` + `error`) to capture and inspect output in tests
- Easy use of `WithWriter(buf)`  to drive deterministic output

### Panic Integration

- `Panic(op, msg)` logs a colored panic banner, captures a stack, then panics with a full `Herror` payload


## Quickstart

```go
package main

import (
  "errors"
  "fmt"

  "github.com/DanielRivasMD/horus"
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

### Error Handling Integration

The horus error-handling library provides a set of powerful functions to wrap, propagate, log, and format errors across your application. Here’s how you can leverage these functions at various layers:

#### 1. Lower-Level Functions

Wrap errors as soon as they occur

For example, when reading a configuration file:

```go
package fileutils

import (
  "fmt"
  "os"

  "github.com/DanielRivasMD/horus"
)

// ReadConfig tries to read a JSON config from disk.
// On failure it wraps the underlying error with full context, category and details.
func ReadConfig(path string) ([]byte, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return nil, horus.PropagateErr(
      "ReadConfig",                              // Op
      "IO_ERROR",                                // Category
      fmt.Sprintf("unable to load config"),      // Message
      err,                                       // underlying error
      map[string]any{                            // Details
        "path": path,
      },
    )
  }
  return data, nil
}
```

What this does:

- Uses Go 1.16+’s `os.ReadFile` instead of the deprecated `ioutil`
- Always returns `nil` or a rich `*Herror`, never a raw `error`
- Stamps on:
  - `Op` = "ReadConfig"
  - `Category` = "IO_ERROR"
  - `Message` = a user-friendly "unable to load config"
  - `Details` = {"path": path}
  - `Stack` trace (captured at the call site)

You’ll now get output like:

```text
Op       ReadConfig,
Message  unable to load config,
Err      open /etc/app.cfg: no such file or directory,
path     /etc/app.cfg,
Category IO_ERROR,

Stack
  fileutils.ReadConfig()
    /Users/.../fileutils/config.go:12
  main.main()
    /Users/.../cmd/app/main.go:23
  runtime.main()
    /usr/local/go/src/runtime/proc.go:250
  ...
```

and if you prefer JSON:

```go
horus.CheckErr(err, horus.WithFormatter(horus.JSONFormatter))
```

will emit something like:

```json
{
  "Op": "ReadConfig",
  "Message": "unable to load config",
  "Err": "open /etc/app.cfg: no such file or directory",
  "Details": { "path": "/etc/app.cfg" },
  "Category": "IO_ERROR",
  "Stack": [ ... ]
}
```

#### 2. Business Logic Functions

When higher-level functions catch errors from lower‐level routines, add domain‐specific context with `PropagateErr` (or `WithDetail`) so every layer contributes its own clues:


```go
package business

import (
  "github.com/your_module/fileutils"
  "github.com/DanielRivasMD/horus"
)

// LoadAndProcessConfig orchestrates reading + validating your config.
// Any I/O or parse failures get wrapped with step-specific context.
func LoadAndProcessConfig(configPath string) error {
  // 1. Call the fileutils helper
  data, err := fileutils.ReadConfig(configPath)
  if err != nil {
    // PropagateErr merges the underlying category/details and stamps on new ones.
    return horus.PropagateErr(
      "LoadAndProcessConfig",           // operation name
      "CONFIG_ERROR",                   // business category
      "unable to load application config", // user-friendly message
      err,                              // the error from ReadConfig
      map[string]any{                   // extra details
        "path":   configPath,
        "service": "business",
      },
    )
  }

  // 2. (Optional) add validation context
  if len(data) == 0 {
    // NewHerror creates a fresh Herror; WithDetail would wrap an existing one.
    return horus.NewHerror(
      "LoadAndProcessConfig",
      "config data is empty",
      nil,
      map[string]any{"path": configPath},
    )
  }

  // 3. process the config...
  return nil
}
```

- `PropagateErr` automatically carries forward any category or details from the lower layer (and merges the new payload)
- We choose a “business” category ("CONFIG_ERROR") that’s orthogonal to the lower-level "IO_ERROR"
- We include both `path` and our own service:"business" detail.
- For purely business‐rule failures (like empty data), we use `NewHerror` to start a fresh error.

#### 3. Centralized Error Reporting and Logging

At the top‐level ofthe app - usually in `main()` - `CheckErr` can be use as one‐stop fatal error handler:

- Format the error (colored table by default) or JSON if you prefer
- Register the error’s category in a global registry (for metrics/observability)
- Exit with a configurable code (default: 1)


```go
package main

import (
  "os"

  "github.com/DanielRivasMD/horus"
  "github.com/your_module/business"
)

func main() {
  // Run your business logic
  err := business.LoadAndProcessConfig("config.json")
  if err != nil {
    // Default: colored table → stderr, exit code 1
    horus.CheckErr(err)

    // Or JSON + code 2 + log to stdout:
    // horus.CheckErr(
    //   err,
    //   horus.WithFormatter(horus.JSONFormatter),
    //   horus.WithExitCode(2),
    //   horus.WithWriter(os.Stdout),
    // )
  }

  // Continue with normal execution...
}
```

Under the hood, `CheckErr` does:

- `RegisterError(err)` – increments a counter for your error’s category
- `fmt.Fprintln(writer, formatter(err))` – prints your chosen format
- `exitFunc(code)` – calls `os.Exit(code)` by default

This ensures that all unhandled, fatal errors flow through a consistent, observable pipeline

#### 4. Integration with External Processes

For commands executed via external processes (e.g., running system commands), use functions like `ExecCmd` or `CaptureExecCmd`, `horus`:

- Shows both `ExecCmd` (streams output) and `CaptureExecCmd` (buffers it)
- Wraps errors with `PropagateErr` for context before calling `CheckErr`
- Logs stdout/stderr details in your error’s `Details` map

```go
package domovoi

import (
  "fmt"

  "github.com/DanielRivasMD/horus"
)

// ListDirectory runs `ls -la <path>` twice: once streamed, once captured.
// It returns an error if anything fails; caller can then CheckErr(err).
func ListDirectory(path string) error {
  // 1) Stream mode
  if err := ExecCmd("ls", "-la", path); err != nil {
    // ExecCmd already wraps in *Herror, but we can add our own op/category
    return horus.NewCategorizedHerror(
      "ListDirectory",                // Op
      "SYS_CMD",                      // Category
      fmt.Sprintf("ls -la %s", path), // Message
      err,
      nil, // no extra details here
    )
  }

  // 2) Capture mode
  stdout, stderr, err := CaptureExecCmd("ls", "-la", path)
  if err != nil {
    // We got a *Herror from CaptureExecCmd—merge in the captured output
    return horus.PropagateErr(
      "ListDirectory",
      "SYS_CMD",
      "failed to capture ls output",
      err,
      map[string]any{"stdout": stdout, "stderr": stderr},
    )
  }

  fmt.Println("=== STDOUT ===")
  fmt.Print(stdout)
  fmt.Println("=== STDERR ===")
  fmt.Print(stderr)
  return nil
}
```

- `ExecCmd(op, category, message, *exec.Cmd)` runs and logs failures immediately
- `CaptureExecCmd` returns `(stdout, stderr string, err error)`
- We propagate errors with `PropagateErr` so our top‐level `CheckErr` shows the full story: operation, category, message, underlying cause, stdout, stderr, and stack trace

#### 5. OS Level Operations

When wrapping system calls like `os.Chdir`, use Horus to enrich and propagate errors with full context:

```go
package domovoi

import (
  "fmt"
  "os"

  "github.com/DanielRivasMD/horus"
)

// ChangeDirectory attempts to chdir into the given path.
// On failure it wraps the underlying os.Chdir error with operation,
// category, user-friendly message, and the path in Details.
func ChangeDirectory(path string) error {
  if err := os.Chdir(path); err != nil {
    return horus.PropagateErr(
      "ChangeDirectory",                               // Op
      "FS_ERROR",                                      // Category
      fmt.Sprintf("unable to change working directory"),// Message
      err,                                             // underlying error
      map[string]any{"path": path},                    // Details
    )
  }
  return nil
}
```

- `Op` = "ChangeDirectory"
- `Category` = "FS_ERROR"
- `Message` = "unable to change working directory"
- `Details` = {"path": "/some/dir"}
- `Stack` trace captured at the call site

That way, any failure bubbles up as a full `*Herror` - complete with stack, category, and details—making your logs and CLI output immediately actionable

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

## License
Copyright (c) 2025

See the [LICENSE](LICENSE) file for license details.
