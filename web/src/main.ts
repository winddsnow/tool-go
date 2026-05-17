// main.ts — 前端应用的入口文件
// Vue 3 应用从这里启动，依次注册 Pinia（状态管理）、Router（路由）、Element Plus（UI 库）
// 最后挂载到 index.html 中的 <div id="app"> 元素上

// Vue 3 核心：createApp 用于创建 Vue 应用实例
import { createApp } from 'vue'
// Pinia：Vue 3 官方推荐的状态管理库，类似 Vuex 但更轻量，用于跨组件共享数据（如用户登录状态）
import { createPinia } from 'pinia'
// Element Plus：基于 Vue 3 的桌面端 UI 组件库，提供表格、表单、弹窗等开箱即用的组件
import ElementPlus from 'element-plus'
// zhCn：Element Plus 的中文语言包，让组件内部文字（如日期选择器的月份）显示为中文
import zhCn from 'element-plus/es/locale/lang/zh-cn'
// 引入 Element Plus 的样式文件
import 'element-plus/dist/index.css'
// 引入 Element Plus 的全部图标组件，并注册为全局组件，这样在任何 .vue 文件中都可以直接使用 <el-icon> 或 <Edit /> 等图标标签
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 根组件 App.vue，所有页面都从这里开始渲染
import App from './App.vue'
// 路由配置文件（src/router/index.ts），定义了 URL 和页面组件的对应关系
import router from './router'
// 全局样式文件，重置浏览器默认样式，设置字体等
import './assets/styles/global.scss'

// 创建 Vue 应用实例
const app = createApp(App)
// 创建 Pinia 状态管理实例
const pinia = createPinia()

// 遍历 Element Plus 图标对象，将每个图标注册为全局组件
// 例如 <Edit /> 会被编译为 element-plus 提供的编辑图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 按顺序安装插件：Pinia -> Router -> Element Plus
app.use(pinia)
app.use(router)
// 安装 Element Plus 并传入中文语言包配置
app.use(ElementPlus, { locale: zhCn })

// 将应用挂载到 index.html 中 id 为 "app" 的 DOM 节点上
app.mount('#app')
