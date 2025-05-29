# HORUS

[![License](https://img.shields.io/badge/license-GPL-blue.svg)](LICENSE)
[![go Version](https://img.shields.io/badge/go-VERSION-green.svg)](LURL)

## Overview

horus: a Go library for error handling and propagation with rich context and stack traces.

---

## Features
- **Context-Rich Errors**: Automatically wrap and enrich errors with context and stack traces.
- **Error Propagation**: Seamlessly pass errors upward with additional details.
- **Custom Formatting**: Supports colored output and JSON representations for logging.
- **Integration Ready**: Easily integrate with external libraries to build robust applications.

---
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
