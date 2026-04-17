# slop

Szymon's Language Overengineered Parser. A Go-based CLI tool for infrastructure orchestration using the `Slopfile`

## Overview

`slop` parses `Slopfile` definitions to automate shell commands, environment variable management, and database seeding.
It is designed to replace repetitive bash scripts with a structured DSL.

## Syntax Specification

### Configuration (`config::`)

Used to set internal state or driver settings.

* `config::name.key["value"]` - Static assignment.
* `config::name.key[$env.VAR]` - Dynamic assignment from system environment.

### Environment Sourcing (`source::`)

* `source::env["path/to/.env"]` - Loads a `.env` file into the process environment for subsequent variable interpolation.

### Variables (`var::`)

* `var::namespace.name["value"]` - Defines a reusable string.
* Variables are referenced using the `$` prefix (e.g., `$namespace.name`).

### Execution (`run::`)

* `run::seed["path/to/file.sql"]` - Executes a MariaDB seed command using the provided SQL file.

### Tasks (`@task-name`)

Tasks currently *only* support run commands

```nginx
# Setup a task
@my-cool-task {
    run::seed["/path/to/seed.sql"]
    run::seed["/path/to/another/seed.sql"]
    run::seed["/path/to/yet/another/seed.sql"]
}

# Run the task like a variable
run::task[$my-cool-task]
```

Slopfile does not populate tasks that are empty
and `slop` command won't recognize them

```nginx
@empty-task {
# won't get rendered anywhere / slop cli won't recognize it
}
```


> [!NOTE]
> Tasks currently cannot be nested in themselves due to parser constraints
> Tasks can invoke an infinite circular cycle: `@a` calling `run::task[$b]` and `@b` calling `run::task[$a]` will cause an infinite runtime loop.

## Example Slopfile

```nginx
# Configure database credentials via environment
config::db.user[$env.DB_USER]
config::db.password[$env.DB_PASSWORD]
# Or via env sourced here
source::env[".env"]

# Define seed paths
var::seed.rbac["/app/seed.sql"]

# Execute seeding
run::seed[$seed.rbac]

```

## Usage

1. Build the binary:
```bash
go build -o slop

```


2. Execute:
```bash
./slop

```


*Default behavior: looks for `Slopfile` in the execution directory.*

## Releases

Releases are automated via [goreleaser](https://goreleaser.com/) on tag pushes.

Cut a new release:

```bash
git tag v0.1.0
git push --tags
```

The GitHub Actions release workflow builds Linux/macOS binaries (amd64 + arm64), publishes them as a GitHub Release, and attaches checksums.

Install from a release:

```bash
curl -L https://github.com/thegroobi/slop/releases/latest/download/slop_Linux_x86_64.tar.gz | tar xz
sudo mv slop /usr/local/bin/
```
