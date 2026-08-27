import path from "path";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import babel from "@rolldown/plugin-babel";
import { defineConfig, loadEnv } from "vite";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");

  return {
    plugins: [react(), tailwindcss(), babel({ presets: [reactCompilerPreset()] })],
    define: {
      "process.env.NODE_ENV": JSON.stringify(env.NODE_ENV ?? "development"),
      // Add more explicit keys here only as needed, e.g.:
      // "process.env.API_URL": JSON.stringify(env.API_URL),
    },
    server: {
      port: 3000,
    },
    resolve: {
      alias: {
        "@": path.resolve(import.meta.dirname, "./src"),
        "@blueprint/openapi": path.resolve(import.meta.dirname, "../../packages/openapi/src"),
        "@blueprint/zod": path.resolve(import.meta.dirname, "../../packages/zod/src"),
      },
    },
  };
});
