package model

import (
	"github.com/bbdshow/bkit/typ"
)

type FindMenuConfigReq struct {
	AppId    string
	ParentId int64
	ActionId int64
	Status   LimitStatus
}
type GetMenuConfigReq struct {
	Id    int64
	AppId string
}

type CreateMenuConfigReq struct {
	AppId    string `json:"appId" binding:"required,len=6"`
	Name     string `json:"name" binding:"required,gte=1,lte=128"`
	Memo     string `json:"memo" binding:"omitempty,lte=128"`
	ParentId int64  `json:"parentId" binding:"omitempty,min=0"`
	Sequence int    `json:"sequence" binding:"omitempty,min=0"`
	Path     string `json:"path" binding:"required,gte=1,lte=256"`
	Typ      int    `json:"typ" binding:"required,min=1,max=2"`
}

type UpdateMenuConfigReq struct {
	typ.IdReq
	Name     string      `json:"name" binding:"omitempty,gte=1,lte=128"`
	Memo     string      `json:"memo" binding:"omitempty,lte=128"`
	ParentId int64       `json:"parentId" binding:"omitempty,min=0"`
	Sequence int         `json:"sequence" binding:"omitempty,min=0"`
	Path     string      `json:"path" binding:"omitempty,gte=1,lte=256"`
	Typ      int         `json:"typ" binding:"omitempty,min=1,max=2"`
	Status   LimitStatus `json:"status" binding:"required,min=1,max=2"` // 1-正常 2-锁定
}

type DelActionConfigReq struct {
	typ.IdReq
}
type ImportActionConfigReq struct {
	AppId  string      `json:"appId" binding:"required,len=6"`
	Name   string      `json:"name" binding:"required,gte=1,lte=128"`
	Path   string      `json:"path" binding:"required,gte=1,lte=256"`
	Method string      `json:"method" binding:"required,oneof=GET POST PUT DELETE"` // GET POST PUT DELETE
	Status LimitStatus `json:"status" binding:"required,min=1,max=2"`               // 1-正常 2-锁定
}

type CreateActionConfigReq struct {
	AppId  string `json:"appId" binding:"required,len=6"`
	Name   string `json:"name" binding:"required,gte=1,lte=128"`
	Path   string `json:"path" binding:"required,gte=1,lte=256"`
	Method string `json:"method" binding:"required,oneof=GET POST PUT DELETE"` // GET POST PUT DELETE
}

type UpdateActionConfigReq struct {
	typ.IdReq
	Name   string      `json:"name" binding:"required,gte=1,lte=128"`
	Path   string      `json:"path" binding:"required,gte=1,lte=256"`
	Method string      `json:"method" binding:"required,oneof=GET POST PUT DELETE"` // GET POST PUT DELETE
	Status LimitStatus `json:"status" binding:"required,min=1,max=2"`               // 1-正常 2-锁定
}

type UpdateMenuConfigActionReq struct {
	MenuId   int64   `json:"menuId" binding:"required,gt=0"`
	ActionId []int64 `json:"actionId"`
}

type ListActionConfigReq struct {
	Id     int64  `json:"id" form:"id" binding:"omitempty,gt=0"`
	AppId  string `json:"appId" form:"appId"`
	Name   string `json:"name" form:"name"`
	Path   string `json:"path" form:"path"`
	Method string `json:"method" form:"method"`
	typ.PageReq
}
type ListActionConfigs []*ListActionConfig

func (asc ListActionConfigs) Len() int           { return len(asc) }
func (asc ListActionConfigs) Swap(i, j int)      { asc[i], asc[j] = asc[j], asc[i] }
func (asc ListActionConfigs) Less(i, j int) bool { return asc[i].Name < asc[j].Name }

type ListActionConfig struct {
	Id        int64       `json:"id"`
	AppId     string      `json:"appId"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Method    string      `json:"method"`
	Status    LimitStatus `json:"status"`
	UpdatedAt int64       `json:"updatedAt"`
}

type GetActionConfigReq struct {
	Id     int64
	AppId  string
	Path   string
	Method string
}

type FindActionConfigReq struct {
	AppId    string  `json:"appId" form:"appId" binding:"required,len=6"`
	ActionId []int64 `json:"actionId" form:"actionId"`
}

type FindActionConfigResp struct {
	Actions []*ActionBase `json:"actions"`
}

type ActionBase struct {
	Id     int64       `json:"id"`
	AppId  string      `json:"appId"`
	Name   string      `json:"name"`
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Status LimitStatus `json:"status"`
}

type GetMenuTreeDirsReq struct {
	AppId string `json:"appId" form:"appId" binding:"required,len=6"`
}

type GetMenuTreeDirsResp struct {
	Dirs MenuTreeDirs `json:"dirs"`
}

type GetMenuActionsReq struct {
	MenuId int64 `json:"menuId" form:"menuId" binding:"required,min=1"`
}

type GetMenuActionsResp struct {
	Actions Actions `json:"actions"`
}

type Actions []*Action
type Action struct {
	Id     int64       `json:"id"`
	AppId  string      `json:"appId"`
	Name   string      `json:"name"`
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Status LimitStatus `json:"status"`
}

type MenuTreeDirs []*MenuTreeDir

func (dirs MenuTreeDirs) Len() int           { return len(dirs) }
func (dirs MenuTreeDirs) Swap(i, j int)      { dirs[i], dirs[j] = dirs[j], dirs[i] }
func (dirs MenuTreeDirs) Less(i, j int) bool { return dirs[i].Sequence < dirs[j].Sequence }

type MenuTreeDir struct {
	Id       int64        `json:"id"`
	AppId    string       `json:"appId"`
	Name     string       `json:"name"`
	Typ      int          `json:"typ"`
	Memo     string       `json:"memo"`
	ParentId int64        `json:"parentId"`
	Status   LimitStatus  `json:"status"`
	Sequence int          `json:"sequence"`
	Path     string       `json:"path"`
	Actions  []int64      `json:"actions"`
	Children MenuTreeDirs `json:"children"`
}

type SwaggerJSONToActionsReq struct {
	AppId      string `json:"appId" binding:"required,len=6"`
	Prefix     string `json:"prefix"`
	SwaggerTxt string `json:"swaggerTxt"`
}

type SwaggerJSON struct {
	BasePath string                   `json:"basePath"`
	Paths    map[string]SwaggerMethod `json:"paths"`
}
type SwaggerMethod map[string]struct {
	Summary string `json:"summary"`
}

type DelMenuConfigReq struct {
	typ.IdReq
}
