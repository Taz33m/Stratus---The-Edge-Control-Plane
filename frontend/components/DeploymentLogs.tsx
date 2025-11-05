'use client'

import { useEffect, useState } from 'react'
import { Tile, Tag } from '@carbon/react'
import { api, DeploymentLog } from '@/lib/api'
import { Checkmark, Error, InProgress } from '@carbon/icons-react'

export function DeploymentLogs() {
  const [logs, setLogs] = useState<DeploymentLog[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    fetchLogs()
    const interval = setInterval(fetchLogs, 5000) // Refresh every 5 seconds
    return () => clearInterval(interval)
  }, [])

  const fetchLogs = async () => {
    try {
      const data = await api.getDeploymentLogs()
      setLogs(data.logs || [])
      setIsLoading(false)
    } catch (error) {
      console.error('Failed to fetch logs:', error)
    }
  }

  const getActionIcon = (action: string) => {
    switch (action) {
      case 'create':
        return <Checkmark size={16} style={{ color: '#42be65' }} />
      case 'update':
        return <InProgress size={16} style={{ color: '#0f62fe' }} />
      case 'delete':
        return <Error size={16} style={{ color: '#da1e28' }} />
      default:
        return null
    }
  }

  const getActionColor = (action: string): any => {
    switch (action) {
      case 'create':
        return 'green'
      case 'update':
        return 'blue'
      case 'delete':
        return 'red'
      default:
        return 'gray'
    }
  }

  if (isLoading) {
    return (
      <Tile style={{ padding: '2rem', textAlign: 'center' }}>
        Loading deployment logs...
      </Tile>
    )
  }

  return (
    <div style={{ marginTop: '2rem' }}>
      <h3 style={{ marginBottom: '1rem', fontSize: '20px', fontWeight: '600' }}>
        Deployment Activity
      </h3>
      
      <Tile style={{ padding: 0 }}>
        <div style={{ maxHeight: '500px', overflowY: 'auto' }}>
          {logs.length === 0 ? (
            <div style={{ padding: '2rem', textAlign: 'center', color: '#c6c6c6' }}>
              No deployment logs yet. Deploy a service to see activity here.
            </div>
          ) : (
            logs.map((log) => (
              <div
                key={log.id}
                style={{
                  padding: '1rem',
                  borderBottom: '1px solid #393939',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '1rem',
                }}
              >
                <div style={{ flexShrink: 0 }}>
                  {getActionIcon(log.action)}
                </div>
                <div style={{ flex: 1 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
                    <span style={{ fontWeight: '600' }}>{log.service_name}</span>
                    <Tag type={getActionColor(log.action)} size="sm">
                      {log.action.toUpperCase()}
                    </Tag>
                  </div>
                  <div style={{ fontSize: '13px', color: '#c6c6c6' }}>
                    {log.message}
                  </div>
                  <div style={{ fontSize: '12px', color: '#8d8d8d', marginTop: '4px' }}>
                    {new Date(log.created_at).toLocaleString()}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </Tile>
    </div>
  )
}
