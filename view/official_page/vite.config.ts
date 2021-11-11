import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
// @ts-ignore
const path = require('path')
// import path from 'path'

function _resolve(dir: string) {
  // @ts-ignore
  return path.resolve(__dirname, dir);
}

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': _resolve('src'),
      // '@assets': _resolve('src/assets'),
      '@comps': _resolve('src/components'),
      '@view': _resolve('src/view'),
      '@router': _resolve('src/router'),
      '@store': _resolve('src/store'),
    }
  },
  plugins: [vue()]
})
