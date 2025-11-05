# Stratus Architecture

## Overview

Stratus is built as a modern microservices control plane with a clear separation between frontend, backend, and data layers. The architecture emphasizes real-time communication, scalability, and developer experience.

## System Architecture

```
                            ┌─────────────────────────────┐
                            │   Users / Administrators    │
                            └──────────────┬──────────────┘
                                          │
                            ┌─────────────┴──────────────┐
                            │  Next.js Frontend (Port 3000) │
                            │  - Dashboard UI                │
                            │  - Real-time updates           │
                            │  - TailwindCSS styling         │
                            └──────────────┬─────────────────┘
                                          │
                     ┌───────────────────┼───────────────────┐
                     │                   │                   │
                  REST API          WebSocket           Static Assets
                     │                   │                   │
                     ▼                   ▼                   │
    ┌────────────────────────────────────────────────┐      │
    │       Go Backend API (Port 8080)               │      │
    │  ┌──────────────────────────────────────────┐  │      │
    │  │           Gin HTTP Server                │  │      │
    │  └────────────────┬─────────────────────────┘  │      │
    │                   │                             │      │
    │  ┌────────────────┼──────────────────────────┐ │      │
    │  │                │                          │ │      │
    │  │  ┌─────────────▼────────┐  ┌─────────────▼─┐│     │
    │  │  │   REST Router        │  │  WebSocket    ││     │
    │  │  │  - Services API      │  │     Hub       ││     │
    │  │  │  - Metrics API       │  │  - Broadcast  ││     │
    │  │  │  - Logs API          │  │  - Clients    ││     │
    │  │  └─────────────┬────────┘  └───────────────┘│     │
    │  │                │                             │     │
    │  │  ┌─────────────▼────────────────────────┐   │     │
    │  │  │        Business Logic Layer          │   │     │
    │  │  │  - ServiceHandler                    │   │     │
    │  │  │  - MetricsHandler                    │   │     │
    │  │  │  - LogsHandler                       │   │     │
    │  │  └─────────────┬────────────────────────┘   │     │
    │  │                │                             │     │
    │  └────────────────┼─────────────────────────────┘     │
    │                   │                                   │
    └───────────────────┼───────────────────────────────────┘
                        │
          ┌─────────────┼─────────────┐
          │             │             │
          ▼             ▼             ▼
  ┌───────────┐  ┌──────────┐  ┌──────────┐
  │PostgreSQL │  │  Redis   │  │  Future  │
  │  (5432)   │  │  (6379)  │  │Services  │
  │           │  │          │  │          │
  │ Services  │  │ Metrics  │  │  Auth    │
  │ Configs   │  │  Cache   │  │  etc.    │
  │   Logs    │  │          │  │          │
  └───────────┘  └──────────┘  └──────────┘
```

## Component Details

### Frontend Layer (Next.js + TypeScript)

**Technology**: Next.js 14, React 18, TypeScript 5, TailwindCSS 3

**Key Files**:
- `app/page.tsx` - Main dashboard
- `components/` - Reusable UI components
- `lib/api.ts` - API client
- `hooks/useWebSocket.ts` - Real-time connection

**Responsibilities**:
- Render interactive dashboard
- Manage WebSocket connections
- Display real-time metrics
- Handle user interactions

**Data Flow**:
1. User triggers action (create service, start/stop)
2. API client sends HTTP request to backend
3. WebSocket receives real-time updates
4. React state updates, UI re-renders

### Backend Layer (Go + Gin)

**Technology**: Go 1.21, Gin framework, Gorilla WebSocket

**Project Structure**:
```
backend/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # DB connections & migrations
│   ├── handlers/          # HTTP request handlers
│   ├── models/            # Data models
│   ├── router/            # Route definitions
│   └── websocket/         # WebSocket hub & clients
```

**Key Components**:

1. **HTTP Server** (`main.go`)
   - Initializes Gin router
   - Sets up middleware (CORS, logging)
   - Graceful shutdown handling

2. **Service Handler** (`handlers/service_handler.go`)
   - CRUD operations for services
   - Status management (start/stop)
   - Deployment log creation

3. **WebSocket Hub** (`websocket/hub.go`)
   - Manages connected clients
   - Broadcasts messages to all clients
   - Handles client registration/unregistration

4. **Database Layer** (`database/`)
   - PostgreSQL for persistent data
   - Redis for real-time metrics cache
   - Automatic migrations on startup

### Data Layer

**PostgreSQL Schema**:

```sql
-- Services table
services (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  region VARCHAR(50) NOT NULL,
  image VARCHAR(255) NOT NULL,
  version VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'stopped',
  uptime BIGINT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

-- Service configurations
service_configs (
  id VARCHAR(36) PRIMARY KEY,
  service_id VARCHAR(36) NOT NULL REFERENCES services(id),
  config JSONB NOT NULL,
  version INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255)
)

-- Deployment logs
deployment_logs (
  id VARCHAR(36) PRIMARY KEY,
  service_id VARCHAR(36) NOT NULL REFERENCES services(id),
  action VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL,
  message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

**Redis Data Structure**:
```
metrics:{service_id} -> List of ServiceMetrics JSON objects
```

## Communication Patterns

### HTTP REST API

**Request Flow**:
```
Client → Gin Router → Handler → Database → Response
```

**Example**:
```go
// Create Service
POST /api/v1/services
Body: { "name": "edge-api", "region": "US-East", ... }
Response: { "id": "uuid", "status": "stopped", ... }
```

### WebSocket Real-time Updates

**Connection Flow**:
```
1. Client connects to ws://localhost:8080/ws
2. Hub registers new client
3. Backend broadcasts events:
   - service_update: Service status changes
   - metrics: Real-time metrics (every 5s)
   - log: Deployment logs
4. Client receives and processes messages
5. UI updates automatically
```

**Message Format**:
```json
{
  "type": "service_update",
  "payload": {
    "id": "service-uuid",
    "status": "running",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## Scalability Considerations

### Current Architecture (Single Instance)

- **Frontend**: Stateless, can scale horizontally
- **Backend**: Single instance, WebSocket connections
- **Database**: PostgreSQL (primary), Redis (cache)

### Future Scaling Path

1. **Horizontal Backend Scaling**
   - Add Redis Pub/Sub for WebSocket message distribution
   - Use sticky sessions or Redis for shared state
   - Deploy multiple backend instances behind load balancer

2. **Database Scaling**
   - PostgreSQL read replicas
   - Redis Cluster for distributed cache
   - Connection pooling (PgBouncer)

3. **CDN Integration**
   - Serve frontend from Cloudflare Pages
   - Cache static assets
   - Edge computing for API responses

## Security Layers

1. **CORS Protection**: Configurable allowed origins
2. **Input Validation**: Request body validation with Gin bindings
3. **SQL Injection Prevention**: Parameterized queries
4. **Graceful Shutdown**: Ensures data integrity on restart

**Future**:
- JWT authentication
- Rate limiting
- API key management
- Role-based access control (RBAC)

## Performance Metrics

| Metric | Target | Current |
|--------|--------|---------|
| API Response Time | < 100ms | ~50ms |
| WebSocket Latency | < 500ms | ~200ms |
| Database Query Time | < 50ms | ~20ms |
| Frontend Load Time | < 2s | ~1.5s |

## Deployment Architecture

### Development
```
Docker Compose → Local containers (Postgres, Redis, Backend, Frontend)
```

### Production (Recommended)
```
Cloudflare Pages (Frontend)
       ↓
AWS ALB / Cloud Load Balancer
       ↓
ECS / Cloud Run (Backend instances)
       ↓
RDS PostgreSQL + ElastiCache Redis
```

## Monitoring & Observability

**Planned**:
- Prometheus metrics export
- Grafana dashboards
- Structured logging (JSON)
- Distributed tracing (OpenTelemetry)
- Health check endpoints

## Technology Decisions

| Decision | Rationale |
|----------|-----------|
| **Go for Backend** | Performance, concurrency, low latency |
| **Gin Framework** | Lightweight, fast, good middleware support |
| **Next.js** | SSR/SSG capabilities, great DX, Vercel/CF deployment |
| **PostgreSQL** | ACID compliance, JSONB support, mature |
| **Redis** | Fast caching, Pub/Sub for real-time features |
| **WebSocket** | True bidirectional real-time communication |
| **Docker Compose** | Simple local orchestration, reproducible env |

---

**Last Updated**: Nov 4, 2024
**Version**: 1.0.0
