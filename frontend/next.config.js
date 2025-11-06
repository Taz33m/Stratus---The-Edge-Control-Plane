/** @type {import('next').NextConfig} */
const nextConfig = {
  // Removed 'output: export' to support WebSocket and server features
  images: {
    unoptimized: true,
  },
}

module.exports = nextConfig
