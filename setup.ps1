# Stratus Setup Script for Windows (PowerShell)
# This script sets up the development environment for Stratus

$ErrorActionPreference = "Stop"

# Header
Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                                                               ║" -ForegroundColor Cyan
Write-Host "║                  ☁️  STRATUS SETUP SCRIPT                     ║" -ForegroundColor Cyan
Write-Host "║          Command your microservices from the edge            ║" -ForegroundColor Cyan
Write-Host "║                                                               ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Function to check if command exists
function Test-Command {
    param($Command)
    $null = Get-Command $Command -ErrorAction SilentlyContinue
    return $?
}

# Step 1: Check prerequisites
Write-Host "=> Checking prerequisites..." -ForegroundColor Cyan

$MissingDeps = 0

if (-not (Test-Command docker)) {
    Write-Host "✗ Docker is not installed" -ForegroundColor Red
    Write-Host "  Install from: https://docs.docker.com/get-docker/" -ForegroundColor Yellow
    $MissingDeps = 1
} else {
    $dockerVersion = docker --version
    Write-Host "✓ Docker found: $dockerVersion" -ForegroundColor Green
}

if (-not (Test-Command docker-compose)) {
    Write-Host "✗ Docker Compose is not installed" -ForegroundColor Red
    Write-Host "  Install from: https://docs.docker.com/compose/install/" -ForegroundColor Yellow
    $MissingDeps = 1
} else {
    $composeVersion = docker-compose --version
    Write-Host "✓ Docker Compose found: $composeVersion" -ForegroundColor Green
}

if (-not (Test-Command go)) {
    Write-Host "⚠ Go is not installed (optional for local development)" -ForegroundColor Yellow
    Write-Host "  Install from: https://golang.org/doc/install" -ForegroundColor Yellow
} else {
    $goVersion = go version
    Write-Host "✓ Go found: $goVersion" -ForegroundColor Green
}

if (-not (Test-Command node)) {
    Write-Host "⚠ Node.js is not installed (optional for local development)" -ForegroundColor Yellow
    Write-Host "  Install from: https://nodejs.org/" -ForegroundColor Yellow
} else {
    $nodeVersion = node --version
    Write-Host "✓ Node.js found: $nodeVersion" -ForegroundColor Green
    
    if (-not (Test-Command npm)) {
        Write-Host "✗ npm is not installed" -ForegroundColor Red
        $MissingDeps = 1
    } else {
        $npmVersion = npm --version
        Write-Host "✓ npm found: $npmVersion" -ForegroundColor Green
    }
}

if ($MissingDeps -eq 1) {
    Write-Host "✗ Please install missing dependencies and run this script again" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Step 2: Setup environment files
Write-Host "=> Setting up environment files..." -ForegroundColor Cyan

# Backend .env
if (-not (Test-Path "backend\.env")) {
    Copy-Item "backend\.env.example" "backend\.env"
    Write-Host "✓ Created backend\.env from .env.example" -ForegroundColor Green
} else {
    Write-Host "⚠ backend\.env already exists, skipping..." -ForegroundColor Yellow
}

# Frontend .env.local
if (-not (Test-Path "frontend\.env.local")) {
    Copy-Item "frontend\.env.local.example" "frontend\.env.local"
    Write-Host "✓ Created frontend\.env.local from .env.local.example" -ForegroundColor Green
} else {
    Write-Host "⚠ frontend\.env.local already exists, skipping..." -ForegroundColor Yellow
}

Write-Host ""

# Step 3: Install Go dependencies
if (Test-Command go) {
    Write-Host "=> Installing Go dependencies..." -ForegroundColor Cyan
    Push-Location backend
    go mod download
    go mod tidy
    Write-Host "✓ Go dependencies installed" -ForegroundColor Green
    Pop-Location
    Write-Host ""
} else {
    Write-Host "⚠ Skipping Go dependencies (Go not installed)" -ForegroundColor Yellow
    Write-Host ""
}

# Step 4: Install Node dependencies
if (Test-Command npm) {
    Write-Host "=> Installing Node.js dependencies..." -ForegroundColor Cyan
    Push-Location frontend
    npm install
    Write-Host "✓ Node.js dependencies installed" -ForegroundColor Green
    Pop-Location
    Write-Host ""
} else {
    Write-Host "⚠ Skipping Node.js dependencies (npm not installed)" -ForegroundColor Yellow
    Write-Host ""
}

# Step 5: Check Docker daemon
Write-Host "=> Checking Docker daemon..." -ForegroundColor Cyan
try {
    docker info | Out-Null
    Write-Host "✓ Docker daemon is running" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker daemon is not running" -ForegroundColor Red
    Write-Host "  Please start Docker and run this script again" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# Step 6: Build Docker images (optional)
$build = Read-Host "Do you want to build Docker images now? (y/N)"
if ($build -eq "y" -or $build -eq "Y") {
    Write-Host "=> Building Docker images..." -ForegroundColor Cyan
    docker-compose build
    Write-Host "✓ Docker images built successfully" -ForegroundColor Green
    Write-Host ""
}

# Summary
Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║                    🎉 SETUP COMPLETE! 🎉                      ║" -ForegroundColor Green
Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""

Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host ""
Write-Host "1️⃣  Start all services:" -ForegroundColor White
Write-Host "   docker-compose up" -ForegroundColor Yellow
Write-Host ""
Write-Host "2️⃣  Or start in detached mode:" -ForegroundColor White
Write-Host "   docker-compose up -d" -ForegroundColor Yellow
Write-Host ""
Write-Host "3️⃣  Access the dashboard:" -ForegroundColor White
Write-Host "   http://localhost:3000" -ForegroundColor Yellow
Write-Host ""
Write-Host "4️⃣  Access the API:" -ForegroundColor White
Write-Host "   http://localhost:8080" -ForegroundColor Yellow
Write-Host ""
Write-Host "5️⃣  View logs:" -ForegroundColor White
Write-Host "   docker-compose logs -f" -ForegroundColor Yellow
Write-Host ""
Write-Host "6️⃣  Stop services:" -ForegroundColor White
Write-Host "   docker-compose down" -ForegroundColor Yellow
Write-Host ""

Write-Host "For local development:" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Backend:  cd backend && go run main.go" -ForegroundColor Yellow
Write-Host "  Frontend: cd frontend && npm run dev" -ForegroundColor Yellow
Write-Host ""

Write-Host "Documentation:" -ForegroundColor Cyan
Write-Host ""
Write-Host "  📖 Full Guide:      README.md" -ForegroundColor Yellow
Write-Host "  🚀 Quick Start:     QUICKSTART.md" -ForegroundColor Yellow
Write-Host "  🏗️  Architecture:    ARCHITECTURE.md" -ForegroundColor Yellow
Write-Host ""

Write-Host "Happy coding! ☁️" -ForegroundColor Green
Write-Host ""
