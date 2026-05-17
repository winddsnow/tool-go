// /// <reference types="vite/client" /> 是 TypeScript 的三斜线指令
// 它告诉 TypeScript 编译器包含 Vite 客户端类型声明，这样 import.meta.env 等 Vite 特有 API 就有类型提示
/// <reference types="vite/client" />

// 声明 .vue 文件的模块类型
// 在 TypeScript 中，导入 .vue 文件时会报错（因为 TS 默认只识别 .ts 文件）
// 这段声明告诉 TS：所有 .vue 文件导出一个 Vue 组件，这样 import App from './App.vue' 就能通过类型检查
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// Vite 环境变量的类型声明（定义在项目根目录的 .env 文件中）
// VITE_API_BASE_URL 是从 .env 文件加载的环境变量，用于配置后端 API 的地址前缀
// 例如 .env 文件中有 VITE_API_BASE_URL=/api/v1，则代码中通过 import.meta.env.VITE_API_BASE_URL 获取
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
}

// 扩展 ImportMeta 接口，加入 env 属性，这样才能使用 import.meta.env.VITE_XXX
interface ImportMeta {
  readonly env: ImportMetaEnv
}
