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

<video src="./stratus.mp4" width="100%" controls>
  Your browser does not support the video tag.
</video>

*[Download Demo Video](./stratus.mp4)*

</div>

---

## 🌟 Features

### Core Capabilities

- **🎛️ Service Registry**: Add, update, or delete service nodes with region assignments
- **🚀 Deployment Control**: Start, stop, and restart simulated microservices
- **📊 Metrics Dashboard**: Real-time CPU/Memory usage, request latency, and error rates
- **🗺️ Global Visualization**: Service health across multiple regions (US-East, US-West, EU-West, APAC)
- **⚙️ Configuration Management**: JSON-based configs with version history
- **🔄 Real-time Updates**: WebSocket-powered live status updates
- **🔐 Production-Ready Architecture**: JWT authentication, CORS, graceful shutdown

---

## 🏗️ Architecture

```
┌────────────────────────────┐
│ Frontend (Next.js + TS)    │
│ • Dashboard UI             │
│ • Real-time metrics        │
└─────────────┬──────────────┘
              │ REST / WebSocket
              ▼
┌─────────────────────────────────────────┐
│ Backend (Go)                            │
│ • API Gateway (Gin)                     │
│ • Service Registry                      │
│ • WebSocket Hub                         │
│ • Metrics Simulator                     │
└──────────────┬──────────────────────────┘
               │
     ┌─────────┴─────────┐
     ▼                   ▼
┌──────────┐      ┌────────────┐
│PostgreSQL│      │   Redis    │
│ Services │      │  Metrics   │
│   Data   │      │   Cache    │
└──────────┘      └────────────┘
```

---

## 🛠️ Tech Stack

| Layer           | Technology                          | Purpose                          |
|-----------------|-------------------------------------|----------------------------------|
| **Frontend**    | Next.js 14, TypeScript, TailwindCSS | Interactive dashboard UI         |
| **Backend**     | Go 1.21 (Gin framework)             | API gateway & business logic     |
| **Database**    | PostgreSQL 16                       | Service registry & configurations|
| **Cache**       | Redis 7                             | Real-time metrics storage        |
| **Real-time**   | WebSocket                           | Live status updates              |
| **Deployment**  | Docker Compose                      | Service orchestration            |
| **Hosting**     | Cloudflare Pages (optional)         | Frontend deployment              |

---

## 🚀 Quick Start

### Prerequisites

- **Docker** & **Docker Compose** (v2.0+)
- **Go** 1.21+ (for local backend development)
- **Node.js** 18+ (for local frontend development)
- **Git**

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/stratus.git
cd stratus
```

### 2. Start with Docker Compose

The easiest way to run the entire stack:

```bash
docker-compose up --build
```

This will start:
- PostgreSQL on `localhost:5432`
- Redis on `localhost:6379`
- Backend API on `localhost:8080`
- Frontend UI on `localhost:3000`

### 3. Access the Dashboard

Open your browser and navigate to:

```
http://localhost:3000
```

You should see the Stratus dashboard with:
- ✅ **Connected** status (green indicator)
- Stats cards showing service metrics
- Empty services table (ready for you to create your first service!)

---

## 💻 Local Development

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run PostgreSQL and Redis
docker-compose up postgres redis -d

# Run the backend
go run main.go
```

Backend will be available at `http://localhost:8080`

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.local.example .env.local

# Run development server
npm run dev
```

Frontend will be available at `http://localhost:3000`

---

## 📡 API Endpoints

### Services

| Method | Endpoint              | Description                    |
|--------|----------------------|--------------------------------|
| GET    | `/api/v1/services`   | List all services              |
| POST   | `/api/v1/services`   | Create a new service           |
| GET    | `/api/v1/services/:id` | Get service by ID            |
| PATCH  | `/api/v1/services/:id` | Update service (start/stop)  |
| DELETE | `/api/v1/services/:id` | Delete a service             |

### Metrics

| Method | Endpoint                    | Description                |
|--------|----------------------------|----------------------------|
| GET    | `/api/v1/metrics/:id`      | Get service metrics        |
| GET    | `/api/v1/metrics/aggregated` | Get aggregated metrics   |

### Logs

| Method | Endpoint                     | Description              |
|--------|----------------------------|--------------------------|
| GET    | `/api/v1/logs/deployment`  | Get deployment logs      |

### WebSocket

```
ws://localhost:8080/ws
```

Receives real-time updates for:
- Service status changes
- Metrics updates
- Deployment logs

---

## 📊 Usage Examples

### Create a Service

```bash
curl -X POST http://localhost:8080/api/v1/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "edge-api",
    "region": "US-East",
    "image": "nginx:alpine",
    "version": "1.0.0"
  }'
```

### Start a Service

```bash
curl -X PATCH http://localhost:8080/api/v1/services/{service-id} \
  -H "Content-Type: application/json" \
  -d '{"status": "running"}'
```

### Get Service Metrics

```bash
curl http://localhost:8080/api/v1/metrics/{service-id}
```

---

## 🎨 UI Features

### Dashboard Components

1. **Header**: Real-time WebSocket connection status
2. **Stats Cards**: 
   - Total Services
   - Running Services
   - Active Regions
   - Average CPU Usage
3. **Service Table**: 
   - Service name & Docker image
   - Region with color indicators
   - Status with live updates
   - Version info
   - Uptime tracking
   - Quick actions (Start/Stop/Delete)
4. **Create Service Dialog**: Easy form to deploy new services

### Design Philosophy

- **Dark Mode First**: Futuristic, minimal aesthetic
- **Real-time Everything**: WebSocket-powered updates (<500ms latency)
- **Responsive**: Works on desktop, tablet, and mobile
- **Accessibility**: Keyboard navigation, screen reader support

---

## 🔧 Configuration

### Backend Configuration (`backend/.env`)

```env
PORT=8080
DATABASE_URL=postgres://stratus:stratus@localhost:5432/stratus?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key-change-in-production
CORS_ORIGINS=http://localhost:3000
ENVIRONMENT=development
```

### Frontend Configuration (`frontend/.env.local`)

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

---

## 🧪 Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

---

## 📦 Deployment

### Deploy Frontend to Cloudflare Pages

1. Build the static export:

```bash
cd frontend
npm run build
```

2. Deploy to Cloudflare Pages:

```bash
npx wrangler pages deploy out
```

### Deploy Backend

Use any cloud provider that supports Docker:

- **Fly.io**: `fly deploy`
- **Railway**: Connect GitHub repo
- **AWS ECS**: Use the provided Dockerfile
- **Google Cloud Run**: `gcloud run deploy`

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🎯 Roadmap

- [ ] **Authentication**: JWT-based user login
- [ ] **Role-Based Access Control**: Viewer/Admin roles
- [ ] **CI/CD Simulation**: Mock deployment pipelines
- [ ] **Advanced Metrics**: Custom dashboards with Grafana
- [ ] **Multi-Cluster Support**: Manage services across Kubernetes clusters
- [ ] **Service Mesh Integration**: Istio/Linkerd compatibility
- [ ] **Alert System**: Slack/Email notifications
- [ ] **GraphQL API**: Alternative to REST

---

## 💡 Inspiration

This project demonstrates:

- Cloudflare's architectural philosophy (UI + control plane at Internet scale)
- Microservices communication and orchestration
- Modern frontend and backend development practices
- Real-world software engineering skills

---

## 📧 Contact

**Tazeem Mahashin** - [GitHub](https://github.com/yourusername)

Project Link: [https://github.com/yourusername/stratus](https://github.com/yourusername/stratus)

---

<div align="center">

**Built with ❤️ for the edge**

⭐ Star this repo if you found it helpful!

</div>
