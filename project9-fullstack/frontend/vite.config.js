import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    // Forward /api requests to the Gin backend during development,
    // so the browser sees a single origin and CORS is not involved.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
