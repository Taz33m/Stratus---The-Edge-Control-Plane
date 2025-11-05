const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export interface Service {
  id: string
  name: string
  region: string
  image: string
  version: string
  status: 'running' | 'stopped' | 'error' | 'starting'
  uptime: number
  created_at: string
  updated_at: string
}

export interface CreateServiceRequest {
  name: string
  region: string
  image: string
  version: string
}

export interface UpdateServiceRequest {
  status?: 'running' | 'stopped' | 'error' | 'starting'
  version?: string
}

export interface ServiceMetrics {
  service_id: string
  timestamp: string
  cpu_usage: number
  memory_usage: number
  request_count: number
  error_rate: number
  p95_latency: number
}

export interface Metrics {
  cpu_usage: number
  memory_usage: number
  request_count: number
  error_rate: number
  p95_latency: number
}

export interface DeploymentLog {
  id: string
  service_id: string
  service_name: string
  action: string
  status: string
  message: string
  created_at: string
}

class ApiClient {
  private baseUrl: string

  constructor() {
    this.baseUrl = API_URL
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    })

    if (!response.ok) {
      throw new Error(`API error: ${response.statusText}`)
    }

    return response.json()
  }

  // Services
  async getServices(filters?: { region?: string; status?: string }): Promise<{ services: Service[] }> {
    const params = new URLSearchParams()
    if (filters?.region) params.set('region', filters.region)
    if (filters?.status) params.set('status', filters.status)
    
    const query = params.toString() ? `?${params}` : ''
    return this.request(`/api/v1/services${query}`)
  }

  async getService(id: string): Promise<Service> {
    return this.request(`/api/v1/services/${id}`)
  }

  async createService(data: CreateServiceRequest): Promise<Service> {
    return this.request('/api/v1/services', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateService(id: string, data: UpdateServiceRequest): Promise<Service> {
    return this.request(`/api/v1/services/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }

  async deleteService(id: string): Promise<{ message: string }> {
    return this.request(`/api/v1/services/${id}`, {
      method: 'DELETE',
    })
  }

  // Metrics
  async getMetrics(serviceId: string): Promise<{ metrics: ServiceMetrics[] }> {
    return this.request(`/api/v1/metrics/${serviceId}`)
  }

  async getAggregatedMetrics(): Promise<any> {
    return this.request('/api/v1/metrics/aggregated')
  }

  // Logs
  async getDeploymentLogs(serviceId?: string): Promise<{ logs: DeploymentLog[] }> {
    const query = serviceId ? `?service_id=${serviceId}` : ''
    return this.request(`/api/v1/logs/deployment${query}`)
  }
}

export const api = new ApiClient()
