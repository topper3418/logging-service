import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import * as path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss()
  ],
  base: "./",
  server: {
    proxy: {
      '/logs': 'http://localhost:8080',
      '/loggers': 'http://localhost:8080'
    },
    fs: {
      allow: [
        path.resolve(__dirname, 'public'), // Allow public assets
        path.resolve(__dirname, 'src'), // Allow source files
      ]
    }
  },
  build: {
    rollupOptions: {
      input: {
        main: 'index.html',
      }
    }
  }
})
