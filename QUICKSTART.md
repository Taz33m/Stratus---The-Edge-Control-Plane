# 🚀 Stratus Quick Start Guide

Get Stratus running in under 5 minutes!

## Prerequisites

- **Docker** and **Docker Compose** installed
- Ports 3000, 8080, 5432, 6379 available

## Step 1: Clone & Start

```bash
# Clone the repository
git clone https://github.com/yourusername/stratus.git
cd stratus

# Start all services
docker-compose up --build
```

Wait for services to start (about 1-2 minutes on first run).

## Step 2: Access the Dashboard

Open your browser to:
```
http://localhost:3000
```

You should see:
- ✅ Green "Connected" status in the header
- Four stats cards (all showing 0 initially)
- An empty services table

## Step 3: Create Your First Service

### Option A: Using the UI

1. Click **"Create Service"** button
2. Fill in the form:
   - **Service Name**: `my-first-service`
   - **Region**: `US-East`
   - **Docker Image**: `nginx:alpine`
   - **Version**: `1.0.0`
3. Click **"Create"**

The service will appear in the table instantly!

### Option B: Using the API

```bash
curl -X POST http://localhost:8080/api/v1/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-first-service",
    "region": "US-East",
    "image": "nginx:alpine",
    "version": "1.0.0"
  }'
```

## Step 4: Start the Service

Click the **"Start"** button next to your service.

Watch the status change from **stopped** → **running** in real-time!

## Step 5: Explore

Try these actions:

- **Create more services** in different regions (US-West, EU-West, APAC)
- **Stop/Start services** to see status updates
- **Delete services** you no longer need
- **Watch the stats cards** update automatically

## API Endpoints

Test the API directly:

```bash
# List all services
curl http://localhost:8080/api/v1/services

# Get a specific service
curl http://localhost:8080/api/v1/services/{service-id}

# Update service status
curl -X PATCH http://localhost:8080/api/v1/services/{service-id} \
  -H "Content-Type: application/json" \
  -d '{"status": "running"}'

# Delete a service
curl -X DELETE http://localhost:8080/api/v1/services/{service-id}

# Get deployment logs
curl http://localhost:8080/api/v1/logs/deployment
```

## WebSocket Connection

Connect to real-time updates:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
  // { type: 'service_update', payload: {...} }
};
```

## Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

## Troubleshooting

### Port already in use

```bash
# Find what's using the port
lsof -i :3000  # or :8080, :5432, :6379

# Kill the process or change ports in docker-compose.yml
```

### Frontend not connecting to backend

Check the environment variables:
```bash
# frontend/.env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

### Database connection errors

```bash
# Restart just the database
docker-compose restart postgres

# Check database logs
docker-compose logs postgres
```

---

**Enjoy building with Stratus!** ☁️
