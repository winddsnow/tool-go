# 新增 5 个工具箱工具 — 设计文档

## 概述

在现有工具箱基础上新增 5 个纯前端工具，工具箱首页增加分类 Tab 切换。所有工具通过弹窗（el-dialog）渲染，与现有模式一致。

## 分类系统

### Tab 定义

| key | 标签 | 包含工具 |
|-----|------|---------|
| `all` | 全部 | 11 个工具 |
| `text` | 文本处理 | JSON 格式化、文本对比、正则测试、大小写转换 |
| `encode` | 编码加密 | Base64 编解码、哈希加密 |
| `generate` | 生成类 | 密码生成器、随机数据生成器、UUID 生成、二维码生成 |
| `convert` | 转换类 | 时间戳转换 |

### 实现方式

- 所有工具注册时增加 `category` 字段
- 使用 `el-tabs` 切换分类
- 选中 Tab 时过滤卡片显示
- Tab 栏跟随顶部布局，在统计栏和开源信息栏之间

## 各工具设计

### 1. 文本对比 (TextDiff)

**功能**：逐行对比两段文本的差异

**UI**：
- 两个 `el-input type="textarea"`：原文本 / 新文本
- 按钮：`对比`
- 输出：逐行 diff 结果，新增行绿色背景 + `+` 前缀，删除行红色背景 + `-` 前缀，不变行无底色
- 空文本或相同时显示友好提示

**移动端**：桌面端左右并排，移动端上下堆叠

**实现**：自行实现 LCS（最长公共子序列）行级 diff 算法，无额外依赖

---

### 2. 正则表达式测试 (RegexTester)

**功能**：输入正则和测试文本，实时显示匹配结果

**UI**：
- `el-input` 输入正则表达式（附加分隔符提示）
- `el-input type="textarea"` 输入测试文本
- 两个 `el-switch`：全局匹配 (`g`)、忽略大小写 (`i`)
- 结果区：匹配总数、每条匹配项的详情（匹配文本 + 位置索引）
- 正则非法时显示错误提示

**移动端**：上下排列，标签左对齐

---

### 3. 大小写/Naming Case 转换 (CaseConverter)

**功能**：在多种命名格式间互相转换

**UI**：
- `el-input type="textarea"` 输入源文本
- 5 个按钮排成一行，桌面端水平，移动端换行：
  - `camelCase` → `helloWorld`
  - `PascalCase` → `HelloWorld`
  - `snake_case` → `hello_world`
  - `kebab-case` → `hello-world`
  - `CONSTANT_CASE` → `HELLO_WORLD`
- 输出 `el-input type="textarea"` 显示结果，带一键复制
- 点击某按钮后自动转换并显示结果

**实现**：纯前端正则处理。分词规则：按空格、连字符、下划线、大小写边界拆分。

---

### 4. UUID 生成 (UuidGenerator)

**功能**：批量生成 UUID

**UI**：
- `el-select` 选择版本：UUID v4（随机）、UUID v1（时间戳，用纯前端模拟）
- `el-input-number` 数量：1-20，默认 5
- 按钮：`生成`
- 输出：列表展示，每行带 UUID 文本 + 复制按钮
- 列表上限显示提示

**实现**：
- 优先使用 `crypto.randomUUID()`（UUID v4）
- 回退方案：`crypto.getRandomValues` + 格式拼接
- v1 用当前时间戳 + MAC 地址模拟格式

---

### 5. 二维码生成 (QrCodeGenerator)

**功能**：输入文本/URL 生成二维码，支持保存为图片

**UI**：
- `el-input` 输入内容
- 按钮：`生成二维码`
- 二维码渲染为 `<canvas>`，居中显示
- 底部：`保存为图片` 按钮（`canvas.toDataURL` → 下载）

**依赖**：安装 `qrcode` npm 包（`npm install qrcode`，约 23KB）

---

## 移动端兼容

所有新工具遵循以下规则：
- 使用现有 `isMobile` 计算属性判断移动端
- 桌面端 dialog 宽度 `800px`，移动端 `95%`
- 并排布局在移动端改为上下堆叠（flex-direction: column）
- 按钮行在移动端允许换行
- textarea 在移动端高度适当缩减
- `el-tabs` 支持横向滚动（移动端 Tab 过多时）

## 依赖变更

新增 `qrcode` 包，生产依赖。

## 变更文件清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `web/src/views/tools/index.vue` | 修改 | 加 Tab 过滤 + 注册 5 个新组件 |
| `web/src/views/tools/TextDiff.vue` | 新增 | 文本对比工具 |
| `web/src/views/tools/RegexTester.vue` | 新增 | 正则测试工具 |
| `web/src/views/tools/CaseConverter.vue` | 新增 | 大小写/命名格式转换 |
| `web/src/views/tools/UuidGenerator.vue` | 新增 | UUID 生成工具 |
| `web/src/views/tools/QrCodeGenerator.vue` | 新增 | 二维码生成工具 |
| `web/package.json` | 修改 | 新增 `qrcode` 依赖 |

## 不变更的

- 后端代码 — 所有新工具纯前端
- 路由配置 — 通过 dialog 弹窗，不新增页面路由
- 现有 6 个工具组件 — 只改注册方式（加 category），不改功能
