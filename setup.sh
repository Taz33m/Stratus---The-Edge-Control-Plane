#!/bin/bash

# Stratus Setup Script
# This script sets up the development environment for Stratus

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Header
echo -e "${BLUE}"
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║                                                               ║"
echo "║                  ☁️  STRATUS SETUP SCRIPT                     ║"
echo "║          Command your microservices from the edge            ║"
echo "║                                                               ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Function to print status messages
print_status() {
    echo -e "${BLUE}==>${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Step 1: Check prerequisites
print_status "Checking prerequisites..."

MISSING_DEPS=0

if ! command_exists docker; then
    print_error "Docker is not installed"
    echo "  Install from: https://docs.docker.com/get-docker/"
    MISSING_DEPS=1
else
    print_success "Docker found: $(docker --version)"
fi

if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
    print_error "Docker Compose is not installed"
    echo "  Install from: https://docs.docker.com/compose/install/"
    MISSING_DEPS=1
else
    if command_exists docker-compose; then
        print_success "Docker Compose found: $(docker-compose --version)"
    else
        print_success "Docker Compose found: $(docker compose version)"
    fi
fi

if ! command_exists go; then
    print_warning "Go is not installed (optional for local development)"
    echo "  Install from: https://golang.org/doc/install"
else
    print_success "Go found: $(go version)"
fi

if ! command_exists node; then
    print_warning "Node.js is not installed (optional for local development)"
    echo "  Install from: https://nodejs.org/"
else
    print_success "Node.js found: $(node --version)"
    if ! command_exists npm; then
        print_error "npm is not installed"
        MISSING_DEPS=1
    else
        print_success "npm found: $(npm --version)"
    fi
fi

if [ $MISSING_DEPS -eq 1 ]; then
    print_error "Please install missing dependencies and run this script again"
    exit 1
fi

echo ""

# Step 2: Setup environment files
print_status "Setting up environment files..."

# Backend .env
if [ ! -f "backend/.env" ]; then
    cp backend/.env.example backend/.env
    print_success "Created backend/.env from .env.example"
else
    print_warning "backend/.env already exists, skipping..."
fi

# Frontend .env.local
if [ ! -f "frontend/.env.local" ]; then
    cp frontend/.env.local.example frontend/.env.local
    print_success "Created frontend/.env.local from .env.local.example"
else
    print_warning "frontend/.env.local already exists, skipping..."
fi

echo ""

# Step 3: Install Go dependencies
if command_exists go; then
    print_status "Installing Go dependencies..."
    cd backend
    go mod download
    go mod tidy
    print_success "Go dependencies installed"
    cd ..
    echo ""
else
    print_warning "Skipping Go dependencies (Go not installed)"
    echo ""
fi

# Step 4: Install Node dependencies
if command_exists npm; then
    print_status "Installing Node.js dependencies..."
    cd frontend
    npm install
    print_success "Node.js dependencies installed"
    cd ..
    echo ""
else
    print_warning "Skipping Node.js dependencies (npm not installed)"
    echo ""
fi

# Step 5: Check Docker daemon
print_status "Checking Docker daemon..."
if ! docker info >/dev/null 2>&1; then
    print_error "Docker daemon is not running"
    echo "  Please start Docker and run this script again"
    exit 1
else
    print_success "Docker daemon is running"
fi

echo ""

# Step 6: Build Docker images (optional)
read -p "$(echo -e ${YELLOW}Do you want to build Docker images now? \(y/N\):${NC} )" -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status "Building Docker images..."
    docker-compose build
    print_success "Docker images built successfully"
    echo ""
fi

# Summary
echo -e "${GREEN}"
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║                    🎉 SETUP COMPLETE! 🎉                      ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

echo -e "${BLUE}Next steps:${NC}"
echo ""
echo "1️⃣  Start all services:"
echo -e "   ${YELLOW}docker-compose up${NC}"
echo ""
echo "2️⃣  Or start in detached mode:"
echo -e "   ${YELLOW}docker-compose up -d${NC}"
echo ""
echo "3️⃣  Access the dashboard:"
echo -e "   ${YELLOW}http://localhost:3000${NC}"
echo ""
echo "4️⃣  Access the API:"
echo -e "   ${YELLOW}http://localhost:8080${NC}"
echo ""
echo "5️⃣  View logs:"
echo -e "   ${YELLOW}docker-compose logs -f${NC}"
echo ""
echo "6️⃣  Stop services:"
echo -e "   ${YELLOW}docker-compose down${NC}"
echo ""

echo -e "${BLUE}For local development:${NC}"
echo ""
echo "  Backend:  ${YELLOW}cd backend && go run main.go${NC}"
echo "  Frontend: ${YELLOW}cd frontend && npm run dev${NC}"
echo ""

echo -e "${BLUE}Documentation:${NC}"
echo ""
echo "  📖 Full Guide:      ${YELLOW}README.md${NC}"
echo "  🚀 Quick Start:     ${YELLOW}QUICKSTART.md${NC}"
echo "  🏗️  Architecture:    ${YELLOW}ARCHITECTURE.md${NC}"
echo ""

echo -e "${GREEN}Happy coding! ☁️${NC}"
echo ""
