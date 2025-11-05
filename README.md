<div align="center">

<img src="https://raw.githubusercontent.com/carbon-design-system/carbon/main/packages/icons/src/svg/32/cloud.svg" width="80" height="80" alt="Stratus Logo" style="margin-bottom: 20px;" />

# Stratus

### The Edge Control Plane

> **Command your microservices from the edge.**

Stratus is an open-source microservices control plane and dashboard designed to simulate how edge services are deployed, configured, and monitored at global scale. Inspired by Cloudflare's control systems, Stratus provides developers with a powerful visual interface to manage distributed services.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)
![Next.js](https://img.shields.io/badge/Next.js-14-black?logo=next.js)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)

### 🎥 Demo Video

[![Watch Demo Video](https://img.shields.io/badge/▶️_Watch_Demo-1.8MB-blue?style=for-the-badge)](https://github.com/Taz33m/Stratus---The-Edge-Control-Plane/raw/main/stratus.mp4)

> **Note:** Click the badge above to watch the full demo video (1.8MB MP4)

</div>

---

## Features

- **Service Registry** — Manage microservices across multiple regions
- **Real-Time Metrics** — Live CPU, memory, latency, and error rate monitoring
- **WebSocket Updates** — Sub-second status synchronization across all clients
- **Deployment Logs** — Complete audit trail of all service operations
- **Carbon Design System** — Enterprise-grade IBM UI components

## Architecture

**Frontend** → Next.js 14, TypeScript, Carbon Design  
**Backend** → Go 1.21, Gin, WebSocket Hub  
**Data** → PostgreSQL 16, Redis 7  
**Deploy** → Docker Compose

## Quick Start

```bash
git clone https://github.com/Taz33m/Stratus---The-Edge-Control-Plane.git
cd Stratus---The-Edge-Control-Plane
docker-compose up --build
```

Access dashboard at `http://localhost:3000`  
API available at `http://localhost:8080`

## API

**REST Endpoints**
- `GET/POST /api/v1/services` — Service CRUD operations
- `GET /api/v1/metrics/:id` — Real-time service metrics
- `GET /api/v1/logs/deployment` — Deployment audit trail

**WebSocket**
- `ws://localhost:8080/ws` — Live status updates

## Configuration

**Backend** (`backend/.env`)
```env
DATABASE_URL=postgres://stratus:stratus@localhost:5432/stratus
REDIS_URL=redis://localhost:6379
PORT=8080
```

**Frontend** (`frontend/.env.local`)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

## License

MIT License - see [LICENSE](LICENSE)

<div align="center">

**Built with ❤️ for the edge**

⭐ Star this repo if you found it helpful!

</div>
