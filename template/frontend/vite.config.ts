import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { resolve } from "@std/path/resolve"
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [svelte(), tailwindcss()],
    build: {
        rollupOptions: {
            input: {
                auth: resolve(import.meta.dirname!,  "entrypoints/auth.html"),
                reset: resolve(import.meta.dirname!,  "entrypoints/reset.html"),
                report: resolve(import.meta.dirname!,  "entrypoints/report.html"),
                settings: resolve(import.meta.dirname!,  "entrypoints/settings.html")
            },
            output: {
                dir: "../static-app"
            }
        }
    }
})
