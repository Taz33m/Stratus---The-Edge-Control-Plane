.PHONY: help install dev build up down logs clean test

# Default target
help:
	@echo "Stratus - Edge Control Plane"
	@echo ""
	@echo "Available commands:"
	@echo "  make install    - Install all dependencies"
	@echo "  make dev        - Run in development mode"
	@echo "  make build      - Build all services"
	@echo "  make up         - Start all services with Docker Compose"
	@echo "  make down       - Stop all services"
	@echo "  make logs       - View service logs"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make test       - Run tests"

# Install dependencies
install:
	@echo "📦 Installing backend dependencies..."
	cd backend && go mod download
	@echo "📦 Installing frontend dependencies..."
	cd frontend && npm install
	@echo "✅ Dependencies installed"

# Run in development mode
dev:
	@echo "🚀 Starting development servers..."
	docker-compose up postgres redis -d
	@echo "Starting backend..."
	cd backend && go run main.go &
	@echo "Starting frontend..."
	cd frontend && npm run dev

# Build all services
build:
	@echo "🔨 Building all services..."
	docker-compose build

# Start services
up:
	@echo "🚀 Starting all services..."
	docker-compose up -d
	@echo "✅ Services started"
	@echo "   Frontend: http://localhost:3000"
	@echo "   Backend:  http://localhost:8080"

# Stop services
down:
	@echo "⏹️  Stopping all services..."
	docker-compose down
	@echo "✅ Services stopped"

# View logs
logs:
	docker-compose logs -f

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	cd backend && go clean
	cd frontend && rm -rf .next out node_modules
	docker-compose down -v
	@echo "✅ Cleaned"

# Run tests
test:
	@echo "🧪 Running backend tests..."
	cd backend && go test ./...
	@echo "🧪 Running frontend tests..."
	cd frontend && npm test
	@echo "✅ Tests complete"
