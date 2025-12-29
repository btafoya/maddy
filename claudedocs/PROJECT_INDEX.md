# Maddy Mail Server - Project Index

## Quick Reference

| Item | Location/Command |
|------|------------------|
| Main Binary | `./build/maddy` |
| Config File | `/etc/maddy/maddy.conf` or `./maddy.conf` |
| Build | `./build.sh` |
| Test | `go test ./...` |
| Lint | `golangci-lint run` |
| Start Server | `maddy run --config maddy.conf` |

## Project Overview

Maddy is a composable, all-in-one mail server implementing:
- **SMTP** (MTA for sending, MX for receiving)
- **IMAP** (mailbox access)
- **Security**: DKIM, SPF, DMARC, DANE, MTA-STS

**Go Version**: 1.23.1+
**License**: GPL v3.0

## Directory Structure

```
maddy/
├── cmd/                      # Binary entrypoints
│   ├── maddy/               # Main server (28 lines)
│   ├── maddy-pam-helper/    # PAM authentication helper
│   └── maddy-shadow-helper/ # Shadow password helper
│
├── framework/               # Reusable module framework
│   ├── module/             # Module interfaces & registry
│   ├── config/             # Configuration parsing
│   ├── cfgparser/          # Config file parser
│   ├── exterrors/          # SMTP-aware error handling
│   ├── hooks/              # Lifecycle hooks
│   ├── log/                # Logging utilities
│   ├── dns/                # DNS utilities
│   ├── address/            # Email address parsing
│   ├── buffer/             # Buffered I/O utilities
│   └── future/             # Async result handling
│
├── internal/               # Core implementations (unstable API)
│   ├── endpoint/           # Protocol listeners
│   │   ├── smtp/          # SMTP/submission server
│   │   ├── imap/          # IMAP server
│   │   ├── openmetrics/   # Prometheus metrics
│   │   └── dovecot_sasld/ # Dovecot SASL daemon
│   │
│   ├── target/            # Delivery targets
│   │   ├── queue/         # Message queuing
│   │   ├── remote/        # Remote SMTP delivery
│   │   ├── smtp/          # SMTP relay
│   │   └── apprise/       # Notification delivery
│   │
│   ├── auth/              # Authentication providers
│   │   ├── pam/           # PAM authentication
│   │   ├── shadow/        # Shadow password file
│   │   ├── ldap/          # LDAP authentication
│   │   ├── netauth/       # NetAuth
│   │   ├── external/      # External helper programs
│   │   ├── pass_table/    # Password table lookup
│   │   ├── plain_separate/# Username/password separation
│   │   └── dovecot_sasl/  # Dovecot SASL proxy
│   │
│   ├── check/             # Message checks
│   │   ├── dkim/          # DKIM verification
│   │   ├── spf/           # SPF verification
│   │   ├── dns/           # DNS checks (MX, PTR)
│   │   ├── dnsbl/         # DNS blocklists
│   │   ├── rspamd/        # Rspamd integration
│   │   ├── milter/        # Milter protocol
│   │   ├── requiretls/    # TLS requirement
│   │   ├── authorize_sender/ # Sender auth
│   │   └── command/       # External command
│   │
│   ├── modify/            # Message modifiers
│   │   ├── dkim/          # DKIM signing
│   │   └── replace_addr.go # Address replacement
│   │
│   ├── storage/           # Storage backends
│   │   ├── imapsql/       # SQL-based IMAP storage
│   │   └── blob/          # Blob storage (fs, s3)
│   │
│   ├── tls/               # TLS configuration
│   │   └── acme/          # ACME automation
│   │
│   ├── table/             # Lookup tables
│   ├── libdns/            # DNS providers for ACME
│   ├── msgpipeline/       # Message processing pipeline
│   ├── dmarc/             # DMARC policy handling
│   ├── dsn/               # Delivery status notifications
│   ├── limits/            # Rate limiting
│   ├── authz/             # Authorization helpers
│   ├── cli/               # CLI utilities
│   ├── imap_filter/       # IMAP filtering
│   └── testutils/         # Test utilities
│
├── docs/                   # Documentation
│   ├── man/               # Man pages (scdoc format)
│   ├── reference/         # Configuration reference
│   ├── tutorials/         # Setup guides
│   └── internals/         # Internal documentation
│
├── tests/                  # Integration tests
├── dist/                   # Distribution files (systemd)
├── contrib/                # Community contributions
└── claudedocs/            # AI-generated documentation
```

## Key Files

| File | Purpose |
|------|---------|
| `maddy.go` | Server initialization, module registration |
| `config.go` | Global configuration handling |
| `signal.go` | Signal handling (graceful shutdown) |
| `directories.go` | Default directory paths |
| `build.sh` | Build and installation script |
| `maddy.conf` | Example configuration |
| `.golangci.yml` | Linter configuration |
| `HACKING.md` | Developer documentation |

## Module System

### Core Interface
```go
type Module interface {
    Init(*config.Map) error
    Name() string
    InstanceName() string
}
```

### Registration Pattern
```go
func init() {
    module.Register("module.name", New)
}
```

### Configuration Pattern
```
module_type instance_name {
    directive value
}

# Reference existing instance
&instance_name
```

## Build Tags

| Tag | Description |
|-----|-------------|
| `libpam` | Enable PAM authentication |
| `nosqlite3` | Disable SQLite support |
| `integration` | Enable integration tests |

## Important Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/emersion/go-smtp` | SMTP implementation |
| `github.com/emersion/go-imap` | IMAP implementation |
| `github.com/emersion/go-msgauth` | DKIM/DMARC/SPF |
| `github.com/caddyserver/certmagic` | ACME automation |
| `github.com/miekg/dns` | DNS queries |
| `github.com/mattn/go-sqlite3` | SQLite driver |
| `github.com/urfave/cli/v2` | CLI framework |
| `go.uber.org/zap` | Logging |

## Cross-References

- [Architecture Notes](./ARCHITECTURE.md)
- [Module Reference](./MODULES.md)
- [API Documentation](./API.md)
- [Configuration Guide](./CONFIGURATION.md)
