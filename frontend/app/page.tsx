'use client'

import { useEffect, useState } from 'react'
import {
  Header,
  HeaderName,
  HeaderGlobalBar,
  HeaderGlobalAction,
  Content,
  Button,
  DataTable,
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableHeader,
  TableBody,
  TableCell,
  TableToolbar,
  TableToolbarContent,
  TableToolbarSearch,
  Tag,
  Modal,
  TextInput,
  Select,
  SelectItem,
  Tile,
  InlineLoading,
  Grid,
  Column,
} from '@carbon/react'
import {
  Add,
  Play,
  Stop,
  TrashCan,
  Cloud,
  Dashboard as DashboardIcon,
  Renew,
  ChartLine,
} from '@carbon/icons-react'
import { api, Service } from '@/lib/api'
import { useWebSocket } from '@/hooks/useWebSocket'
import { MetricsDashboard } from '@/components/MetricsDashboard'
import { DeploymentLogs } from '@/components/DeploymentLogs'

export default function Dashboard() {
  const [services, setServices] = useState<Service[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedServiceForMetrics, setSelectedServiceForMetrics] = useState<Service | null>(null)
  const [showMetrics, setShowMetrics] = useState(false)
  const [showLogs, setShowLogs] = useState(false)
  const [newService, setNewService] = useState({
    name: '',
    region: 'US-East',
    image: '',
    version: '',
  })
  const { isConnected, lastMessage } = useWebSocket()

  useEffect(() => {
    fetchServices()
  }, [])

  useEffect(() => {
    if (lastMessage?.type === 'service_update') {
      // Refresh services when an update is received
      fetchServices()
    }
  }, [lastMessage])

  const fetchServices = async () => {
    try {
      const data = await api.getServices()
      setServices(data.services || [])
    } catch (error) {
      console.error('Failed to fetch services:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateService = async (data: {
    name: string
    region: string
    image: string
    version: string
  }) => {
    try {
      await api.createService(data)
      await fetchServices()
    } catch (error) {
      console.error('Failed to create service:', error)
    }
  }

  const handleStartService = async (id: string) => {
    try {
      await api.updateService(id, { status: 'running' })
      await fetchServices()
    } catch (error) {
      console.error('Failed to start service:', error)
    }
  }

  const handleStopService = async (id: string) => {
    try {
      await api.updateService(id, { status: 'stopped' })
      await fetchServices()
    } catch (error) {
      console.error('Failed to stop service:', error)
    }
  }

  const handleDeleteService = async (id: string) => {
    if (!confirm('Are you sure you want to delete this service?')) return

    try {
      await api.deleteService(id)
      await fetchServices()
    } catch (error) {
      console.error('Failed to delete service:', error)
    }
  }

  const handleCreateSubmit = async () => {
    await handleCreateService(newService)
    setIsModalOpen(false)
    setNewService({ name: '', region: 'US-East', image: '', version: '' })
  }

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { type: any; label: string }> = {
      running: { type: 'green', label: 'Running' },
      stopped: { type: 'gray', label: 'Stopped' },
      error: { type: 'red', label: 'Error' },
      starting: { type: 'blue', label: 'Starting' },
    }
    const config = statusMap[status] || { type: 'gray', label: status }
    return <Tag type={config.type}>{config.label}</Tag>
  }

  const getRegionTag = (region: string) => {
    const colors: Record<string, any> = {
      'US-East': 'blue',
      'US-West': 'purple',
      'EU-West': 'teal',
      'APAC': 'magenta',
    }
    return <Tag type={colors[region] || 'gray'}>{region}</Tag>
  }

  const totalServices = services.length
  const runningServices = services.filter(s => s.status === 'running').length
  const regions = new Set(services.map(s => s.region)).size

  return (
    <>
      <Header aria-label="Stratus Control Plane">
        <HeaderName href="#" prefix="">
          <Cloud size={20} style={{ marginRight: '8px' }} />
          Stratus
        </HeaderName>
        <HeaderGlobalBar>
          <HeaderGlobalAction 
            aria-label="Metrics Dashboard" 
            onClick={() => setShowMetrics(!showMetrics)}
            isActive={showMetrics}
          >
            <ChartLine size={20} />
          </HeaderGlobalAction>
          <HeaderGlobalAction 
            aria-label="Deployment Logs" 
            onClick={() => setShowLogs(!showLogs)}
            isActive={showLogs}
          >
            <DashboardIcon size={20} />
          </HeaderGlobalAction>
          <HeaderGlobalAction
            aria-label={isConnected ? 'Connected' : 'Disconnected'}
            isActive={isConnected}
            tooltipAlignment="end"
          >
            <Renew size={20} />
          </HeaderGlobalAction>
        </HeaderGlobalBar>
      </Header>

      <Content>
        <Grid fullWidth style={{ padding: '2rem 0' }}>
          {/* Stats Cards */}
          <Column lg={16}>
            <Grid>
              <Column lg={5} md={4} sm={4}>
                <Tile style={{ marginBottom: '1rem', height: '120px' }}>
                  <div style={{ fontSize: '14px', color: '#c6c6c6', marginBottom: '8px' }}>
                    Total Services
                  </div>
                  <div style={{ fontSize: '32px', fontWeight: '600' }}>{totalServices}</div>
                </Tile>
              </Column>
              <Column lg={5} md={4} sm={4}>
                <Tile style={{ marginBottom: '1rem', height: '120px' }}>
                  <div style={{ fontSize: '14px', color: '#c6c6c6', marginBottom: '8px' }}>
                    Running Services
                  </div>
                  <div style={{ fontSize: '32px', fontWeight: '600', color: '#42be65' }}>
                    {runningServices}
                  </div>
                </Tile>
              </Column>
              <Column lg={6} md={4} sm={4}>
                <Tile style={{ marginBottom: '1rem', height: '120px' }}>
                  <div style={{ fontSize: '14px', color: '#c6c6c6', marginBottom: '8px' }}>
                    Active Regions
                  </div>
                  <div style={{ fontSize: '32px', fontWeight: '600' }}>{regions}</div>
                </Tile>
              </Column>
            </Grid>
          </Column>

          {/* Services Table */}
          <Column lg={16}>
            <DataTable
              rows={services.map((s) => ({ ...s, id: s.id }))}
              headers={[
                { key: 'name', header: 'Service Name' },
                { key: 'status', header: 'Status' },
                { key: 'region', header: 'Region' },
                { key: 'image', header: 'Image' },
                { key: 'version', header: 'Version' },
                { key: 'actions', header: 'Actions' },
              ]}
            >
              {({ rows, headers, getTableProps, getHeaderProps, getRowProps }) => (
                <TableContainer
                  title="Services"
                  description="Manage your distributed microservices"
                >
                  <TableToolbar>
                    <TableToolbarContent>
                      <TableToolbarSearch persistent />
                      <Button
                        renderIcon={Add}
                        onClick={() => setIsModalOpen(true)}
                      >
                        Create Service
                      </Button>
                    </TableToolbarContent>
                  </TableToolbar>
                  <Table {...getTableProps()}>
                    <TableHead>
                      <TableRow>
                        {headers.map((header) => (
                          <TableHeader {...getHeaderProps({ header })}>
                            {header.header}
                          </TableHeader>
                        ))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {isLoading ? (
                        <TableRow>
                          <TableCell colSpan={6}>
                            <InlineLoading description="Loading services..." />
                          </TableCell>
                        </TableRow>
                      ) : rows.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={6}>
                            <div style={{ textAlign: 'center', padding: '2rem' }}>
                              No services yet. Create one to get started!
                            </div>
                          </TableCell>
                        </TableRow>
                      ) : (
                        rows.map((row) => {
                          const service = services.find((s) => s.id === row.id)
                          return (
                            <TableRow {...getRowProps({ row })}>
                              <TableCell>{service?.name}</TableCell>
                              <TableCell>{getStatusTag(service?.status || '')}</TableCell>
                              <TableCell>{getRegionTag(service?.region || '')}</TableCell>
                              <TableCell>{service?.image}</TableCell>
                              <TableCell>{service?.version}</TableCell>
                              <TableCell>
                                <div style={{ display: 'flex', gap: '8px' }}>
                                  <Button
                                    kind="ghost"
                                    size="sm"
                                    renderIcon={ChartLine}
                                    iconDescription="View Metrics"
                                    hasIconOnly
                                    onClick={() => {
                                      setSelectedServiceForMetrics(service || null)
                                      setShowMetrics(true)
                                      setShowLogs(false)
                                    }}
                                  />
                                  {service?.status !== 'running' ? (
                                    <Button
                                      kind="ghost"
                                      size="sm"
                                      renderIcon={Play}
                                      iconDescription="Start"
                                      hasIconOnly
                                      onClick={() => handleStartService(service?.id || '')}
                                    />
                                  ) : (
                                    <Button
                                      kind="ghost"
                                      size="sm"
                                      renderIcon={Stop}
                                      iconDescription="Stop"
                                      hasIconOnly
                                      onClick={() => handleStopService(service?.id || '')}
                                    />
                                  )}
                                  <Button
                                    kind="danger--ghost"
                                    size="sm"
                                    renderIcon={TrashCan}
                                    iconDescription="Delete"
                                    hasIconOnly
                                    onClick={() => handleDeleteService(service?.id || '')}
                                  />
                                </div>
                              </TableCell>
                            </TableRow>
                          )
                        })
                      )}
                    </TableBody>
                  </Table>
                </TableContainer>
              )}
            </DataTable>
          </Column>

          {/* Metrics Dashboard */}
          {showMetrics && selectedServiceForMetrics && (
            <Column lg={16}>
              <MetricsDashboard 
                serviceId={selectedServiceForMetrics.id} 
                serviceName={selectedServiceForMetrics.name}
              />
            </Column>
          )}

          {/* Deployment Logs */}
          {showLogs && (
            <Column lg={16}>
              <DeploymentLogs />
            </Column>
          )}
        </Grid>
      </Content>

      {/* Create Service Modal */}
      <Modal
        open={isModalOpen}
        modalHeading="Create New Service"
        primaryButtonText="Create"
        secondaryButtonText="Cancel"
        onRequestClose={() => setIsModalOpen(false)}
        onRequestSubmit={handleCreateSubmit}
      >
        <div style={{ marginBottom: '1rem' }}>
          <TextInput
            id="service-name"
            labelText="Service Name"
            placeholder="e.g., edge-api"
            value={newService.name}
            onChange={(e) => setNewService({ ...newService, name: e.target.value })}
          />
        </div>
        <div style={{ marginBottom: '1rem' }}>
          <Select
            id="region"
            labelText="Region"
            value={newService.region}
            onChange={(e) => setNewService({ ...newService, region: e.target.value })}
          >
            <SelectItem value="US-East" text="US-East" />
            <SelectItem value="US-West" text="US-West" />
            <SelectItem value="EU-West" text="EU-West" />
            <SelectItem value="APAC" text="APAC" />
          </Select>
        </div>
        <div style={{ marginBottom: '1rem' }}>
          <TextInput
            id="image"
            labelText="Docker Image"
            placeholder="e.g., nginx:alpine"
            value={newService.image}
            onChange={(e) => setNewService({ ...newService, image: e.target.value })}
          />
        </div>
        <div style={{ marginBottom: '1rem' }}>
          <TextInput
            id="version"
            labelText="Version"
            placeholder="e.g., 1.0.0"
            value={newService.version}
            onChange={(e) => setNewService({ ...newService, version: e.target.value })}
          />
        </div>
      </Modal>
    </>
  )
}
