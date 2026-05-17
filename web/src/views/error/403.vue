<template>
  <div class="error-container">
    <div class="error-content">
      <div class="error-code">403</div>
      <div class="error-icon">
        <svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
          <circle cx="100" cy="100" r="80" fill="none" stroke="#e6e6e6" stroke-width="4" />
          <path d="M70 70 L130 130 M130 70 L70 130" stroke="#ff6b6b" stroke-width="8" stroke-linecap="round" />
        </svg>
      </div>
      <h1>抱歉，您没有权限访问此页面</h1>
      <p>请联系管理员获取相应权限</p>
      <el-button type="primary" @click="goBack">返回首页</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.push('/')
}
</script>

<style scoped lang="scss">
.error-container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
}

.error-content {
  text-align: center;
  
  .error-code {
    font-size: 120px;
    font-weight: bold;
    color: #e6e6e6;
    line-height: 1;
    margin-bottom: -40px;
  }
  
  .error-icon {
    width: 200px;
    height: 200px;
    margin: 0 auto 30px;
    
    svg {
      width: 100%;
      height: 100%;
    }
  }
  
  h1 {
    font-size: 24px;
    color: #303133;
    margin-bottom: 12px;
  }
  
  p {
    font-size: 14px;
    color: #909399;
    margin-bottom: 30px;
  }
}
</style>
