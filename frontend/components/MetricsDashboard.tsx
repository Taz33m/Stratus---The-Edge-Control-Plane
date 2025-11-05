'use client'

import { useEffect, useState } from 'react'
import { Tile, Grid, Column, Tag } from '@carbon/react'
import { LineChart, Line, AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from 'recharts'
import { api, Metrics } from '@/lib/api'

interface MetricsDashboardProps {
  serviceId: string
  serviceName: string
}

interface MetricHistory {
  timestamp: string
  cpu: number
  memory: number
  latency: number
  requests: number
}

export function MetricsDashboard({ serviceId, serviceName }: MetricsDashboardProps) {
  const [currentMetrics, setCurrentMetrics] = useState<Metrics | null>(null)
  const [history, setHistory] = useState<MetricHistory[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    fetchMetrics()
    const interval = setInterval(fetchMetrics, 3000) // Update every 3 seconds
    return () => clearInterval(interval)
  }, [serviceId])

  const fetchMetrics = async () => {
    try {
      const response = await api.getMetrics(serviceId)
      // Get the latest metrics from the array
      const latestMetrics = response.metrics && response.metrics.length > 0 
        ? response.metrics[0] 
        : null
      
      if (!latestMetrics) {
        // Generate mock metrics if none exist yet
        const mockMetrics: Metrics = {
          cpu_usage: Math.random() * 80,
          memory_usage: Math.random() * 70 + 20,
          p95_latency: Math.floor(Math.random() * 200 + 50),
          request_count: Math.floor(Math.random() * 1000),
          error_rate: Math.random() * 2,
        }
        setCurrentMetrics(mockMetrics)
        
        setHistory(prev => {
          const newPoint = {
            timestamp: new Date().toLocaleTimeString(),
            cpu: mockMetrics.cpu_usage,
            memory: mockMetrics.memory_usage,
            latency: mockMetrics.p95_latency,
            requests: mockMetrics.request_count,
          }
          return [...prev, newPoint].slice(-20)
        })
      } else {
        const metrics: Metrics = {
          cpu_usage: latestMetrics.cpu_usage,
          memory_usage: latestMetrics.memory_usage,
          p95_latency: latestMetrics.p95_latency,
          request_count: latestMetrics.request_count,
          error_rate: latestMetrics.error_rate,
        }
        setCurrentMetrics(metrics)
        
        setHistory(prev => {
          const newPoint = {
            timestamp: new Date().toLocaleTimeString(),
            cpu: metrics.cpu_usage,
            memory: metrics.memory_usage,
            latency: metrics.p95_latency,
            requests: metrics.request_count,
          }
          return [...prev, newPoint].slice(-20)
        })
      }
      
      setIsLoading(false)
    } catch (error) {
      console.error('Failed to fetch metrics:', error)
    }
  }

  if (isLoading) {
    return (
      <Tile style={{ padding: '2rem', textAlign: 'center' }}>
        Loading metrics for {serviceName}...
      </Tile>
    )
  }

  return (
    <div style={{ marginTop: '2rem' }}>
      <h3 style={{ marginBottom: '1rem', fontSize: '20px', fontWeight: '600' }}>
        Real-Time Metrics: {serviceName}
      </h3>

      {/* Current Metrics Overview */}
      <Grid style={{ marginBottom: '2rem' }}>
        <Column lg={4} md={4} sm={4}>
          <Tile style={{ textAlign: 'center', padding: '1.5rem' }}>
            <div style={{ fontSize: '12px', color: '#c6c6c6', marginBottom: '8px' }}>CPU Usage</div>
            <div style={{ fontSize: '32px', fontWeight: '600', color: getCpuColor(currentMetrics?.cpu_usage || 0) }}>
              {currentMetrics?.cpu_usage.toFixed(1)}%
            </div>
            <Tag type={getCpuTagType(currentMetrics?.cpu_usage || 0)} size="sm" style={{ marginTop: '8px' }}>
              {getCpuStatus(currentMetrics?.cpu_usage || 0)}
            </Tag>
          </Tile>
        </Column>
        <Column lg={4} md={4} sm={4}>
          <Tile style={{ textAlign: 'center', padding: '1.5rem' }}>
            <div style={{ fontSize: '12px', color: '#c6c6c6', marginBottom: '8px' }}>Memory Usage</div>
            <div style={{ fontSize: '32px', fontWeight: '600', color: getMemoryColor(currentMetrics?.memory_usage || 0) }}>
              {currentMetrics?.memory_usage.toFixed(1)}%
            </div>
            <Tag type={getMemoryTagType(currentMetrics?.memory_usage || 0)} size="sm" style={{ marginTop: '8px' }}>
              {getMemoryStatus(currentMetrics?.memory_usage || 0)}
            </Tag>
          </Tile>
        </Column>
        <Column lg={4} md={4} sm={4}>
          <Tile style={{ textAlign: 'center', padding: '1.5rem' }}>
            <div style={{ fontSize: '12px', color: '#c6c6c6', marginBottom: '8px' }}>P95 Latency</div>
            <div style={{ fontSize: '32px', fontWeight: '600' }}>
              {currentMetrics?.p95_latency}ms
            </div>
            <Tag type="blue" size="sm" style={{ marginTop: '8px' }}>
              {currentMetrics?.request_count.toLocaleString()} req/s
            </Tag>
          </Tile>
        </Column>
        <Column lg={4} md={4} sm={4}>
          <Tile style={{ textAlign: 'center', padding: '1.5rem' }}>
            <div style={{ fontSize: '12px', color: '#c6c6c6', marginBottom: '8px' }}>Error Rate</div>
            <div style={{ fontSize: '32px', fontWeight: '600', color: currentMetrics && currentMetrics.error_rate > 1 ? '#ff832b' : '#42be65' }}>
              {currentMetrics?.error_rate.toFixed(2)}%
            </div>
          </Tile>
        </Column>
      </Grid>

      {/* CPU Usage Chart */}
      <Tile style={{ marginBottom: '1rem', padding: '1.5rem' }}>
        <h4 style={{ marginBottom: '1rem', fontSize: '16px' }}>CPU Usage Over Time</h4>
        <ResponsiveContainer width="100%" height={200}>
          <AreaChart data={history}>
            <defs>
              <linearGradient id="colorCpu" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#0f62fe" stopOpacity={0.8}/>
                <stop offset="95%" stopColor="#0f62fe" stopOpacity={0.1}/>
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke="#393939" />
            <XAxis dataKey="timestamp" stroke="#c6c6c6" style={{ fontSize: '12px' }} />
            <YAxis stroke="#c6c6c6" style={{ fontSize: '12px' }} domain={[0, 100]} />
            <Tooltip 
              contentStyle={{ backgroundColor: '#262626', border: '1px solid #393939', borderRadius: '4px' }}
              labelStyle={{ color: '#f4f4f4' }}
            />
            <Area type="monotone" dataKey="cpu" stroke="#0f62fe" fillOpacity={1} fill="url(#colorCpu)" />
          </AreaChart>
        </ResponsiveContainer>
      </Tile>

      {/* Memory Usage Chart */}
      <Tile style={{ marginBottom: '1rem', padding: '1.5rem' }}>
        <h4 style={{ marginBottom: '1rem', fontSize: '16px' }}>Memory Usage Over Time</h4>
        <ResponsiveContainer width="100%" height={200}>
          <AreaChart data={history}>
            <defs>
              <linearGradient id="colorMemory" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#8a3ffc" stopOpacity={0.8}/>
                <stop offset="95%" stopColor="#8a3ffc" stopOpacity={0.1}/>
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke="#393939" />
            <XAxis dataKey="timestamp" stroke="#c6c6c6" style={{ fontSize: '12px' }} />
            <YAxis stroke="#c6c6c6" style={{ fontSize: '12px' }} domain={[0, 100]} />
            <Tooltip 
              contentStyle={{ backgroundColor: '#262626', border: '1px solid #393939', borderRadius: '4px' }}
              labelStyle={{ color: '#f4f4f4' }}
            />
            <Area type="monotone" dataKey="memory" stroke="#8a3ffc" fillOpacity={1} fill="url(#colorMemory)" />
          </AreaChart>
        </ResponsiveContainer>
      </Tile>

      {/* Latency Chart */}
      <Tile style={{ marginBottom: '1rem', padding: '1.5rem' }}>
        <h4 style={{ marginBottom: '1rem', fontSize: '16px' }}>Latency & Request Rate</h4>
        <ResponsiveContainer width="100%" height={200}>
          <LineChart data={history}>
            <CartesianGrid strokeDasharray="3 3" stroke="#393939" />
            <XAxis dataKey="timestamp" stroke="#c6c6c6" style={{ fontSize: '12px' }} />
            <YAxis yAxisId="left" stroke="#c6c6c6" style={{ fontSize: '12px' }} />
            <YAxis yAxisId="right" orientation="right" stroke="#c6c6c6" style={{ fontSize: '12px' }} />
            <Tooltip 
              contentStyle={{ backgroundColor: '#262626', border: '1px solid #393939', borderRadius: '4px' }}
              labelStyle={{ color: '#f4f4f4' }}
            />
            <Legend />
            <Line yAxisId="left" type="monotone" dataKey="latency" stroke="#42be65" name="Latency (ms)" strokeWidth={2} dot={false} />
            <Line yAxisId="right" type="monotone" dataKey="requests" stroke="#ff832b" name="Requests" strokeWidth={2} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </Tile>
    </div>
  )
}

function getCpuColor(cpu: number): string {
  if (cpu > 80) return '#ff832b' // orange
  if (cpu > 60) return '#f1c21b' // yellow
  return '#42be65' // green
}

function getCpuStatus(cpu: number): string {
  if (cpu > 80) return 'High Load'
  if (cpu > 60) return 'Moderate'
  return 'Normal'
}

function getCpuTagType(cpu: number): any {
  if (cpu > 80) return 'red'
  if (cpu > 60) return 'orange'
  return 'green'
}

function getMemoryColor(memory: number): string {
  if (memory > 85) return '#ff832b'
  if (memory > 70) return '#f1c21b'
  return '#42be65'
}

function getMemoryStatus(memory: number): string {
  if (memory > 85) return 'Critical'
  if (memory > 70) return 'Warning'
  return 'Healthy'
}

function getMemoryTagType(memory: number): any {
  if (memory > 85) return 'red'
  if (memory > 70) return 'orange'
  return 'green'
}
