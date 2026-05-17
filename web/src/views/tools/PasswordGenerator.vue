<template>
  <div class="password-generator">
    <div class="options">
      <el-form :inline="true" size="small">
        <el-form-item label="长度">
          <el-input-number v-model="length" :min="4" :max="64" :step="1" />
        </el-form-item>
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

    <div class="password-list" v-if="passwords.length > 0">
      <div class="password-item" v-for="(pw, idx) in passwords" :key="idx">
        <span class="index">#{{ idx + 1 }}</span>
        <span class="password">{{ pw }}</span>
        <el-button size="small" type="primary" link @click="copyPassword(pw)">
          复制
        </el-button>
      </div>
    </div>

    <el-empty v-else description="点击「生成」按钮生成密码" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const length = ref(16)
const useUpper = ref(true)
const useLower = ref(true)
const useDigits = ref(true)
const useSpecial = ref(true)
const passwords = ref<string[]>([])
const COUNT = 10

function generateOne(): string {
  const upper = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  const lower = 'abcdefghijklmnopqrstuvwxyz'
  const digits = '0123456789'
  const special = '!@#$%^&*()_+-=[]{}|;:,.<>?'

  let pool = ''
  if (useUpper.value) pool += upper
  if (useLower.value) pool += lower
  if (useDigits.value) pool += digits
  if (useSpecial.value) pool += special

  if (!pool) {
    ElMessage.warning('请至少选择一种字符类型')
    return ''
  }

  let result = ''
  for (let i = 0; i < length.value; i++) {
    result += pool[Math.floor(Math.random() * pool.length)]
  }
  return result
}

function generatePasswords() {
  const list: string[] = []
  for (let i = 0; i < COUNT; i++) {
    const pw = generateOne()
    if (!pw) return
    list.push(pw)
  }
  passwords.value = list
}

function copyPassword(pw: string) {
  navigator.clipboard.writeText(pw).then(() => {
    ElMessage.success('已复制')
  })
}
</script>

<style scoped>
.password-generator {
  padding: 20px;
}

.options {
  margin-bottom: 16px;
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
