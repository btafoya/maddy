# Maddy API Documentation

## Framework Packages

### framework/module

Core module system interfaces and registry.

#### Key Interfaces

```go
// Base module interface
type Module interface {
    Init(*config.Map) error
    Name() string
    InstanceName() string
}

// Delivery target
type DeliveryTarget interface {
    Module
    Start(ctx context.Context, msgMeta *MsgMetadata, mailFrom string) (Delivery, error)
}

// Message delivery
type Delivery interface {
    AddRcpt(ctx context.Context, rcptTo string, opts smtp.RcptOptions) error
    Body(ctx context.Context, header textproto.MIMEHeader, body buffer.Buffer) error
    Abort(ctx context.Context) error
    Commit(ctx context.Context) error
}

// Message check
type Check interface {
    Module
    CheckStateForMsg(ctx context.Context, msgMeta *MsgMetadata) (CheckState, error)
}

// Check state per message
type CheckState interface {
    CheckConnection(ctx context.Context) CheckResult
    CheckSender(ctx context.Context, mailFrom string) CheckResult
    CheckRcpt(ctx context.Context, rcptTo string) CheckResult
    CheckBody(ctx context.Context, header textproto.MIMEHeader, body buffer.Buffer) CheckResult
    Close() error
}

// Message modifier
type Modifier interface {
    Module
    ModStateForMsg(ctx context.Context, msgMeta *MsgMetadata) (ModifierState, error)
}

// Authentication provider
type PlainAuth interface {
    Module
    AuthPlain(username, password string) error
}

// Storage backend
type Storage interface {
    Module
    GetOrCreateIMAPAcct(username string) (Backend, error)
}
```

#### Registry Functions

```go
// Register a module factory
func Register(name string, factory FuncNewModule)

// Register an endpoint factory
func RegisterEndpoint(name string, factory FuncNewEndpoint)

// Get a module factory
func Get(name string) FuncNewModule

// Get an endpoint factory
func GetEndpoint(name string) FuncNewEndpoint

// Register a module instance
func RegisterInstance(inst Module, cfg *config.Map)

// Check if instance exists
func HasInstance(name string) bool

// Register an alias for an instance
func RegisterAlias(alias, original string)
```

### framework/config

Configuration parsing and directive handling.

#### config.Map

```go
// Create a new config map
func NewMap(globals map[string]interface{}, node Node) *Map

// Parse string directive
func (m *Map) String(name string, required, inherited bool, defaultVal string, dest *string)

// Parse bool directive
func (m *Map) Bool(name string, required, inherited bool, defaultVal *bool, dest *bool)

// Parse int directive
func (m *Map) Int(name string, required, inherited bool, defaultVal int, dest *int)

// Parse duration directive
func (m *Map) Duration(name string, required, inherited bool, defaultVal time.Duration, dest *time.Duration)

// Parse string list directive
func (m *Map) StringList(name string, required, inherited bool, defaultVal []string, dest *[]string)

// Custom directive handler
func (m *Map) Custom(name string, required, inherited bool, defaultVal interface{}, handler func(*Map, *Node) (interface{}, error), dest interface{})

// Allow unknown directives
func (m *Map) AllowUnknown()

// Process all directives
func (m *Map) Process() ([]Node, error)
```

#### config.Node

```go
type Node struct {
    Name     string   // Directive name
    Args     []string // Arguments
    Children []Node   // Child directives
    File     string   // Source file
    Line     int      // Source line
}

// Create error with location
func NodeErr(node Node, format string, args ...interface{}) error
```

### framework/exterrors

Extended error handling with SMTP codes.

#### SMTPError

```go
type SMTPError struct {
    Code       int                    // SMTP code (e.g., 550)
    EnhCode    string                 // Enhanced code (e.g., "5.7.1")
    Message    string                 // Client-visible message
    Err        error                  // Wrapped error
    Reason     string                 // Override Err.Error()
    CheckName  string                 // Check module name
    TargetName string                 // Target module name
    Misc       map[string]interface{} // Additional context
}
```

#### Helper Functions

```go
// Wrap error with fields
func WithFields(err error, fields map[string]interface{}) error

// Mark error as temporary
func WithTemporary(err error, temp bool) error

// Check if error is temporary
func IsTemporary(err error) bool

// Select SMTP code based on temporary status
func SMTPCode(err error, tempCode, permCode int) int
```

### framework/hooks

Lifecycle event management.

#### Events

```go
const (
    EventShutdown  = "shutdown"  // Server shutdown
    EventLogRotate = "logrotate" // Log rotation signal
)
```

#### Functions

```go
// Add a hook for an event
func AddHook(event string, fn func())

// Run all hooks for an event
func RunHooks(event string)
```

### framework/log

Logging utilities.

```go
var DefaultLogger *Logger

type Logger struct {
    Out   io.WriteCloser
    Debug bool
    Name  string
}

func (l *Logger) Printf(format string, args ...interface{})
func (l *Logger) Println(args ...interface{})
func (l *Logger) Debugf(format string, args ...interface{})
func (l *Logger) Debugln(args ...interface{})
```

### framework/buffer

Buffered I/O for message bodies.

```go
type Buffer interface {
    Open() (io.ReadCloser, error)
    Len() int
}

// Create memory buffer
func MemoryBuffer(data []byte) Buffer

// Create file buffer
func FileBuffer(path string) Buffer
```

### framework/dns

DNS utilities.

```go
// Lookup MX records
func LookupMX(ctx context.Context, domain string) ([]*net.MX, error)

// Lookup TXT records
func LookupTXT(ctx context.Context, domain string) ([]string, error)

// Lookup A/AAAA records
func LookupIP(ctx context.Context, host string) ([]net.IP, error)
```

### framework/address

Email address parsing.

```go
// Parse email address
func Parse(addr string) (mailbox, domain string, err error)

// Normalize email address
func Normalize(addr string) (string, error)

// Split into local and domain parts
func Split(addr string) (local, domain string)

// Check if address is valid
func Valid(addr string) bool
```

## Internal Packages

### internal/msgpipeline

Message processing pipeline.

```go
type MsgPipeline struct {
    // Configuration
    Checks    []module.Check
    Modifiers []module.Modifier
    Targets   []module.DeliveryTarget
}

func (p *MsgPipeline) Start(ctx context.Context, msgMeta *module.MsgMetadata, mailFrom string) (module.Delivery, error)
```

### internal/limits

Rate and concurrency limiting.

```go
type Group struct {
    // Limiters
    Rate        *limiters.Rate
    Concurrency *limiters.Concurrency
    Bucket      *limiters.Bucket
}

func (g *Group) Take(key string) error
func (g *Group) Release(key string)
```

### internal/dmarc

DMARC policy handling.

```go
type Verifier struct {
    // ...
}

func (v *Verifier) Verify(ctx context.Context, header textproto.MIMEHeader) (Result, error)

type Result struct {
    AuthResult   string // "pass", "fail", "none"
    Policy       string // "none", "quarantine", "reject"
    Disposition  string // Applied action
}
```

### internal/authz

Authorization utilities.

```go
// Normalize username for authentication
func NormalizeUsername(u string, method int) (string, error)

// Normalization methods
const (
    NormalizeAuto     = iota
    NormalizePrecis
    NormalizeCasefold
)
```

## CLI Interface

### maddycli Package

Located in `internal/cli/`.

```go
// Add global flag
func AddGlobalFlag(flag cli.Flag)

// Add subcommand
func AddSubcommand(cmd *cli.Command)

// Run CLI application
func Run()
```

### Built-in Commands

```
maddy run [options]          # Start server
maddy version                # Show version info

maddyctl users list          # List users
maddyctl users create <user> # Create user
maddyctl users delete <user> # Delete user

maddyctl imap-acct list      # List IMAP accounts
maddyctl imap-acct create    # Create IMAP account

maddyctl hash bcrypt <pass>  # Hash password
```

## Message Metadata

```go
type MsgMetadata struct {
    ID        string            // Unique message ID
    Conn      *smtp.ConnectionState
    OriginalFrom string         // Original MAIL FROM
    AuthUser  string            // Authenticated username
    Quarantine bool             // Message quarantined
    SrcHostname string          // Source hostname
    SrcProto  string            // Source protocol
    DontTraceSender bool        // Don't add sender trace
    SMTPOpts  smtp.MailOptions  // SMTP MAIL options
}
```

## Check Results

```go
type CheckResult struct {
    Reject bool          // Reject message
    Quarantine bool      // Quarantine message
    AuthResult []string  // Authentication-Results values
    Header textproto.MIMEHeader // Headers to add
    Reason error         // Rejection reason
}
```

## Configuration Patterns

### Module Reference

```go
// Reference existing module in config
modconfig.ModuleFromNode(name, globals, node, interfaces...)

// Example usage in Init():
var target module.DeliveryTarget
cfg.Custom("deliver_to", true, false, nil, modconfig.DeliveryTarget(&target))
```

### Table Reference

```go
// Get table from config
var table module.Table
modconfig.Table(cfg, "table", required, inherited, defaultVal, &table)
```

### TLS Configuration

```go
// Get TLS config from node
tlsCfg, err := tls.TLSDirective(cfg, node)
```
