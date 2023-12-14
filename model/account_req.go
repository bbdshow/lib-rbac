package model

import (
	"github.com/bbdshow/bkit/typ"
	"github.com/bbdshow/gin-rabc/pkg/types"
)

type ListAccountReq struct {
	AppId    string            `json:"appId" form:"appId"`
	Nickname string            `json:"nickname" form:"nickname"`
	Username string            `json:"username" form:"username"`
	Status   types.LimitStatus `json:"status" form:"status"`
	typ.PageReq
}

type ListAccount struct {
	Id        int64             `json:"id"`
	Nickname  string            `json:"nickname"`
	Username  string            `json:"username"`
	PwdWrong  int               `json:"pwdWrong"`
	LoginLock int64             `json:"loginLock"`
	Memo      string            `json:"memo"`
	Status    types.LimitStatus `json:"status"`
	Roles     []RoleBase        `json:"roles"`
	UpdatedAt int64             `json:"updatedAt"`
	CreatedAt int64             `json:"createdAt"`
}

type RoleBase struct {
	Id      int64             `json:"id"`
	Name    string            `json:"name"`
	AppId   string            `json:"appId"`
	AppName string            `json:"appName"`
	Status  types.LimitStatus `json:"status"`
}

type GetAccountReq struct {
	Id       int64
	Username string
}

type GetAccountAppActivateReq struct {
	Id        int64
	AccountId int64
	AppId     string
	Token     string
}

type FindAccountReq struct {
	Status types.LimitStatus
}

type FindAccountAppActivateReq struct {
	AccountId int64
	AppId     string
}

type CreateAccountReq struct {
	Nickname string `json:"nickname" binding:"required,lte=64"`
	Username string `json:"username" binding:"required,lte=64"`
	Password string `json:"password" binding:"required,len=32"`
	Memo     string `json:"memo" binding:"omitempty,lte=128"`
}

type UpdateAccountReq struct {
	typ.IdReq
	Nickname string            `json:"nickname"`
	Memo     string            `json:"memo"`
	Status   types.LimitStatus `json:"status"`
}

type DelAccountReq struct {
	typ.IdReq
}

type GetAccountMenuAuthReq struct {
	Token string `json:"-"`
}

type GetAccountMenuAuthResp struct {
	IsRoot bool         `json:"isRoot"`
	Dirs   MenuTreeDirs `json:"dirs"`
}

type LoginAccountReq struct {
	AppId    string `json:"appId" binding:"required,len=6"`
	Username string `json:"username" binding:"required,lte=64"`
	Password string `json:"password" binding:"required,len=32"`
}

type LoginAccountResp struct {
	Token        string `json:"token"`
	TokenExpired int64  `json:"tokenExpired"`
	Nickname     string `json:"nickname"`
}

type LoginOutAccountReq struct {
	Token string `json:"-"`
}

type UpdateAccountPasswordReq struct {
	Token       string `json:"-"`
	OldPassword string `json:"oldPassword" binding:"required,len=32"`
	NewPassword string `json:"newPassword" binding:"required,len=32"`
}

type ResetAccountPasswordReq struct {
	typ.IdReq
	Password string `json:"password" binding:"required,len=32"`
}

type UpdateAccountRoleReq struct {
	typ.IdReq
	Roles []int64 `json:"roles" binding:"required"`
}

type VerifyAccountTokenResp struct {
	Verify    bool   // false-验证不通过
	Message   string // 验证不通过原因
	AccountId int64
	Nickname  string
	Username  string
	AppId     string
}
