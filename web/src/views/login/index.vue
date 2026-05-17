<template>
  <div class="login-container">
    <div class="login-background">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>
    <div class="login-card">
      <div class="login-header">
        <div class="logo-icon">
          <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="40" height="40" rx="8" fill="url(#gradient)" />
            <path d="M12 20L18 26L28 14" stroke="white" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
            <defs>
              <linearGradient id="gradient" x1="0" y1="0" x2="40" y2="40">
                <stop stop-color="#667eea" />
                <stop offset="1" stop-color="#764ba2" />
              </linearGradient>
            </defs>
          </svg>
        </div>
        <h1>瓦特的工具站</h1>
        <p class="subtitle">Enterprise Management System</p>
      </div>
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin" class="login-form">
        <el-form-item prop="username">
          <el-input 
            v-model="form.username" 
            placeholder="请输入用户名" 
            size="large"
            :prefix-icon="User"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input 
            v-model="form.password" 
            type="password" 
            placeholder="请输入密码" 
            size="large"
            :prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            size="large" 
            class="login-btn" 
            @click="handleLogin"
            :loading="loading"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>
      <div class="login-footer">
        <el-button type="primary" link size="small" @click="router.push('/')">返回工具首页</el-button>
        <span>© 2026 瓦特的工具站. All rights reserved.</span>
        <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer" class="icp-link">粤ICP备2025511523号</a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// ============================================================
// Vue 3 组合式 API —— 本文件是登录页面
// ref ：创建一个响应式值（底层用 Proxy 实现数据劫持）
// reactive：创建一个响应式对象（深层响应，修改任意属性都会触发更新）
// ref 适合基础类型（字符串、数字），reactive 适合对象/数组
// ============================================================
import { ref, reactive } from 'vue'

// useRouter：路由实例，用于页面跳转（如登录成功 -> 首页）
import { useRouter } from 'vue-router'

// Element Plus 图标：User（用户）、Lock（锁）
import { User, Lock } from '@element-plus/icons-vue'

// ElMessage：全局消息提示（成功/错误）
// FormInstance：el-form 的实例类型（用于调用 validate 方法）
// FormRules：表单验证规则的 TypeScript 类型
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

// Pinia store：用户状态管理（存储 token、用户信息、角色等）
import { useUserStore } from '@/store/modules/user'

// authApi：登录/登出的 API 请求函数封装
import { authApi } from '@/api/auth'

const router = useRouter()
const userStore = useUserStore()

// formRef：el-form 的模板引用（用于调用 validate() 表单验证方法）
// ref<FormInstance>() 尖括号是 TypeScript 泛型，声明类型
const formRef = ref<FormInstance>()

// loading：登录中状态（v-loading 指令绑定此值，提交时显示加载动画）
const loading = ref(false)

// ----------------------------------------------------------
// form：登录表单的数据模型（reactive 对象）
// v-model="form.username" 将输入框与数据双向绑定
// 用户输入时会自动更新 form.username
// ----------------------------------------------------------
const form = reactive({
  username: '',
  password: '',
})

// ----------------------------------------------------------
// rules：表单验证规则
// required: true  -> 必填项
// message         -> 验证失败的提示文字
// trigger: 'blur' -> 当输入框失去焦点时触发验证
// el-form 会自动读取 rules 并执行验证
// ----------------------------------------------------------
const rules = reactive<FormRules>({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
})

// ----------------------------------------------------------
// handleLogin：登录按钮点击事件
// async/await：异步函数的语法糖，让代码看起来像同步执行
// @submit.prevent 阻止表单默认提交行为（防止页面刷新）
// ----------------------------------------------------------
const handleLogin = async () => {
  // 如果模板引用不存在，直接返回（类型安全）
  if (!formRef.value) return

  // 调用 el-form 的 validate 方法，执行所有字段的验证规则
  // valid: boolean —— 全部通过为 true，否则 false
  await formRef.value.validate(async (valid) => {
    if (!valid) return  // 验证不通过，停止提交

    loading.value = true  // 显示加载状态，按钮禁用 + 显示"登录中..."
    try {
      // 调用后端登录接口（POST 请求）
      const res = await authApi.login({
        username: form.username,
        password: form.password,
      })

      // 将 JWT token 存入 Pinia store（同时会存到 localStorage 持久化）
      userStore.setToken(res.token)
      // 保存用户信息（ID、用户名、昵称、角色列表）
      userStore.setUserInfo({
        userId: res.user_id,
        username: res.username,
        nickname: res.nickname,
        roles: res.roles,
      })

      // 弹出成功提示
      ElMessage.success('登录成功')
      // 跳转到首页
      router.push('/')
    } catch (err: any) {
      // 登录失败，显示后端返回的错误消息
      ElMessage.error(err.message || '登录失败')
    } finally {
      // 无论成功或失败，finally 块都会执行 —— 关闭加载状态
      loading.value = false
    }
  })
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  position: relative;
  overflow: hidden;
}

.login-background {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
  
  .circle {
    position: absolute;
    border-radius: 50%;
    opacity: 0.1;
  }
  
  .circle-1 {
    width: 400px;
    height: 400px;
    background: #667eea;
    top: -100px;
    right: -100px;
    animation: float 8s ease-in-out infinite;
  }
  
  .circle-2 {
    width: 300px;
    height: 300px;
    background: #764ba2;
    bottom: -50px;
    left: -50px;
    animation: float 6s ease-in-out infinite reverse;
  }
  
  .circle-3 {
    width: 200px;
    height: 200px;
    background: #f093fb;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    animation: float 10s ease-in-out infinite;
  }
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-30px) rotate(10deg); }
}

.login-card {
  width: 420px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  position: relative;
  z-index: 1;
  animation: slideUp 0.5s ease-out;

  @media (max-width: 480px) {
    width: calc(100% - 32px);
    padding: 24px 20px;
    border-radius: 12px;

    h1 { font-size: 22px; }
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
  
  .logo-icon {
    width: 64px;
    height: 64px;
    margin: 0 auto 16px;
    
    svg {
      width: 100%;
      height: 100%;
    }
  }
  
  h1 {
    font-size: 26px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 8px;
  }
  
  .subtitle {
    font-size: 13px;
    color: #8c8c8c;
    margin: 0;
    letter-spacing: 1px;
  }
}

.login-form {
  .login-btn {
    width: 100%;
    height: 44px;
    font-size: 16px;
    font-weight: 500;
    border-radius: 8px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    
    &:hover {
      opacity: 0.9;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
    }
  }
}

.login-footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  margin-top: 24px;
  font-size: 12px;
  color: #8c8c8c;

  .icp-link {
    color: #8c8c8c;
    text-decoration: none;

    &:hover {
      color: #667eea;
    }
  }
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 4px 12px;
  
  &.is-focus {
    box-shadow: 0 0 0 1px #667eea inset;
  }
}
</style>
