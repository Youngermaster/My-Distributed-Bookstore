import path from "path";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 5173,
    proxy: {
      // Route cart requests to cart service
      "/api/v1/cart": {
        target: "http://localhost:8083",
        changeOrigin: true,
      },
      // Route order requests to order service
      "/api/v1/orders": {
        target: "http://localhost:8084",
        changeOrigin: true,
      },
      // Route catalog requests to catalog service
      "/api/v1/catalog": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      // Route user/auth requests to user service
      "/api/v1/users": {
        target: "http://localhost:8081",
        changeOrigin: true,
      },
      "/api/v1/auth": {
        target: "http://localhost:8081",
        changeOrigin: true,
      },
    },
  },
});
