<template>
  <!--
    router-view 是 Vue Router 提供的占位组件
    当 URL 变化时，Vue Router 会根据路由配置自动匹配对应的页面组件，并渲染在此处
    可以理解为"页面的切换出口"，所有页面都在这里显示
  -->
  <router-view />
</template>

<script setup lang="ts">
// <script setup lang="ts"> 是 Vue 3 的组合式 API 语法糖，写法更简洁
// lang="ts" 表示使用 TypeScript，可以利用类型检查避免运行时错误
</script>

<style>
/* 全局基础样式：清除默认边距和内边距，让页面充满整个浏览器窗口 */
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
}
</style>
