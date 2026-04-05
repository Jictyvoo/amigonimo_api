# Amigonimo API

Amigonimo is a Secret Friend (Secret Santa) management API. It lets users create groups, invite
participants, manage wishlists and denylists, and run an automatic draw that respects all pairing
constraints.

## Documentation

| Document                     | Description                                                                     |
|------------------------------|---------------------------------------------------------------------------------|
| [FEATURES.md](FEATURES.md)   | What the API currently provides — endpoints, business rules, and infrastructure |
| [USER_FLOW.md](USER_FLOW.md) | Step-by-step flow from group creation to result reveal, with a flowchart        |
| [ROADMAP.md](ROADMAP.md)     | Outstanding work items                                                          |

---

## Architecture

The project follows a clean architecture approach with three main layers:

- **`pkg/web/handlers`** — HTTP handlers; entry point for all external requests.
- **`internal/domain`** — Core business logic and use cases; has no HTTP dependency.
- **`internal/infra`** — Infrastructure adapters: database repositories, mailer backends.

## Tech Stack

| Concern              | Tool                                              |
|----------------------|---------------------------------------------------|
| Language             | Go                                                |
| Web framework        | [Fuego](https://fuego.go-lab.io/)                 |
| Dependency injection | [Remy](https://github.com/wrapped-owls/goremy-di) |
| Schema-as-code       | [Ent](https://entgo.io/)                          |
| Migrations           | [Atlas](https://atlasgo.io/)                      |
| Type-safe queries    | [sqlc](https://sqlc.dev/)                         |
| DTO mapping          | [Goverter](https://github.com/jmattheis/goverter) |
| Auth                 | JWT (RSA key pairs)                               |
| Database             | MariaDB                                           |
| Task runner          | [Taskfile](https://taskfile.dev)                  |
| Containerization     | Docker & Docker Compose                           |

## Database & Data Layer

- **Ent** defines the schema as Go code (`build/entschema/`).
- **Atlas** generates versioned migration files from the schema.
- **sqlc** generates type-safe query functions and table models from `.sql` files
  (`internal/infra/repositories/mysqlrepo/internal/dbgen/`).
- Table models are internal to the repository layer. **Goverter** generates compile-time-checked
  mappers that convert them to domain entities, preventing forgotten-field bugs.

## Authentication & Per-Request Context

- JWT tokens are validated by a middleware that extracts the current user and stores it in the
  request context.
- **Remy** reads that context and injects the current user into each use case automatically. Every
  request gets a fresh, stateless use case instance with the correct user already available — no
  manual threading of user arguments through every call.

## Development Workflow

```bash
# Start the database and prepare the environment
task dev-env:up

# Apply migrations
cd build && task migration:apply

# Build
task build

# Run all tests
go test ./...

# Run formatters
task run:formatters
```

