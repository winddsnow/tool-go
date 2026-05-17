// ============================================================
// package v1 — API v1 的工具类接口数据结构
// ------------------------------------------------------------
// 工具类接口提供一些辅助功能，不涉及核心业务。
// 如生成测试用的模拟数据（mock data）。
// ============================================================
package v1

import "github.com/gogf/gf/v2/frame/g"

// ============================================================
// MockDataReq — 生成模拟数据请求
// ------------------------------------------------------------
//   Types  []string — 要生成的数据类型列表（可多个，如 ["name","phone"]）
//                     v:"required#..." 确保至少选一个类型
//
//   Count  int      — 要生成的数据条数
//                     v:"between:1,100#数量范围1-100" — 范围验证规则。
//                     between:1,100 表示值必须在 1 到 100 之间（含两端）。
//                     `#` 后面是自定义错误信息。
//                     多条验证规则用空格分隔：
//                     d:"10" 是默认值，v 是验证规则。
//                     注意：d 和 v 是独立的两个 tag，不是同一个。
// ============================================================
type MockDataReq struct {
	g.Meta `path:"/tools/mock-data" method:"post" tags:"Tools" summary:"生成模拟数据"`
	Types  []string `json:"types" v:"required#请选择至少一个数据类型" dc:"数据类型: name,phone,email,id_card,passport,address,ip,datetime,bank_card"`
	Count  int      `json:"count" d:"10" v:"between:1,100#数量范围1-100" dc:"生成数量"`
}

// ============================================================
// MockDataRes — 模拟数据响应
// ------------------------------------------------------------
//   Columns []string            — 数据列名列表，如 ["name", "phone", "email"]
//   Data    []map[string]string — 数据行列表。
//                                 map[string]string 是 Go 中的映射类型（类似字典），
//                                 key 是列名，value 是生成的模拟数据值。
//                                 []map[string]string 即"字符串到字符串映射"的切片，
//                                 表示多行数据。
//
// 注意：DC tag 通常用于描述字段含义，这里省略了 dc，因为字段名已经足够说明用途。
// 但在生产项目中建议添加 dc tag 以自动生成更好的 API 文档。
// ============================================================
type MockDataRes struct {
	Columns []string            `json:"columns"`
	Data    []map[string]string `json:"data"`
}
