/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['localhost', 'kith.com'],
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
    unoptimized: process.env.NODE_ENV !== 'production',
  },
};

export default nextConfig;