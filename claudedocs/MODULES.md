# Maddy Module Reference

## Module System Overview

All Maddy functionality is provided by modules that implement the `module.Module` interface:

```go
type Module interface {
    Init(*config.Map) error  // Initialize with configuration
    Name() string            // Module type name
    InstanceName() string    // Instance name (from config)
}
```

## Endpoints

Endpoints are special modules that listen for incoming connections.

### smtp

**Location**: `internal/endpoint/smtp/`

SMTP server endpoint supporting MTA, MX, and submission roles.

```
smtp tcp://0.0.0.0:25 {
    hostname mx.example.com
    tls &tls_config

    deliver_to &pipeline
}
```

**Key Directives**:
| Directive | Description |
|-----------|-------------|
| `hostname` | Server hostname for EHLO |
| `tls` | TLS configuration reference |
| `auth` | Authentication provider |
| `deliver_to` | Delivery target |
| `source` | Message pipeline configuration |
| `limits` | Connection/rate limits |

### imap

**Location**: `internal/endpoint/imap/`

IMAP server for mailbox access.

```
imap tls://0.0.0.0:993 {
    tls &tls_config
    auth &auth_provider
    storage &local_storage
}
```

### openmetrics

**Location**: `internal/endpoint/openmetrics/`

Prometheus-compatible metrics endpoint.

```
openmetrics tcp://127.0.0.1:9749 {}
```

### dovecot_sasld

**Location**: `internal/endpoint/dovecot_sasld/`

Dovecot SASL authentication daemon.

## Targets

Targets handle message delivery.

### target.queue

**Location**: `internal/target/queue/`

Queues messages for asynchronous delivery.

```
target.queue outbound_queue {
    target &remote_target

    bounce {
        destination &local_mailboxes
    }
}
```

### target.remote

**Location**: `internal/target/remote/`

Delivers messages to remote SMTP servers via MX lookup.

```
target.remote outbound {
    hostname mx.example.com

    limits {
        destination rate 20 1s
    }
}
```

### target.smtp

**Location**: `internal/target/smtp/`

Relays messages to a specific SMTP server.

```
target.smtp relay {
    targets tcp://smtp.relay.com:25
    auth plain user password
}
```

### target.apprise

**Location**: `internal/target/apprise/` (if present)

Sends notifications via Apprise.

```
target.apprise notifications {
    apprise_url http://localhost:8000/notify
}
```

## Authentication Providers

### auth.pam

**Location**: `internal/auth/pam/`

PAM authentication (requires `libpam` build tag).

```
auth.pam system_auth {
    service maddy
}
```

### auth.shadow

**Location**: `internal/auth/shadow/`

Shadow password file authentication.

```
auth.shadow system_auth {
    # Uses maddy-shadow-helper
}
```

### auth.ldap

**Location**: `internal/auth/ldap/`

LDAP authentication.

```
auth.ldap ldap_auth {
    urls ldap://ldap.example.com
    bind_dn cn=maddy,ou=services,dc=example,dc=com
    bind_password secret

    base_dn ou=users,dc=example,dc=com
    filter (&(objectClass=posixAccount)(uid={username}))
}
```

### auth.pass_table

**Location**: `internal/auth/pass_table/`

Password table-based authentication.

```
auth.pass_table local_auth {
    table sql_table {
        driver sqlite3
        dsn credentials.db
        lookup "SELECT password FROM users WHERE email = $1"
    }
}
```

### auth.external

**Location**: `internal/auth/external/`

External helper program authentication.

```
auth.external custom_auth {
    helper /usr/local/bin/auth-helper
}
```

## Message Checks

### check.dkim

**Location**: `internal/check/dkim/`

DKIM signature verification.

```
check.dkim {
    require_signature false
    allow_body_subset true
}
```

### check.spf

**Location**: `internal/check/spf/`

SPF record verification.

```
check.spf {
    fail_action reject
    softfail_action quarantine
}
```

### check.dnsbl

**Location**: `internal/check/dnsbl/`

DNS blocklist checking.

```
check.dnsbl {
    reject_threshold 1

    dnsbl zen.spamhaus.org
    dnsbl bl.spamcop.net
}
```

### check.rspamd

**Location**: `internal/check/rspamd/`

Rspamd spam filtering integration.

```
check.rspamd {
    api_url http://127.0.0.1:11333
}
```

### check.milter

**Location**: `internal/check/milter/`

Milter protocol integration.

```
check.milter clamav {
    endpoint tcp://127.0.0.1:3310
}
```

### check.authorize_sender

**Location**: `internal/check/authorize_sender/`

Ensures authenticated user matches sender.

```
check.authorize_sender {
    prepare_email identity
}
```

### check.requiretls

**Location**: `internal/check/requiretls/`

Requires TLS for delivery.

```
check.requiretls {
    # No configuration needed
}
```

### check.command

**Location**: `internal/check/command/`

External command check.

```
check.command custom_check {
    cmd /usr/local/bin/check-message
}
```

## Message Modifiers

### modify.dkim

**Location**: `internal/modify/dkim/`

DKIM signature addition.

```
modify.dkim {
    domain example.com
    selector default
    key_path /var/lib/maddy/dkim_keys/example.com.key
}
```

### modify.replace_addr

**Location**: `internal/modify/`

Address replacement/mapping.

```
modify.replace_addr sender {
    table static {
        entry old@example.com new@example.com
    }
}
```

## Storage

### storage.imapsql

**Location**: `internal/storage/imapsql/`

SQL-based IMAP storage.

```
storage.imapsql local_storage {
    driver sqlite3
    dsn /var/lib/maddy/imapsql/imapsql.db

    appendlimit 32M

    imap_filter {
        # IMAP-side filtering rules
    }
}
```

**Supported Drivers**:
- `sqlite3` - SQLite (default)
- `postgres` - PostgreSQL
- `mysql` - MySQL/MariaDB

### storage.blob.fs

**Location**: `internal/storage/blob/fs/`

Filesystem blob storage.

```
storage.blob.fs blob_store {
    path /var/lib/maddy/blobs
}
```

### storage.blob.s3

**Location**: `internal/storage/blob/s3/`

S3-compatible blob storage.

```
storage.blob.s3 blob_store {
    endpoint s3.amazonaws.com
    bucket maddy-blobs
    access_key AKIA...
    secret_key ...
}
```

## TLS Configuration

### tls.loader

**Location**: `internal/tls/`

TLS certificate configuration.

```
tls file /etc/maddy/certs/cert.pem /etc/maddy/certs/key.pem

# Or with multiple certificates
tls {
    loader file {
        certificates /etc/maddy/certs
    }
}
```

### tls.acme

**Location**: `internal/tls/acme/`

ACME (Let's Encrypt) certificate automation.

```
tls {
    loader acme {
        staging off
        email admin@example.com

        challenge dns-01 {
            provider cloudflare {
                api_token token123
            }
        }
    }
}
```

## Tables

Tables provide lookup functionality for address mapping, routing, etc.

### Common Table Types

**static** - Static key-value mapping:
```
table static {
    entry key1 value1
    entry key2 value2
}
```

**file** - File-based lookup:
```
table file /etc/maddy/aliases
```

**sql_table** - SQL-based lookup:
```
table sql_table {
    driver sqlite3
    dsn /var/lib/maddy/tables.db
    lookup "SELECT value FROM aliases WHERE key = $1"
}
```

**regexp** - Regular expression mapping:
```
table regexp {
    entry "(.+)@old.com" "$1@new.com"
}
```

## Module Registration Pattern

All modules follow this registration pattern:

```go
package mymodule

import (
    "github.com/foxcpp/maddy/framework/config"
    "github.com/foxcpp/maddy/framework/module"
)

type MyModule struct {
    modName  string
    instName string
    // ... fields
}

func init() {
    module.Register("my_module", New)
}

func New(modName, instName string, aliases, inlineArgs []string) (module.Module, error) {
    return &MyModule{
        modName:  modName,
        instName: instName,
    }, nil
}

func (m *MyModule) Name() string         { return m.modName }
func (m *MyModule) InstanceName() string { return m.instName }

func (m *MyModule) Init(cfg *config.Map) error {
    // Parse configuration
    cfg.String("my_option", false, false, "default", &m.option)

    if _, err := cfg.Process(); err != nil {
        return err
    }

    // Initialize module
    return nil
}
```

## Creating New Modules

### Check Module Template

Use `internal/check/skeleton.go` as starting point:

```go
// Implement module.Check interface
type MyCheck struct {
    modName  string
    instName string
}

// Per-message state
type myCheckState struct {
    c   *MyCheck
    // ... per-message state
}

func (c *MyCheck) CheckStateForMsg(ctx context.Context, msgMeta *module.MsgMetadata) (module.CheckState, error) {
    return &myCheckState{c: c}, nil
}

func (s *myCheckState) CheckConnection(ctx context.Context) module.CheckResult { ... }
func (s *myCheckState) CheckSender(ctx context.Context, from string) module.CheckResult { ... }
func (s *myCheckState) CheckRcpt(ctx context.Context, rcpt string) module.CheckResult { ... }
func (s *myCheckState) CheckBody(ctx context.Context, header textproto.MIMEHeader, body buffer.Buffer) module.CheckResult { ... }
func (s *myCheckState) Close() error { ... }
```

### Stateless Check

For simple checks without per-message state:

```go
type MyCheck struct {
    check.StatelessCheck
    // ... fields
}

func (c *MyCheck) CheckConnection(ctx context.Context, state *smtp.ConnectionState) error { ... }
func (c *MyCheck) CheckSender(ctx context.Context, from string) error { ... }
func (c *MyCheck) CheckRcpt(ctx context.Context, rcpt string) error { ... }
func (c *MyCheck) CheckBody(ctx context.Context, state *smtp.ConnectionState, header textproto.MIMEHeader, body buffer.Buffer) error { ... }
```
