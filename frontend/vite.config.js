import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  base: '/',
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    },
    extensions: ['.vue', '.ts', '.js']
  },
  // 配置代理，将API请求转发到后端
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8000', // 后端服务地址
        changeOrigin: true,
        //rewrite: (path) => path.replace(/^\/api/, '')
      },
      // 添加uploads路径的代理配置
    '/uploads': {
      target: 'http://localhost:8000', // 后端服务地址
      changeOrigin: true
    },
    // 添加站点地图代理配置
    '/sitemap.xml': {
      target: 'http://localhost:8000', // 后端服务地址
      changeOrigin: true
    }
    }
  }
})