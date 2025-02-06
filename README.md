# Cogniboard

Cogniboard is a full-stack task management application built with modern technologies and best practices.

## Tech Stack

### Backend
- Go
- PostgreSQL (via Docker)
- OpenAPI 3.0 specification
- Domain-Driven Design patterns

### Frontend
- React
- TypeScript
- Tailwind CSS
- React Query (with auto-generated hooks via Orval)
- Vite
- Biome (for code formatting and linting)

### Development Tools
- [Devbox](https://www.jetpack.io/devbox) - Instant, easy, and predictable development environments
- [Task](https://taskfile.dev) - Task runner / build tool
- Docker Compose - Container orchestration for development services

## Project Structure

```
.
├── cmd/                    # Command line applications
│   └── http/              # HTTP API server
├── internal/              # Private application code
│   ├── decorator/         # Decorators for logging and instrumentation
│   ├── postgres/          # Database layer
│   └── project/          # Project domain logic
│       ├── adapters/     # Infrastructure adapters
│       └── app/         # Application services
├── web/                   # Frontend application
│   ├── app/              # React application code
│   ├── components/       # Reusable UI components
│   └── hooks/           # Custom React hooks
└── openapi3.yaml         # API specification
```

## Getting Started

### Prerequisites

- [Task](https://taskfile.dev) - Task runner
- [Docker](https://www.docker.com/) - For running PostgreSQL

### Development Setup

1. Copy the example environment file and modify if needed:
```bash
cp .env.example .env
```

2. Install development dependencies:
```bash
task install
```

2. Set up the development environment:
```bash
task setup
```

3. Start the development services (PostgreSQL):
```bash
task docker
```

4. Run the API service:
```bash
task run-api
```

The API will be available at `http://localhost:8000` by default.

### Frontend Development

1. Navigate to the web directory:
```bash
cd web
```

2. Install dependencies:
```bash
bun install
```

3. Generate React Query hooks from OpenAPI spec:
```bash
task gen-hooks
```

4. Start the development server:
```bash
bun dev
```

The frontend will be available at `http://localhost:5173` by default.

## Available Tasks

- `task install` - Install devbox if not already installed
- `task setup` - Setup development environment
- `task docker` - Start Docker services
- `task run-api` - Run the API service
- `task gen-hooks` - Generate React Query hooks from OpenAPI spec
- `task dev` - Start development environment (includes Docker and API)
- `task web` - Start the web server

## Environment Variables

### API
- `SERVICE_POSTGRES_DSN` - PostgreSQL connection string (default: `postgres://cogniboard:password@localhost:5432/cogniboard?sslmode=disable`)

## Contributing

1. Ensure you have all prerequisites installed
2. Fork the repository
3. Create your feature branch
4. Commit your changes
5. Push to the branch
6. Create a new Pull Request