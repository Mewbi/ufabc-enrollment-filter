import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig(({ command, mode, isSsrBuild, isPreview }) => {
  if (command === "serve") {
    return {
      plugins: [vue()],
      base: "/ufabc-enrollment-filter/",
    };
  } else {
    // command === 'build'
    return {
      plugins: [vue()],
      base: "/ufabc-enrollment-filter/",
    };
  }
});
