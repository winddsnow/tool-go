package v1

import "github.com/gogf/gf/v2/frame/g"

type MockDataReq struct {
	g.Meta `path:"/tools/mock-data" method:"post" tags:"Tools" summary:"生成模拟数据"`
	Types  []string `json:"types" v:"required#请选择至少一个数据类型" dc:"数据类型: name,phone,email,id_card,passport,address,ip,datetime,bank_card"`
	Count  int      `json:"count" d:"10" v:"between:1,100#数量范围1-100" dc:"生成数量"`
}

type MockDataRes struct {
	Columns []string            `json:"columns"`
	Data    []map[string]string `json:"data"`
}
