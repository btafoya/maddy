# Maddy Mail Server (`github.com/foxcpp/maddy`)

## Gemini Autonomous Issue Workflow (Required)

Gemini is responsible for resolving all open issues documented as `ISSUE{number}.md` in this repository. Treat each issue file as the source of truth for scope, acceptance criteria, and tracking.

### Operating Principles
- Work **autonomously**: do not ask for human confirmation unless the issue is ambiguous or lacks required information to proceed safely.
- Handle **one issue at a time** from start to finish before moving to the next.
- Prefer **minimal, correct changes** that align with Maddy’s architecture and style.
- If a proposed fix affects public behavior or compatibility, document the impact explicitly in the issue file and ensure tests cover it.

### Standard Flow Per Issue
For each `ISSUE{number}.md`:

1. **Read & interpret**
   - Extract: problem statement, repro steps (if any), expected behavior, acceptance criteria, constraints, and affected areas.
   - If repro steps are missing, derive them from code context and add them to the issue file.

2. **Assess & plan**
   - Identify root cause and the smallest viable fix.
   - Identify test strategy: unit tests, integration tests, or both.
   - Note any risks (regressions, performance, security, compatibility).

3. **Implement**
   - Use context7 mcp for reference.
   - Make changes following existing project conventions in:
     - `cmd/maddy/main.go` (entrypoint) as needed
     - `internal/` packages for core logic
     - `framework/` for module/config mechanics
   - Keep changes focused and avoid unrelated refactors.

4. **Test**
   - Run all relevant tests for the change:
     - `go test ./...`
     - If applicable, `cd tests && ./run.sh`
   - If tests are not applicable or cannot run in the current environment, document precisely:
     - what could not be run
     - why
     - what was done instead (e.g., targeted compilation, narrow test, static reasoning)

5. **Build verification**
   - Verify the build completes successfully:
     - `./build.sh`
   - If the issue concerns static builds or packaging, also run:
     - `./build.sh --static` (when feasible)

6. **Document the resolution in `ISSUE{number}.md`**
   Update the issue file with a clear resolution record including:
   - **Root cause** (what was wrong and where)
   - **Fix summary** (what changed at a high level)
   - **Files changed** (list)
   - **How to reproduce** (before/after if relevant)
   - **How to verify** (commands run and outcomes)
   - **Test results** (what passed; include failures and remediation if any)
   - **Notes / follow-ups** (optional), including any edge cases, limitations, or future improvements

7. **Github Commit, push, and pull request**
   - For each issue git commit, push, and create pull request for each issues. Do not submit the issue to the parent repo!!!

8. **Close-out criteria**
   - The issue is considered resolved only when:
     - the solution is implemented
     - tests relevant to the change pass
     - the build succeeds
     - `ISSUE{number}.md` is updated with the complete solution details

### Issue Ordering
- If `PR.md` defines priority, address issues in that order.
- Otherwise, address issues in **ascending issue number** (`ISSUE1.md`, `ISSUE2.md`, …).
- If a later issue blocks an earlier one due to dependencies, document that dependency in both issue files.

### Output Expectations
- Every resolved issue must leave behind a high-quality `ISSUE{number}.md` update that enables another engineer to:
  - understand the bug
  - understand the fix
  - reproduce and verify confidently

### Safety & Scope
- Do not introduce insecure defaults. If an issue touches authentication, TLS, DKIM/SPF/DMARC, or mail routing, include a brief security impact note in the issue file.
- Avoid configuration-breaking changes unless the issue explicitly requires it; if unavoidable, document migration steps.

## Project Overview

Maddy is a composable, all-in-one mail server written in Go. It implements all the functionality required to run a modern, secure e-mail server. It is designed to be a single daemon that replaces complex stacks of software like Postfix, Dovecot, OpenDKIM, OpenSPF, and OpenDMARC. This simplifies configuration and maintenance.

Key features include:
- SMTP for sending (MTA) and receiving (MX) messages.
- IMAP for message storage and access.
- Support for modern security standards like DKIM, SPF, DMARC, DANE, and MTA-STS.
- Modular architecture with a clear configuration file format.

The project is structured as a standard Go application, with a main entrypoint in `cmd/maddy/main.go`, core logic separated into `internal/` packages, and a reusable `framework/` for the module system and configuration parsing.

## Building and Running

### Building from Source

The project includes a shell script for building the binaries.

- **To build the server:**
  ```sh
  ./build.sh
  ```
  This will create the `maddy` executable in the `./build` directory.

- **To create a static, self-contained build:**
  ```sh
  ./build.sh --static
  ```

- **To install the server:**
  ```sh
  ./build.sh install
  ```
  This will install the binary to `/usr/local/bin` and the configuration to `/etc/maddy` by default. These paths can be customized with `--prefix` and `--destdir`.

### Running the Server

The server is started using the `run` subcommand.

```sh
maddy run --config /path/to/maddy.conf
```

The default configuration file path is `/etc/maddy/maddy.conf`. An example `maddy.conf` is available in the root of the repository.

### Running with Docker

A `Dockerfile` is provided for building and running Maddy in a container.

- **To build the Docker image:**
  ```sh
  docker build -t maddy-server .
  ```

- **To run the container:**
  ```sh
  docker run -p 25:25 -p 143:143 -p 465:465 -p 587:587 -p 993:993 -v /path/to/data:/data maddy-server
  ```
  The container expects a volume at `/data` for persistent state, including configuration (`/data/maddy.conf`), certificates, and mail storage.

## Development

### Testing

The project has both unit/module tests and a suite of integration tests.

- **Run unit and module tests:**
  ```sh
  go test ./...
  ```

- **Run integration tests:**
  ```sh
  cd tests/
  ./run.sh
  ```

### Linting

The project uses `golangci-lint` to enforce code quality and style. The configuration is in `.golangci.yml`.

- **To run the linter:**
  ```sh
  golangci-lint run
  ```

### Dependencies

The project uses Go modules for dependency management. Dependencies are declared in the `go.mod` file.


## Doumentation

### PR

Review PR.md to determine the focus on the current tasks to address any open issues in order to provide pull requests to the parent repo.

### Issues

Docuemnt all issues and their status using markdown with them naming format: ISSUE{number}.md

