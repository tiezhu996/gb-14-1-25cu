import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Vite 配置：开发代理转发 /api 与 /healthz 到后端
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 8010,
    proxy: {
      '/api': { target: 'http://localhost:3010', changeOrigin: true },
      '/healthz': { target: 'http://localhost:3010', changeOrigin: true },
    },
  },
  build: {
    chunkSizeWarningLimit: 4000,
  },
})
