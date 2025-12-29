# Maddy Architecture Documentation

## System Overview

```
                    ┌─────────────────────────────────────────┐
                    │            Configuration                │
                    │           (maddy.conf)                  │
                    └─────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Module Registry                                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │  Endpoints   │  │   Targets    │  │    Checks    │  │  Modifiers   │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │    Auth      │  │   Storage    │  │    Tables    │  │     TLS      │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
              ┌───────────────────────┼───────────────────────┐
              ▼                       ▼                       ▼
     ┌─────────────┐         ┌─────────────┐         ┌─────────────┐
     │    SMTP     │         │    IMAP     │         │  Metrics    │
     │  Endpoint   │         │  Endpoint   │         │  Endpoint   │
     └─────────────┘         └─────────────┘         └─────────────┘
```

## Initialization Sequence

```
1. CLI Parsing
   cmd/maddy/main.go → maddycli.Run()
        │
        ▼
2. Configuration Loading
   maddy.go:Run() → parser.Read() → cfg []config.Node
        │
        ▼
3. Global Directives
   ReadGlobals() → state_dir, runtime_dir, hostname, tls, logging
        │
        ▼
4. Directory Setup
   InitDirs() → create state/runtime dirs, change working dir
        │
        ▼
5. Module Registration
   RegisterModules() → create instances, register in global map
        │
        ├─► Endpoints: created, added to endpoints list
        └─► Other modules: created, registered in instances map
        │
        ▼
6. Module Initialization
   initModules() → endpoints.Init(), lazy init for others
        │
        ▼
7. Ready State
   systemdStatus(SDReady) → signal systemd
        │
        ▼
8. Signal Handling
   handleSignals() → wait for SIGTERM/SIGINT
        │
        ▼
9. Graceful Shutdown
   hooks.RunHooks(EventShutdown) → close all modules
```

## Message Processing Pipeline

### Inbound Mail (MX)

```
Remote Server
      │
      ▼
┌─────────────────┐
│  SMTP Endpoint  │ ← Connection checks
│   (port 25)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  CheckConnection│ ← IP reputation, DNSBL
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   MAIL FROM     │
│  CheckSender    │ ← SPF, sender verification
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    RCPT TO      │
│   CheckRcpt     │ ← Recipient validation
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│      DATA       │
│   CheckBody     │ ← DKIM, content filters, rspamd
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Modifiers     │ ← Header manipulation
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Target(s)    │ ← Local storage, forwarding
└─────────────────┘
```

### Outbound Mail (Submission)

```
Mail Client
      │
      ▼
┌─────────────────┐
│  SMTP Endpoint  │ ← TLS required
│  (port 587/465) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Authentication  │ ← SASL (PLAIN, LOGIN)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Authorization  │ ← Sender == authenticated user
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Modifiers     │ ← DKIM signing
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  target.queue   │ ← Queued for delivery
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ target.remote   │ ← MX lookup, SMTP delivery
└─────────────────┘
```

### IMAP Access

```
Mail Client
      │
      ▼
┌─────────────────┐
│  IMAP Endpoint  │ ← TLS required
│  (port 993/143) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Authentication  │ ← SASL authentication
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  IMAP Filters   │ ← Server-side filtering
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Storage Backend │ ← imapsql (SQLite/PG/MySQL)
└─────────────────┘
```

## Module Relationships

### Dependency Graph

```
smtp endpoint
    ├── auth provider (for submission)
    ├── message pipeline
    │   ├── checks (dkim, spf, dnsbl, etc.)
    │   ├── modifiers (dkim signing)
    │   └── targets
    │       ├── local storage
    │       └── queue → remote
    └── tls config

imap endpoint
    ├── auth provider
    ├── storage backend
    └── tls config

storage.imapsql
    ├── blob storage (fs or s3)
    └── database (sqlite/pg/mysql)
```

### Module Interface Hierarchy

```
module.Module (base)
    │
    ├── module.DeliveryTarget
    │   └── target.queue, target.remote, target.smtp
    │
    ├── module.Check
    │   └── check.dkim, check.spf, check.dnsbl, etc.
    │
    ├── module.Modifier
    │   └── modify.dkim, modify.replace_addr
    │
    ├── module.Table
    │   └── Various lookup tables
    │
    ├── module.PlainAuth
    │   └── auth.pam, auth.shadow, auth.ldap, etc.
    │
    └── module.Storage
        └── storage.imapsql
```

## Configuration Model

### Block Types

```
# Global directive
directive value

# Module definition
module_type instance_name {
    directive value
    nested_module {
        directive value
    }
}

# Module reference
&instance_name

# Inline module
parent {
    child_directive module_type args {
        config
    }
}
```

### Resolution Flow

```
Configuration Text
       │
       ▼
┌──────────────┐
│  cfgparser   │ → []config.Node (AST)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  config.Map  │ → Type-safe directive parsing
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Module     │ → Initialized module instance
└──────────────┘
```

## Error Handling Model

### Error Flow

```
Low-level error (I/O, network, etc.)
       │
       ▼
┌──────────────────────────┐
│    exterrors.SMTPError   │
│  ├── Code (550, 421...)  │
│  ├── EnhCode (5.7.1...)  │
│  ├── Message             │
│  ├── Err (wrapped)       │
│  ├── CheckName           │
│  └── Misc (context)      │
└──────────────────────────┘
       │
       ▼
Module returns error
       │
       ▼
SMTP response to client
```

### Temporary vs Permanent

```go
// Temporary error (4xx) - retry later
exterrors.WithTemporary(err, true)

// Permanent error (5xx) - don't retry
exterrors.WithTemporary(err, false)

// Auto-detect from SMTP code
exterrors.SMTPCode(err, 421, 521) // 421 if temp, 521 if perm
```

## Concurrency Model

### Per-Connection

```
SMTP Connection
    │
    └─► Goroutine
        ├── Read commands
        ├── Process message
        │   ├── Checks (can be parallel)
        │   ├── Modifiers (sequential)
        │   └── Targets (can be parallel)
        └── Send responses
```

### Global Coordination

```
┌────────────────────────────────────────┐
│            Main Goroutine              │
│  ├── Signal handling                   │
│  └── Shutdown coordination             │
└────────────────────────────────────────┘
           │
           │ hooks.EventShutdown
           ▼
┌────────────────────────────────────────┐
│         Per-Module Cleanup             │
│  ├── Close listeners                   │
│  ├── Wait for in-flight requests       │
│  └── Release resources                 │
└────────────────────────────────────────┘
```

## State Management

### Directory Layout

```
/var/lib/maddy/ (state_dir)
    ├── imapsql/        # IMAP database
    │   ├── imapsql.db  # SQLite database
    │   └── blobs/      # Message blobs
    ├── queue/          # Message queue
    ├── dkim_keys/      # DKIM private keys
    └── acme/           # ACME certificates

/run/maddy/ (runtime_dir)
    ├── maddy.sock      # Control socket
    └── *.pid           # PID files
```

### Persistence Points

| Component | Storage | Format |
|-----------|---------|--------|
| IMAP mailboxes | SQLite/PG/MySQL | imapsql schema |
| Message blobs | Filesystem/S3 | Raw message |
| Message queue | Filesystem | JSON + blob |
| DKIM keys | Filesystem | PEM |
| ACME certs | Filesystem | certmagic format |
