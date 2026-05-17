package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

type User struct {
	Id        uint64      `json:"id"        orm:"id"        description:"User ID"`
	Username  string      `json:"username"  orm:"username"  description:"Username"`
	Password  string      `json:"password"  orm:"password"  description:"Password"`
	Nickname  string      `json:"nickname"  orm:"nickname"  description:"Nickname"`
	Email     string      `json:"email"     orm:"email"     description:"Email"`
	Phone     string      `json:"phone"     orm:"phone"     description:"Phone"`
	Status    uint        `json:"status"    orm:"status"    description:"Status: 1=active, 0=disabled"`
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" description:"Created at"`
	UpdatedAt *gtime.Time `json:"updated_at" orm:"updated_at" description:"Updated at"`
	DeletedAt *gtime.Time `json:"deleted_at" orm:"deleted_at" description:"Deleted at"`
}
