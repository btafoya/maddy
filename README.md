Maddy Mail Server
=====================
> Composable all-in-one mail server.

Maddy Mail Server implements all functionality required to run a e-mail
server. It can send messages via SMTP (works as MTA), accept messages via SMTP
(works as MX) and store messages while providing access to them via IMAP.
In addition to that it implements auxiliary protocols that are mandatory
to keep email reasonably secure (DKIM, SPF, DMARC, DANE, MTA-STS).

It replaces Postfix, Dovecot, OpenDKIM, OpenSPF, OpenDMARC and more with one
daemon with uniform configuration and minimal maintenance cost.

**Note:** IMAP storage is "beta". If you are looking for stable and
feature-packed implementation you may want to use Dovecot instead. maddy still
can handle message delivery business.

[![CI status](https://img.shields.io/github/actions/workflow/status/foxcpp/maddy/cicd.yml?style=flat-square)](https://github.com/foxcpp/maddy/actions/workflows/cicd.yml)
[![Issues tracker](https://img.shields.io/github/issues/foxcpp/maddy?style=flat-square)](https://github.com/foxcpp/maddy)

* [Setup tutorial](https://maddy.email/tutorials/setting-up/)
* [Documentation](https://maddy.email/)

* [IRC channel](https://webchat.oftc.net/?channels=maddy&uio=MT11bmRlZmluZWQb1)
* [Mailing list](https://lists.sr.ht/~foxcpp/maddy)

## Features

- **SMTP/IMAP Server** - Full MTA, MX, and IMAP support in a single binary
- **Security Standards** - Built-in DKIM signing/verification, SPF, DMARC, DANE, MTA-STS
- **Modular Architecture** - Composable modules for custom mail pipelines
- **SQLite Storage** - Default storage backend with no external dependencies
- **Notifications** - Apprise integration for delivery notifications
- **Authentication** - Multiple auth backends including PAM, shadow, and LDAP

## Building from Source

### Prerequisites

- Go 1.18 or later
- C compiler (for SQLite support)
- Optional: `scdoc` for man pages
- Optional: `libpam-dev` for PAM authentication

### Basic Build

```bash
# Clone the repository
git clone https://github.com/foxcpp/maddy.git
cd maddy

# Build the server
./build.sh

# The binary will be in ./build/maddy
```

### Static Build

For a self-contained binary without external dependencies:

```bash
./build.sh --static
```

### Build Tags

Optional features can be enabled with build tags:

| Tag | Description | Dependency |
|-----|-------------|------------|
| `libpam` | PAM authentication support | `libpam-dev` |
| `nosqlite3` | Disable SQLite (use external DB) | None |

Example with PAM support:
```bash
go build -tags libpam ./cmd/maddy
```

### Installation

```bash
# Install to /usr/local (default)
./build.sh install

# Custom prefix
./build.sh --prefix=/opt/maddy install

# Custom destination root (for packaging)
./build.sh --destdir=/tmp/package install
```

## Optional Components

### PAM Helper

The `maddy-pam-helper` binary enables PAM authentication. It requires the PAM development headers.

**Building with PAM support:**

```bash
# Install PAM development headers first:
# Debian/Ubuntu: apt install libpam0g-dev
# Fedora/RHEL:   dnf install pam-devel
# Alpine:        apk add linux-pam-dev
# Arch:          pacman -S pam

# Build the helper
go build -tags libpam ./cmd/maddy-pam-helper
```

**Without PAM headers**, the helper builds as a stub that displays installation instructions:

```
maddy-pam-helper: PAM support not available

This binary was built without PAM support. To enable PAM authentication:
1. Install PAM development headers
2. Rebuild with: go build -tags libpam ./cmd/maddy-pam-helper
```

See [PAM authentication documentation](https://maddy.email/reference/auth/pam/) for setup details.

### Apprise Notifications

The `target.apprise` module sends notifications via [Apprise](https://github.com/caronc/apprise) when mail is delivered.

**Configuration example:**

```
target.apprise local_notify {
    apprise_url http://localhost:8000/notify
}

# Use in delivery pipeline
deliver_to &local_notify
```

**Features:**
- Send notifications to any Apprise-supported service (Slack, Discord, email, SMS, etc.)
- Multiple notification URLs supported
- Includes message snippet in notification body

## Running

```bash
# Start the server
maddy run --config /path/to/maddy.conf

# Default config location: /etc/maddy/maddy.conf
```

See the [setup tutorial](https://maddy.email/tutorials/setting-up/) for complete configuration guidance.

## Testing

```bash
# Run unit tests
go test ./...

# Run integration tests
cd tests && ./run.sh
```

## Project Structure

```
maddy/
├── cmd/
│   ├── maddy/              # Main server binary
│   ├── maddy-pam-helper/   # PAM authentication helper
│   └── maddy-shadow-helper/ # Shadow password helper
├── internal/               # Core implementation
│   ├── auth/              # Authentication modules
│   ├── check/             # Message checks (SPF, DKIM, etc.)
│   ├── modify/            # Message modifiers
│   ├── storage/           # Storage backends
│   └── target/            # Delivery targets (including apprise)
├── framework/             # Module system and config
└── docs/                  # Documentation sources
```

## Contributing

Contributions are welcome! Please see the [contribution guidelines](https://maddy.email/dev/) for details.

## License

Maddy is licensed under the GNU General Public License v3.0. See [LICENSE](LICENSE) for details.
