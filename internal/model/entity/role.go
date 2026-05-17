package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type Role struct {
	Id        uint64      `json:"id"        orm:"id"        description:"Role ID"`
	Name      string      `json:"name"      orm:"name"      description:"Role name"`
	Code      string      `json:"code"      orm:"code"      description:"Role code"`
	Sort      int         `json:"sort"      orm:"sort"      description:"Sort order"`
	Status    uint        `json:"status"    orm:"status"    description:"Status: 1=active, 0=disabled"`
	Desc      string      `json:"desc"      orm:"desc"      description:"Description"`
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" description:"Created at"`
	UpdatedAt *gtime.Time `json:"updated_at" orm:"updated_at" description:"Updated at"`
	DeletedAt *gtime.Time `json:"deleted_at" orm:"deleted_at" description:"Deleted at"`
}
