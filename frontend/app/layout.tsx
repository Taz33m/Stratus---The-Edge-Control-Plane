import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Stratus - Edge Control Plane',
  description: 'Command your microservices from the edge',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="cds--g100">{children}</body>
    </html>
  )
}
