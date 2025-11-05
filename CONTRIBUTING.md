# Contributing to Stratus

First off, thank you for considering contributing to Stratus! It's people like you that make Stratus such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps to reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed and what you expected**
* **Include screenshots if possible**
* **Include your environment details** (OS, Docker version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the enhancement**
* **Describe the current behavior and expected behavior**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Follow the coding style used throughout the project
* Include appropriate test cases
* Update documentation as needed
* End all files with a newline

## Development Process

### 1. Fork & Clone

```bash
# Fork the repository on GitHub, then:
git clone https://github.com/your-username/stratus.git
cd stratus
```

### 2. Set Up Development Environment

```bash
# Run the setup script
chmod +x setup.sh
./setup.sh

# Or manually:
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local
cd backend && go mod download && cd ..
cd frontend && npm install && cd ..
```

### 3. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 4. Make Your Changes

#### Backend (Go)

```bash
cd backend

# Run tests
go test ./...

# Run linter
go vet ./...

# Format code
go fmt ./...

# Run the server
go run main.go
```

#### Frontend (TypeScript/Next.js)

```bash
cd frontend

# Run tests
npm test

# Run linter
npm run lint

# Format code (if prettier is configured)
npm run format

# Run dev server
npm run dev
```

### 5. Commit Your Changes

We follow conventional commits:

```bash
# Feature
git commit -m "feat(api): add support for service health checks"

# Bug fix
git commit -m "fix(ui): resolve websocket reconnection issue"

# Documentation
git commit -m "docs: update API endpoint documentation"

# Refactor
git commit -m "refactor(backend): improve database connection pooling"
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

### 6. Push & Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:
* Clear title and description
* Reference to any related issues
* Screenshots/GIFs for UI changes
* List of changes made

## Coding Standards

### Go (Backend)

* Follow [Effective Go](https://golang.org/doc/effective_go)
* Use `gofmt` for formatting
* Write meaningful variable and function names
* Add comments for exported functions
* Keep functions small and focused
* Write unit tests for new features

Example:
```go
// GetServiceByID retrieves a service from the database by its unique identifier.
// Returns an error if the service is not found or if a database error occurs.
func (h *ServiceHandler) GetServiceByID(id string) (*models.Service, error) {
    // Implementation
}
```

### TypeScript (Frontend)

* Follow TypeScript best practices
* Use functional components with hooks
* Use meaningful component and variable names
* Add proper type annotations
* Keep components small and reusable
* Write unit tests for components

Example:
```typescript
interface ServiceCardProps {
  service: Service
  onStart: (id: string) => void
  onStop: (id: string) => void
}

export function ServiceCard({ service, onStart, onStop }: ServiceCardProps) {
  // Implementation
}
```

### Database Migrations

* Always create reversible migrations
* Test migrations both up and down
* Document schema changes

### API Design

* Follow RESTful conventions
* Use proper HTTP status codes
* Return consistent error formats
* Version your APIs

## Testing

### Backend Tests

```bash
cd backend
go test ./... -v
go test -race ./...
go test -cover ./...
```

### Frontend Tests

```bash
cd frontend
npm test
npm run test:coverage
```

### Integration Tests

```bash
# Start services
docker-compose up -d

# Run integration tests
./scripts/integration-tests.sh
```

## Documentation

* Update README.md for user-facing changes
* Update ARCHITECTURE.md for architectural changes
* Add JSDoc/GoDoc comments for public APIs
* Include examples in documentation

## Release Process

1. Update version in `package.json` and `go.mod`
2. Update CHANGELOG.md
3. Create a git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push tag: `git push origin v1.0.0`
5. Create GitHub release with notes

## Questions?

Feel free to open an issue with the `question` label or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Stratus! 🎉
