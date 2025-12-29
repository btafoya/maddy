# Maddy Development Guide

## Quick Start

```bash
# Clone repository
git clone https://github.com/foxcpp/maddy.git
cd maddy

# Build
./build.sh

# Run tests
go test ./...

# Run linter
golangci-lint run

# Run server
./build/maddy run --config maddy.conf
```

## Development Environment Setup

### Prerequisites

- Go 1.18+ (project uses 1.23.1)
- C compiler (for SQLite support)
- `golangci-lint` for linting
- Optional: `scdoc` for man pages
- Optional: `libpam-dev` for PAM support

### IDE Configuration

The project includes `.editorconfig`:
- Tab indentation for Go files
- LF line endings
- Trim trailing whitespace

## Project Layout

```
/                     # Root package - server initialization
├── cmd/maddy/       # Main binary entry point
├── framework/       # Reusable, stable API packages
└── internal/        # Implementation packages (unstable API)
```

### Framework vs Internal

**framework/** packages:
- Stable API
- Designed for reuse
- Module interfaces
- Configuration parsing
- Error handling

**internal/** packages:
- Implementation details
- Unstable API
- Module implementations
- Not for external use

## Adding New Features

### Creating a New Check

1. Create directory: `internal/check/mycheck/`
2. Copy skeleton: `internal/check/skeleton.go`
3. Implement interfaces
4. Register in `maddy.go` imports

```go
// internal/check/mycheck/mycheck.go
package mycheck

import (
    "context"
    "net/textproto"

    "github.com/foxcpp/maddy/framework/buffer"
    "github.com/foxcpp/maddy/framework/config"
    "github.com/foxcpp/maddy/framework/module"
)

const modName = "check.mycheck"

type Check struct {
    modName  string
    instName string
    // config fields
}

func init() {
    module.Register(modName, New)
}

func New(_, instName string, _, _ []string) (module.Module, error) {
    return &Check{
        modName:  modName,
        instName: instName,
    }, nil
}

func (c *Check) Name() string         { return c.modName }
func (c *Check) InstanceName() string { return c.instName }

func (c *Check) Init(cfg *config.Map) error {
    // Parse configuration
    return nil
}

// Per-message state
type checkState struct {
    c *Check
}

func (c *Check) CheckStateForMsg(ctx context.Context, msgMeta *module.MsgMetadata) (module.CheckState, error) {
    return &checkState{c: c}, nil
}

func (s *checkState) CheckConnection(ctx context.Context) module.CheckResult {
    return module.CheckResult{}
}

func (s *checkState) CheckSender(ctx context.Context, mailFrom string) module.CheckResult {
    return module.CheckResult{}
}

func (s *checkState) CheckRcpt(ctx context.Context, rcptTo string) module.CheckResult {
    return module.CheckResult{}
}

func (s *checkState) CheckBody(ctx context.Context, header textproto.MIMEHeader, body buffer.Buffer) module.CheckResult {
    // Implement check logic
    return module.CheckResult{}
}

func (s *checkState) Close() error {
    return nil
}
```

5. Add import to `maddy.go`:
```go
_ "github.com/foxcpp/maddy/internal/check/mycheck"
```

### Creating a New Modifier

Similar pattern to checks, but implement `module.Modifier`:

```go
type Modifier interface {
    Module
    ModStateForMsg(ctx context.Context, msgMeta *MsgMetadata) (ModifierState, error)
}

type ModifierState interface {
    RewriteSender(ctx context.Context, from string) (string, error)
    RewriteRcpt(ctx context.Context, to string) (string, error)
    RewriteBody(ctx context.Context, header *textproto.Header, body buffer.Buffer) error
    Close() error
}
```

### Creating a New Target

Implement `module.DeliveryTarget`:

```go
type DeliveryTarget interface {
    Module
    Start(ctx context.Context, msgMeta *MsgMetadata, mailFrom string) (Delivery, error)
}

type Delivery interface {
    AddRcpt(ctx context.Context, rcptTo string, opts smtp.RcptOptions) error
    Body(ctx context.Context, header textproto.MIMEHeader, body buffer.Buffer) error
    Abort(ctx context.Context) error
    Commit(ctx context.Context) error
}
```

## Error Handling

### Use exterrors Package

```go
import "github.com/foxcpp/maddy/framework/exterrors"

// Good: SMTP-aware error
return exterrors.SMTPError{
    Code:      550,
    EnhCode:   exterrors.EnhancedCode{5, 7, 1},
    Message:   "Access denied",
    Err:       originalErr,
    CheckName: "check.mycheck",
    Misc: map[string]interface{}{
        "detail": "additional context",
    },
}

// Good: Add context to existing error
return exterrors.WithFields(err, map[string]interface{}{
    "module": "mycheck",
    "key":    value,
})

// Good: Mark as temporary
return exterrors.WithTemporary(err, true)
```

### SMTP Codes

| Code | Type | Usage |
|------|------|-------|
| 2xx | Success | Delivery accepted |
| 4xx | Temporary | Retry later |
| 5xx | Permanent | Don't retry |

## Concurrency Guidelines

### Per-Message State

Checks and modifiers create per-message state objects:
- State objects are NOT shared between messages
- Multiple messages can be processed concurrently
- Use `CheckStateForMsg` / `ModStateForMsg` for state

### Goroutine Safety

```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("panic recovered: %v", r)
        }
    }()
    // work here
}()
```

### Module Cleanup

Implement `io.Closer` for cleanup:

```go
func (m *MyModule) Close() error {
    // Stop goroutines
    // Close connections
    // Release resources
    return nil
}
```

## Testing

### Unit Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/check/mycheck/

# With coverage
go test -cover ./internal/check/mycheck/

# Verbose
go test -v ./internal/check/mycheck/
```

### Integration Tests

```bash
cd tests/
./run.sh
```

### Test Utilities

Located in `internal/testutils/`:

```go
import "github.com/foxcpp/maddy/internal/testutils"

// Create test message
msg := testutils.NewTestMessage(t)

// Mock DNS
resolver := testutils.MockResolver{
    MX: map[string][]*net.MX{
        "example.com": {{Host: "mx.example.com", Pref: 10}},
    },
}
```

## Configuration Testing

### Test Configuration Parsing

```go
func TestMyModuleConfig(t *testing.T) {
    mod, err := New("check.mycheck", "test", nil, nil)
    if err != nil {
        t.Fatal(err)
    }

    cfg := config.NewMap(nil, config.Node{
        Children: []config.Node{
            {Name: "option", Args: []string{"value"}},
        },
    })

    if err := mod.Init(cfg); err != nil {
        t.Fatal(err)
    }
}
```

## Debugging

### Enable Debug Logging

```bash
maddy --debug run --config maddy.conf
```

Or in config:
```
debug on
```

### pprof Profiling

Build with debug tags:
```go
// maddy_debug.go (already in repo)
// +build debugflags
```

Run with:
```bash
go build -tags debugflags ./cmd/maddy
./maddy --debug.pprof 127.0.0.1:6060 run --config maddy.conf
```

Access at: http://127.0.0.1:6060/debug/pprof/

## Code Style

### Imports

Use goimports for automatic formatting:
```bash
goimports -w .
```

Order:
1. Standard library
2. Third-party
3. Local (github.com/foxcpp/maddy/...)

### Naming

- Use clear, descriptive names
- Follow Go conventions (camelCase)
- Use `modName`/`instName` pattern for modules

### Comments

- Document exported functions
- Explain non-obvious logic
- Reference relevant RFCs

## Linting

```bash
# Run all linters
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

Enabled linters (from `.golangci.yml`):
- gosimple
- errcheck
- staticcheck
- ineffassign
- typecheck
- govet
- unused
- goimports
- prealloc
- unconvert
- misspell
- whitespace
- nakedret
- dogsled
- copyloopvar

## Documentation

### User Documentation

Located in `docs/` (MkDocs format):
- `docs/reference/` - Configuration reference
- `docs/tutorials/` - Setup guides
- `docs/man/` - Man pages (scdoc format)

### Building Man Pages

Requires `scdoc`:
```bash
scdoc < docs/man/maddy.1.scd > maddy.1
```

### Building Docs Site

```bash
mkdocs serve  # Local preview
mkdocs build  # Build static site
```

## Contributing Workflow

1. Fork repository
2. Create feature branch
3. Make changes
4. Run tests: `go test ./...`
5. Run linter: `golangci-lint run`
6. Commit (no AI attribution)
7. Create pull request

### Commit Messages

```
type(scope): description

[optional body]
```

Types: feat, fix, docs, refactor, test, chore

Example:
```
feat(check): add custom header injection

Allows checks to add custom headers to messages
for downstream processing.
```
