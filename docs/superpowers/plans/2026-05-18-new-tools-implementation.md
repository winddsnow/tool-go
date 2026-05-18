# 新增 5 个工具箱工具 — 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在工具箱首页增加分类 Tab（el-tabs）并新增 5 个纯前端工具组件

**Architecture:** 所有工具通过 dialog 弹窗渲染，与现有 6 个工具模式一致。新增 `category` 字段按 Tab 分类过滤卡片。TextDiff 自实现 LCS 行级 diff；CaseConverter、RegexTester、UuidGenerator 纯前端无依赖；QrCodeGenerator 依赖 `qrcode` npm 包。

**Tech Stack:** Vue 3 + TypeScript + Element Plus + `qrcode`

---

### Task 1: 安装 qrcode 依赖

**Files:**
- Modify: `web/package.json`

- [ ] **安装 qrcode 包**

```bash
cd /home/walter/myopencode/tool-go/web && npm install qrcode
```

Expected: `qrcode` 包被添加到 `dependencies`，package.json 自动更新。

---

### Task 2: 修改 tools/index.vue — 增加 Tab 分类 + 注册 5 个新组件

**Files:**
- Modify: `web/src/views/tools/index.vue`

Interface `Tool` 增加 `category` 字段，tools 注册表增加新工具和分类，模板增加 `el-tabs` 过滤逻辑。

- [ ] **更新 Tool 接口和 tools 注册表**

在 `web/src/views/tools/index.vue` 中，将 `Tool` 接口改为：

```typescript
interface Tool {
  id: string
  title: string
  category: string
  component: any
}
```

新增 5 个组件导入：

```typescript
import TextDiff from './TextDiff.vue'
import RegexTester from './RegexTester.vue'
import CaseConverter from './CaseConverter.vue'
import UuidGenerator from './UuidGenerator.vue'
import QrCodeGenerator from './QrCodeGenerator.vue'
```

更新 tools 注册表，所有 11 个工具带上 `category`：

```typescript
const tools: Record<string, Tool> = {
  timestamp: { id: 'timestamp', title: '时间戳转换工具', category: 'convert', component: TimestampConverter },
  json: { id: 'json', title: 'JSON 格式化工具', category: 'text', component: JsonFormatter },
  hash: { id: 'hash', title: '哈希加密工具', category: 'encode', component: HashEncryptor },
  base64: { id: 'base64', title: 'Base64 编解码工具', category: 'encode', component: Base64Converter },
  password: { id: 'password', title: '密码生成器', category: 'generate', component: PasswordGenerator },
  mockdata: { id: 'mockdata', title: '随机数据生成器', category: 'generate', component: MockDataGenerator },
  textdiff: { id: 'textdiff', title: '文本对比工具', category: 'text', component: TextDiff },
  regex: { id: 'regex', title: '正则表达式测试工具', category: 'text', component: RegexTester },
  caseconv: { id: 'caseconv', title: '大小写/Naming Case 转换', category: 'text', component: CaseConverter },
  uuid: { id: 'uuid', title: 'UUID 生成器', category: 'generate', component: UuidGenerator },
  qrcode: { id: 'qrcode', title: '二维码生成器', category: 'generate', component: QrCodeGenerator },
}
```

- [ ] **添加 category 和 activeTab 状态**

在 `currentTool` 下方添加：

```typescript
const categories = [
  { key: 'all', label: '全部' },
  { key: 'text', label: '文本处理' },
  { key: 'encode', label: '编码加密' },
  { key: 'generate', label: '生成类' },
  { key: 'convert', label: '转换类' },
]
const activeTab = ref('all')
```

- [ ] **修改模板 — 在 page-header 和 tools-grid 之间插入 el-tabs**

替换原有的 `tools-grid` 部分：用分类过滤的计算属性驱动显示。在开源信息栏（`oss-banner`）下方添加：

```html
<div class="category-tabs">
  <el-tabs v-model="activeTab" @tab-click="activeTab = $event.props.name">
    <el-tab-pane v-for="cat in categories" :key="cat.key" :label="cat.label" :name="cat.key" />
  </el-tabs>
</div>
```

- [ ] **添加 filteredTools 计算属性**

在 `topUserLabel` 附近添加：

```typescript
const filteredTools = computed(() => {
  return Object.values(tools).filter(t => activeTab.value === 'all' || t.category === activeTab.value)
})
```

将模板中 `tools-grid` 的 `v-for` 改为遍历 `filteredTools`，卡片改为动态渲染：

```html
<div class="tools-grid">
  <div
    v-for="tool in filteredTools"
    :key="tool.id"
    class="tool-card"
    @click="openTool(tool.id)"
  >
    <!-- 保持现有卡片样式，内容改为动态 -->
  </div>
</div>
```

由于每个卡片图标不同，改为动态图标渲染。为了简化，给 Tool 接口新增 `icon` 字段：

```typescript
interface Tool {
  id: string
  title: string
  category: string
  icon: string
  component: any
}
```

注册每个工具的图标名：

```typescript
timestamp: { id: 'timestamp', title: '时间戳转换工具', category: 'convert', icon: 'Clock', component: TimestampConverter },
json: { id: 'json', title: 'JSON 格式化工具', category: 'text', icon: 'Edit', component: JsonFormatter },
hash: { id: 'hash', title: '哈希加密工具', category: 'encode', icon: 'Lock', component: HashEncryptor },
base64: { id: 'base64', title: 'Base64 编解码工具', category: 'encode', icon: 'Document', component: Base64Converter },
password: { id: 'password', title: '密码生成器', category: 'generate', icon: 'Key', component: PasswordGenerator },
mockdata: { id: 'mockdata', title: '随机数据生成器', category: 'generate', icon: 'List', component: MockDataGenerator },
textdiff: { id: 'textdiff', title: '文本对比工具', category: 'text', icon: 'EditPen', component: TextDiff },
regex: { id: 'regex', title: '正则表达式测试工具', category: 'text', icon: 'Search', component: RegexTester },
caseconv: { id: 'caseconv', title: '大小写/Naming Case 转换', category: 'text', icon: 'EditPen', component: CaseConverter },
uuid: { id: 'uuid', title: 'UUID 生成器', category: 'generate', icon: 'Link', component: UuidGenerator },
qrcode: { id: 'qrcode', title: '二维码生成器', category: 'generate', icon: 'Ticket', component: QrCodeGenerator },
```

卡片模板改为动态图标：

```html
<div v-for="tool in filteredTools" :key="tool.id" class="tool-card" @click="openTool(tool.id)">
  <div class="tool-icon">
    <el-icon :size="40"><component :is="icons[tool.icon]" /></el-icon>
  </div>
  <div class="tool-info">
    <h3>{{ tool.title }}</h3>
    <p>{{ toolDescriptions[tool.id] }}</p>
  </div>
  <div class="tool-arrow">
    <el-icon><ArrowRight /></el-icon>
  </div>
</div>
```

导入新图标（从 `@element-plus/icons-vue` 导入 `EditPen`, `Search`, `Ticket`, `Link` 等，加上原有的 `Clock`, `Edit`, `Lock`, `Document`, `Key`, `List`）：

```typescript
import { Clock, Edit, Lock, Document, Key, List, EditPen, Search, Link, Ticket, ArrowRight } from '@element-plus/icons-vue'
```

添加 `icons` 映射：

```typescript
const icons: Record<string, any> = { Clock, Edit, Lock, Document, Key, List, EditPen, Search, Link, Ticket }
```

添加 `toolDescriptions` 映射：

```typescript
const toolDescriptions: Record<string, string> = {
  timestamp: 'Unix 时间戳与日期时间双向转换',
  json: 'JSON 数据美化与校验',
  hash: 'MD5, SHA1, SHA256 加密工具',
  base64: 'Base64 编码与解码工具',
  password: '随机密码生成，支持自定义规则',
  mockdata: '生成姓名/手机/身份证/护照等模拟数据',
  textdiff: '逐行对比两段文本差异',
  regex: '正则表达式匹配与测试',
  caseconv: 'camelCase/snake_case/kebab-case 互转',
  uuid: '批量生成 UUID v1/v4',
  qrcode: '将文本或 URL 生成二维码',
}
```

移除原有的硬编码卡片模板（时间戳、JSON、哈希、Base64、密码、Mock 数据那 6 个静态 `.tool-card` div，因为现在由 `v-for` 动态渲染）。

- [ ] **为 `el-tabs` 添加样式**

在 `<style scoped lang="scss">` 中添加：

```scss
.category-tabs {
  margin-bottom: 24px;
  background: #fff;
  border-radius: 12px;
  padding: 8px 16px 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #f0f0f0;

  .el-tabs {
    .el-tabs__header {
      margin: 0;
    }
  }

  @media (max-width: 480px) {
    padding: 8px 8px 0;
    overflow-x: auto;
  }
}
```

---

### Task 3: 创建 TextDiff.vue

**Files:**
- Create: `web/src/views/tools/TextDiff.vue`

```vue
<template>
  <div class="text-diff">
    <div class="diff-inputs" :class="{ 'is-mobile': isMobile }">
      <div class="input-area">
        <span class="label">原文本</span>
        <el-input v-model="oldText" type="textarea" :rows="8" placeholder="请输入原文本" size="large" />
      </div>
      <div class="input-area">
        <span class="label">新文本</span>
        <el-input v-model="newText" type="textarea" :rows="8" placeholder="请输入新文本" size="large" />
      </div>
    </div>
    <div class="action-bar">
      <el-button type="primary" @click="computeDiff">对比</el-button>
      <el-button @click="clear">清空</el-button>
    </div>
    <div class="diff-output" v-if="diffResult.length > 0">
      <div
        v-for="(line, i) in diffResult"
        :key="i"
        class="diff-line"
        :class="line.type"
      >
        <span class="line-prefix">{{ line.prefix }}</span>
        <span class="line-text">{{ line.text }}</span>
      </div>
    </div>
    <div class="no-diff" v-else-if="compared">
      <el-empty :description="oldText === newText ? '两段文本完全相同' : '请输入文本进行对比'" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const oldText = ref('')
const newText = ref('')
const diffResult = ref<{ type: string; prefix: string; text: string }[]>([])
const compared = ref(false)
const isMobile = ref(window.innerWidth < 768)

window.addEventListener('resize', () => { isMobile.value = window.innerWidth < 768 })

function computeDiff() {
  compared.value = true
  if (!oldText.value && !newText.value) {
    diffResult.value = []
    return
  }
  const oldLines = oldText.value.split('\n')
  const newLines = newText.value.split('\n')
  const lcs = computeLCS(oldLines, newLines)
  diffResult.value = buildDiff(oldLines, newLines, lcs)
}

function computeLCS(a: string[], b: string[]): number[][] {
  const m = a.length, n = b.length
  const dp: number[][] = Array.from({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      if (a[i - 1] === b[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1
      } else {
        dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1])
      }
    }
  }
  return dp
}

function buildDiff(a: string[], b: string[], dp: number[][]): { type: string; prefix: string; text: string }[] {
  const result: { type: string; prefix: string; text: string }[] = []
  let i = a.length, j = b.length
  const temp: { type: string; prefix: string; text: string }[] = []
  while (i > 0 || j > 0) {
    if (i > 0 && j > 0 && a[i - 1] === b[j - 1]) {
      temp.push({ type: 'equal', prefix: '  ', text: a[i - 1] })
      i--; j--
    } else if (j > 0 && (i === 0 || dp[i][j - 1] >= dp[i - 1][j])) {
      temp.push({ type: 'add', prefix: '+ ', text: b[j - 1] })
      j--
    } else {
      temp.push({ type: 'remove', prefix: '- ', text: a[i - 1] })
      i--
    }
  }
  return temp.reverse()
}

function clear() {
  oldText.value = ''
  newText.value = ''
  diffResult.value = []
  compared.value = false
}
</script>

<style scoped lang="scss">
.text-diff {
  .diff-inputs {
    display: flex;
    gap: 16px;
    margin-bottom: 16px;

    &.is-mobile {
      flex-direction: column;
    }

    .input-area {
      flex: 1;
      .label {
        display: block;
        font-size: 14px;
        font-weight: 600;
        color: #303133;
        margin-bottom: 8px;
      }
    }
  }

  .action-bar {
    margin-bottom: 16px;
  }

  .diff-output {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 12px;
    max-height: 400px;
    overflow-y: auto;
    font-family: 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.6;

    .diff-line {
      display: flex;
      padding: 2px 8px;
      border-radius: 4px;
      margin-bottom: 1px;

      &.add {
        background: #e6ffed;
        .line-prefix { color: #28a745; }
      }
      &.remove {
        background: #ffeef0;
        .line-prefix { color: #d73a49; }
      }
      &.equal {
        .line-prefix { color: #6a737d; }
      }

      .line-prefix {
        width: 24px;
        flex-shrink: 0;
      }
      .line-text {
        white-space: pre-wrap;
        word-break: break-all;
      }
    }
  }

  .no-diff {
    margin-top: 16px;
  }
}
</style>
```

---

### Task 4: 创建 RegexTester.vue

**Files:**
- Create: `web/src/views/tools/RegexTester.vue`

```vue
<template>
  <div class="regex-tester">
    <div class="input-area">
      <span class="label">正则表达式</span>
      <el-input v-model="pattern" placeholder="/ 输入正则表达式 /" size="large" @input="testRegex" />
    </div>
    <div class="test-text-area">
      <span class="label">测试文本</span>
      <el-input v-model="testText" type="textarea" :rows="6" placeholder="请输入测试文本" size="large" @input="testRegex" />
    </div>
    <div class="options">
      <el-switch v-model="globalFlag" active-text="全局匹配 (g)" @change="testRegex" />
      <el-switch v-model="ignoreCaseFlag" active-text="忽略大小写 (i)" @change="testRegex" />
    </div>
    <div class="error-msg" v-if="error">{{ error }}</div>
    <div class="result" v-if="!error && matches.length > 0">
      <div class="match-count">匹配 {{ matches.length }} 处</div>
      <div class="match-list">
        <div v-for="(m, i) in matches" :key="i" class="match-item">
          <span class="match-index">#{{ i + 1 }}</span>
          <span class="match-text">{{ m[0] }}</span>
          <span class="match-pos">位置: {{ m.index }}</span>
        </div>
      </div>
    </div>
    <div class="no-match" v-else-if="!error && tested">无匹配结果</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const pattern = ref('')
const testText = ref('')
const globalFlag = ref(true)
const ignoreCaseFlag = ref(false)
const matches = ref<RegExpExecArray[]>([])
const error = ref('')
const tested = ref(false)

function testRegex() {
  error.value = ''
  matches.value = []
  tested.value = false
  if (!pattern.value || !testText.value) return
  tested.value = true
  try {
    const flags = `${globalFlag.value ? 'g' : ''}${ignoreCaseFlag.value ? 'i' : ''}`
    const regex = new RegExp(pattern.value, flags)
    const results: RegExpExecArray[] = []
    let match: RegExpExecArray | null
    if (globalFlag.value) {
      while ((match = regex.exec(testText.value)) !== null) {
        results.push(match)
      }
    } else {
      match = regex.exec(testText.value)
      if (match) results.push(match)
    }
    matches.value = results
  } catch (e: any) {
    error.value = `正则错误: ${e.message}`
  }
}
</script>

<style scoped lang="scss">
.regex-tester {
  .input-area, .test-text-area {
    margin-bottom: 16px;
    .label {
      display: block;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 8px;
    }
  }
  .options {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }
  .error-msg {
    color: #f56c6c;
    font-size: 13px;
    margin-bottom: 12px;
    padding: 8px 12px;
    background: #fef0f0;
    border-radius: 6px;
  }
  .result {
    .match-count {
      font-size: 14px;
      font-weight: 600;
      color: #409eff;
      margin-bottom: 12px;
    }
    .match-list {
      max-height: 300px;
      overflow-y: auto;
      .match-item {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 6px 12px;
        background: #f0f9ff;
        border-radius: 6px;
        margin-bottom: 4px;
        font-size: 13px;
        .match-index {
          color: #909399;
          min-width: 28px;
        }
        .match-text {
          font-family: 'Courier New', monospace;
          color: #303133;
          flex: 1;
          word-break: break-all;
        }
        .match-pos {
          color: #909399;
          font-size: 12px;
          flex-shrink: 0;
        }
      }
    }
  }
  .no-match {
    color: #909399;
    font-size: 14px;
    padding: 12px 0;
  }
}
</style>
```

---

### Task 5: 创建 CaseConverter.vue

**Files:**
- Create: `web/src/views/tools/CaseConverter.vue`

```vue
<template>
  <div class="case-converter">
    <div class="input-area">
      <span class="label">源文本</span>
      <el-input v-model="inputText" type="textarea" :rows="4" placeholder="请输入要转换的文本，如 hello world" size="large" />
    </div>
    <div class="actions" :class="{ 'is-mobile': isMobile }">
      <el-button v-for="conv in converters" :key="conv.key" @click="convert(conv.key)">
        <span class="conv-label">{{ conv.label }}</span>
        <span class="conv-example">{{ conv.example }}</span>
      </el-button>
    </div>
    <div class="output-area">
      <span class="label">转换结果</span>
      <el-input v-model="outputText" type="textarea" :rows="4" placeholder="点击上方按钮查看转换结果" size="large" readonly />
    </div>
    <el-button v-if="outputText" size="small" link type="primary" @click="copyResult">复制结果</el-button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const inputText = ref('')
const outputText = ref('')
const isMobile = ref(window.innerWidth < 768)
window.addEventListener('resize', () => { isMobile.value = window.innerWidth < 768 })

const converters = [
  { key: 'camel', label: 'camelCase', example: '→ helloWorld' },
  { key: 'pascal', label: 'PascalCase', example: '→ HelloWorld' },
  { key: 'snake', label: 'snake_case', example: '→ hello_world' },
  { key: 'kebab', label: 'kebab-case', example: '→ hello-world' },
  { key: 'constant', label: 'CONSTANT_CASE', example: '→ HELLO_WORLD' },
]

function splitWords(text: string): string[] {
  const trimmed = text.trim()
  if (!trimmed) return []
  // 按空格、连字符、下划线、大小写边界拆分
  const parts = trimmed.split(/[\s_-]+/)
  const words: string[] = []
  for (const part of parts) {
    if (!part) continue
    // 进一步按大小写边界拆分：HelloWorld → Hello, World
    const sub = part.replace(/([A-Z])/g, ' $1').trim().split(/\s+/)
    words.push(...sub.map(w => w.toLowerCase()))
  }
  return words.filter(Boolean)
}

function convert(type: string) {
  const words = splitWords(inputText.value)
  if (words.length === 0) {
    ElMessage.warning('请输入要转换的文本')
    return
  }
  switch (type) {
    case 'camel':
      outputText.value = words[0] + words.slice(1).map(w => w.charAt(0).toUpperCase() + w.slice(1)).join('')
      break
    case 'pascal':
      outputText.value = words.map(w => w.charAt(0).toUpperCase() + w.slice(1)).join('')
      break
    case 'snake':
      outputText.value = words.join('_')
      break
    case 'kebab':
      outputText.value = words.join('-')
      break
    case 'constant':
      outputText.value = words.map(w => w.toUpperCase()).join('_')
      break
  }
}

function copyResult() {
  navigator.clipboard.writeText(outputText.value)
  ElMessage.success('已复制')
}
</script>

<style scoped lang="scss">
.case-converter {
  .input-area, .output-area {
    margin-bottom: 16px;
    .label {
      display: block;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 8px;
    }
  }
  .actions {
    display: flex;
    gap: 8px;
    margin-bottom: 16px;
    flex-wrap: wrap;

    .el-button {
      display: flex;
      flex-direction: column;
      align-items: center;
      height: auto;
      padding: 8px 16px;

      .conv-label {
        font-size: 13px;
        font-weight: 600;
      }
      .conv-example {
        font-size: 11px;
        color: #909399;
        margin-top: 2px;
      }
    }

    &.is-mobile {
      .el-button {
        flex: 1;
        min-width: 80px;
      }
    }
  }
}
</style>
```

---

### Task 6: 创建 UuidGenerator.vue

**Files:**
- Create: `web/src/views/tools/UuidGenerator.vue`

```vue
<template>
  <div class="uuid-generator">
    <div class="options">
      <div class="option-item">
        <span class="label">UUID 版本</span>
        <el-select v-model="version" size="large" style="width: 160px">
          <el-option label="UUID v4 (随机)" value="v4" />
          <el-option label="UUID v1 (时间戳)" value="v1" />
        </el-select>
      </div>
      <div class="option-item">
        <span class="label">生成数量</span>
        <el-input-number v-model="count" :min="1" :max="20" size="large" />
      </div>
    </div>
    <div class="action-bar">
      <el-button type="primary" @click="generate">生成</el-button>
    </div>
    <div class="result" v-if="uuids.length > 0">
      <div class="uuid-item" v-for="(uuid, i) in uuids" :key="i">
        <span class="uuid-text">{{ uuid }}</span>
        <el-button size="small" link type="primary" @click="copy(uuid)">复制</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const version = ref('v4')
const count = ref(5)
const uuids = ref<string[]>([])

function generate() {
  uuids.value = []
  for (let i = 0; i < count.value; i++) {
    uuids.value.push(version.value === 'v4' ? generateV4() : generateV1())
  }
}

function generateV4(): string {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID()
  }
  const hex = '0123456789abcdef'
  const sections = [8, 4, 4, 4, 12]
  return sections.map(len => {
    let s = ''
    for (let i = 0; i < len; i++) s += hex[Math.floor(Math.random() * 16)]
    return s
  }).join('-')
}

function generateV1(): string {
  const now = Date.now()
  const timeHex = now.toString(16).padStart(12, '0')
  const mac = Math.floor(Math.random() * 0xffffff).toString(16).padStart(6, '0')
  const clockSeq = Math.floor(Math.random() * 0x3fff).toString(16).padStart(4, '0')
  return `${timeHex.slice(0, 8)}-${timeHex.slice(8, 12)}-1${timeHex.slice(12, 15)}-${clockSeq}-${mac}${Math.floor(Math.random() * 0xffffff).toString(16).padStart(6, '0')}`
}

function copy(text: string) {
  navigator.clipboard.writeText(text)
  ElMessage.success('已复制')
}
</script>

<style scoped lang="scss">
.uuid-generator {
  .options {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
    flex-wrap: wrap;

    .option-item {
      display: flex;
      align-items: center;
      gap: 12px;

      .label {
        font-size: 14px;
        font-weight: 600;
        color: #303133;
        flex-shrink: 0;
      }
    }
  }

  .action-bar {
    margin-bottom: 16px;
  }

  .result {
    max-height: 400px;
    overflow-y: auto;

    .uuid-item {
      display: flex;
      align-items: center;
      padding: 8px 12px;
      background: #f8f9fa;
      border-radius: 6px;
      margin-bottom: 4px;

      .uuid-text {
        font-family: 'Courier New', monospace;
        font-size: 13px;
        color: #303133;
        flex: 1;
        word-break: break-all;
      }
    }
  }
}
</style>
```

---

### Task 7: 创建 QrCodeGenerator.vue

**Files:**
- Create: `web/src/views/tools/QrCodeGenerator.vue`

```vue
<template>
  <div class="qr-code-generator">
    <div class="input-area">
      <span class="label">文本或 URL</span>
      <el-input v-model="text" placeholder="请输入文本或 URL" size="large" @keyup.enter="generate" />
    </div>
    <div class="action-bar">
      <el-button type="primary" @click="generate">生成二维码</el-button>
    </div>
    <div class="qr-result" v-if="showQR">
      <canvas ref="canvasRef"></canvas>
      <div class="qr-actions">
        <el-button type="success" size="small" @click="saveImage">保存为图片</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import QRCode from 'qrcode'

const text = ref('')
const showQR = ref(false)
const canvasRef = ref<HTMLCanvasElement | null>(null)

async function generate() {
  if (!text.value) {
    ElMessage.warning('请输入文本或 URL')
    return
  }
  showQR.value = true
  await nextTick()
  if (canvasRef.value) {
    try {
      await QRCode.toCanvas(canvasRef.value, text.value, { width: 256, margin: 2 })
    } catch {
      ElMessage.error('生成二维码失败')
    }
  }
}

function saveImage() {
  if (!canvasRef.value) return
  const link = document.createElement('a')
  link.download = 'qrcode.png'
  link.href = canvasRef.value.toDataURL('image/png')
  link.click()
}
</script>

<style scoped lang="scss">
.qr-code-generator {
  .input-area {
    margin-bottom: 16px;
    .label {
      display: block;
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 8px;
    }
  }

  .action-bar {
    margin-bottom: 16px;
  }

  .qr-result {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;

    canvas {
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
  }
}
</style>
```

---

### Task 8: 验证构建

- [ ] **运行 vue-tsc 类型检查**

```bash
cd /home/walter/myopencode/tool-go/web && npx vue-tsc --noEmit
```

Expected: 无类型错误

- [ ] **运行 vite 构建**

```bash
cd /home/walter/myopencode/tool-go/web && npx vite build
```

Expected: 构建成功，无报错
