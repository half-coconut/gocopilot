import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import eslint from "vite-plugin-eslint";
// import path from "path"; // 引入 path 模块

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), eslint()],
  optimizeDeps: {
    include: ["date-fns"],
  },
});

// export default defineConfig(({ mode }) => {
//   // Load env file based on `mode`
//   const env = loadEnv(mode, process.cwd(), "");

//   return {
//     plugins: [react(), eslint()], // 使用 react 和 eslint 插件
//     optimizeDeps: {
//       include: ["date-fns"],
//     },
//     define: {
//       "process.env": env, // 直接将 env 对象注入到 process.env
//     },
//     resolve: {
//       alias: {
//         "@": path.resolve(__dirname, "src"), //配置src的别名
//       },
//     },
//   };
// });
