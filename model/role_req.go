package model

import (
	"github.com/bbdshow/bkit/typ"
)

type ListRoleConfigReq struct {
	AppId string `json:"appId" form:"appId"`
	Name  string `json:"name" form:"name"`
	typ.PageReq
}

type ListRoleConfig struct {
	Id        int64       `json:"id"`
	AppId     string      `json:"appId"`
	AppName   string      `json:"appName"`
	Name      string      `json:"name"`
	IsRoot    int32       `json:"isRoot"`
	Memo      string      `json:"memo"`
	Status    LimitStatus `json:"status"`
	UpdatedAt int64       `json:"updatedAt"`
}

type FindRoleConfigReq struct {
	RoleId []int64
}

type GetRoleConfigReq struct {
	Id   int64 `json:"id" form:"id"`
	Name string
}

type GetRoleMenuActionReq struct {
	RoleId int64 `json:"roleId" form:"roleId" binding:"required,min=1"`
}

type GetRoleMenuActionResp struct {
	MenuActions []MenuAction `json:"menuActions"`
}

type CreateRoleConfigReq struct {
	AppId  string `json:"appId" binding:"required,len=6"`
	Name   string `json:"name" binding:"required,gte=1,lte=128"`
	IsRoot int32  `json:"isRoot"`
	Memo   string `json:"memo"`
}

type UpdateRoleConfigReq struct {
	typ.IdReq
	Name   string      `json:"name"`
	IsRoot int32       `json:"isRoot"`
	Memo   string      `json:"memo"`
	Status LimitStatus `json:"status"`
}

type DelRoleConfigReq struct {
	typ.IdReq
}

type UpsertRoleMenuActionReq struct {
	RoleId      int64 `json:"roleId" form:"roleId" binding:"required,min=1"`
	MenuActions []MenuAction
}

type MenuAction struct {
	MenuId  int64   `json:"menuId" binding:"required,min=1"`
	Actions []int64 `json:"actions"`
}

type Roles []Role

type Role struct {
	RoleId  int64
	Actions Actions
}
