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
