# Agents Documentation

## Overall

- Always review and comprehend the necessary files before suggesting edits. Avoid speculating on code that hasn't been
  inspected.
- Prioritize readability over backward compatibility or migrations. Significant changes are acceptable if they lead to
  cleaner code.
- Focus on requested changes and maintain simplicity. Feel free to refactor surrounding code if it substantially
  improves clarity or correctness.

## Tooling

- `go task` is used as the task runner (see `Taskfile.yml`, `Taskfile.build.yml`, `Taskfile.docker.yml`, and
  `build/Taskfile.yml`).
- `sqlc` is used for generating type-safe Go code from SQL (see `sqlc.yaml`).
- `atlas` and `ent` are used for database migrations and schema management.
- `golines` and `gofumpt` are used for formatting Go code.
- `golangci-lint` is the linter that must be setup and used for the project.
- `mockgen` is used for generating mocks.
- `goverter` is used for generating type converters.
- `git` is used for version control.

## Development Workflow

These should generally be run from the root of the repo unless specified otherwise.

- `task install:deps` to install all needed dependencies to run the project.
- `task gen:code` to generate all needed code (`go generate` and `sqlc generate`).
- `task run:formatters` to run `golines` and `gofumpt` on the entire project.
- `task test` to execute tests.
- `task dev-env:up` to start the development environment using Docker.
- In the `build/` directory:
    - `task migration:gen DIFF_NAME=<name>` to generate a new migration using Atlas and Ent.
    - `task migration:apply` to apply versioned migrations to the database.

**Important**: Run `task run:formatters` frequently to format your code. You should also ensure `golangci-lint` passes.

- Do NOT fix formatting issues manually if they can be handled by `task run:formatters`.
- When spawning subagents, use foreground ones, not background subagents.

## Project Structure

This is a Go project structured as follows:

- `cmd/` - Main applications (e.g., `cmd/api` for the main API).
- `internal/` - Private application and library code.
    - `internal/entities` - Shared domain entities used across the project.
    - `internal/domain` - Core business logic, following a vertical slicing pattern.
        - `internal/domain/usecases` - Specific application actions or business flows (orchestrate
          services/repositories).
        - `internal/domain/services` - Reusable domain-specific logic, agnostic of usecases.
        - **Self-Enclosed Packages**: Each usecase and service should be self-enclosed in its own package (vertical
          slicing).
        - **Interface & Type Definition**: Each package must define its own interfaces for external dependencies (e.g.,
          repositories or other services) and its own types for input/output.
    - `internal/infra/repositories/mysqlrepo` - MySQL repository implementations.
        - `internal/infra/repositories/mysqlrepo/internal/queries` - SQL queries for `sqlc`.
        - `internal/infra/repositories/mysqlrepo/internal/dbgen` - Generated code by `sqlc`.
- `pkg/` - Library code that is okay to be used by external applications.
    - `pkg/web/handlers` - HTTP request handlers.
- `api/` - API definitions and related files.
- `build/` - Build-related configurations and scripts.
    - `build/migrations` - SQL migration files managed by Atlas.
    - `build/entschema` - Ent schema definitions.
- `docs/` - Documentation files.
