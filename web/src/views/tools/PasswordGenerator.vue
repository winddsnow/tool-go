<template>
  <div class="password-generator">
    <!-- 配置区域：使用 el-form+inline 实现水平布局 -->
    <div class="options">
      <el-form :inline="true" size="small">
        <!-- el-input-number：带步进按钮的数字输入框，控制密码长度（4～64） -->
        <el-form-item label="长度">
          <el-input-number v-model="length" :min="4" :max="64" :step="1" />
        </el-form-item>
        <!-- el-switch：开关组件，控制是否包含某种字符类型 -->
        <el-form-item label="大写字母">
          <el-switch v-model="useUpper" />
        </el-form-item>
        <el-form-item label="小写字母">
          <el-switch v-model="useLower" />
        </el-form-item>
        <el-form-item label="数字">
          <el-switch v-model="useDigits" />
        </el-form-item>
        <el-form-item label="特殊字符">
          <el-switch v-model="useSpecial" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="generatePasswords">生成</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- v-if="passwords.length > 0"：有密码时显示列表 -->
    <div class="password-list" v-if="passwords.length > 0">
      <div class="password-item" v-for="(pw, idx) in passwords" :key="idx">
        <span class="index">#{{ idx + 1 }}</span>
        <span class="password">{{ pw }}</span>
        <el-button size="small" type="primary" link @click="copyPassword(pw)">
          复制
        </el-button>
      </div>
    </div>

    <!-- el-empty：Element Plus 的空状态占位组件，当列表为空时显示提示文字 -->
    <el-empty v-else description="点击「生成」按钮生成密码" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

// 密码选项：长度（默认 16）、四种字符类型开关（默认全部启用）
const length = ref(16)
const useUpper = ref(true)
const useLower = ref(true)
const useDigits = ref(true)
const useSpecial = ref(true)
// 生成的密码列表
const passwords = ref<string[]>([])
// 每次生成 10 条密码
const COUNT = 10

// 生成单条随机密码
function generateOne(): string {
  // 定义四类字符池
  const upper = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  const lower = 'abcdefghijklmnopqrstuvwxyz'
  const digits = '0123456789'
  const special = '!@#$%^&*()_+-=[]{}|;:,.<>?'

  // 根据用户开关，动态拼接可用的字符池
  let pool = ''
  if (useUpper.value) pool += upper
  if (useLower.value) pool += lower
  if (useDigits.value) pool += digits
  if (useSpecial.value) pool += special

  // 如果所有开关都关掉了，字符池为空，提示用户
  if (!pool) {
    ElMessage.warning('请至少选择一种字符类型')
    return ''
  }

  // 从拼接好的字符池中随机取 length 次字符，组成密码
  let result = ''
  for (let i = 0; i < length.value; i++) {
    result += pool[Math.floor(Math.random() * pool.length)]
  }
  return result
}

// 生成 COUNT 条密码并更新列表
function generatePasswords() {
  const list: string[] = []
  for (let i = 0; i < COUNT; i++) {
    const pw = generateOne()
    if (!pw) return // 如果字符池为空则停止
    list.push(pw)
  }
  passwords.value = list
}

// 复制单条密码到剪贴板
function copyPassword(pw: string) {
  navigator.clipboard.writeText(pw).then(() => {
    ElMessage.success('已复制')
  })
}
</script>

<style scoped>
.password-generator {
  padding: 20px;

  @media (max-width: 480px) {
    padding: 12px 0;
  }
}

.options {
  margin-bottom: 16px;

  @media (max-width: 480px) {
    :deep(.el-form-item) {
      display: block;
      margin-right: 0;
      margin-bottom: 8px;
    }
    :deep(.el-form-item__label) {
      display: inline-block;
      width: 80px;
    }
    :deep(.el-form-item__content) {
      display: inline-block;
    }
  }
}

.password-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.password-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.index {
  color: #909399;
  font-size: 13px;
  min-width: 28px;
}

.password {
  flex: 1;
  font-family: monospace;
  font-size: 14px;
  color: #409eff;
  word-break: break-all;
  user-select: all;
}
</style>
