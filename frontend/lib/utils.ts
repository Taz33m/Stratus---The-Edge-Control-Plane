import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatUptime(seconds: number): string {
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h`
  return `${Math.floor(seconds / 86400)}d`
}

export function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes.toFixed(0)} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

export function getStatusColor(status: string): string {
  switch (status) {
    case 'running':
      return 'text-green-500'
    case 'stopped':
      return 'text-gray-500'
    case 'error':
      return 'text-red-500'
    case 'starting':
      return 'text-yellow-500'
    default:
      return 'text-gray-500'
  }
}

export function getRegionColor(region: string): string {
  const colors: Record<string, string> = {
    'US-East': 'bg-blue-500',
    'US-West': 'bg-purple-500',
    'EU-West': 'bg-green-500',
    'APAC': 'bg-orange-500',
  }
  return colors[region] || 'bg-gray-500'
}
