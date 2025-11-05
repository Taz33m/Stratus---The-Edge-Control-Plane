'use client'

import { Service } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Server, Activity, Cpu, Globe } from 'lucide-react'

interface StatsCardsProps {
  services: Service[]
}

export function StatsCards({ services }: StatsCardsProps) {
  const runningServices = services.filter((s) => s.status === 'running').length
  const totalServices = services.length
  const regions = new Set(services.map((s) => s.region)).size

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card className="bg-gradient-to-br from-blue-500/10 to-blue-600/10 border-blue-500/20">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Total Services</CardTitle>
          <Server className="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">{totalServices}</div>
          <p className="text-xs text-muted-foreground mt-1">Deployed instances</p>
        </CardContent>
      </Card>

      <Card className="bg-gradient-to-br from-green-500/10 to-green-600/10 border-green-500/20">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Running</CardTitle>
          <Activity className="h-4 w-4 text-green-500 animate-pulse" />
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">{runningServices}</div>
          <p className="text-xs text-muted-foreground mt-1">Active services</p>
        </CardContent>
      </Card>

      <Card className="bg-gradient-to-br from-purple-500/10 to-purple-600/10 border-purple-500/20">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Regions</CardTitle>
          <Globe className="h-4 w-4 text-purple-500" />
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">{regions}</div>
          <p className="text-xs text-muted-foreground mt-1">Global coverage</p>
        </CardContent>
      </Card>

      <Card className="bg-gradient-to-br from-orange-500/10 to-orange-600/10 border-orange-500/20">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium">Avg CPU</CardTitle>
          <Cpu className="h-4 w-4 text-orange-500" />
        </CardHeader>
        <CardContent>
          <div className="text-3xl font-bold">45.3%</div>
          <p className="text-xs text-muted-foreground mt-1">System usage</p>
        </CardContent>
      </Card>
    </div>
  )
}
