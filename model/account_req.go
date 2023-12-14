package model

import (
	"fmt"
	"github.com/bbdshow/bkit/typ"
)

type ListAccountReq struct {
	AppNo    string      `json:"app_no" form:"app_no"`
	Nickname string      `json:"nickname" form:"nickname"`
	Username string      `json:"username" form:"username"`
	Status   LimitStatus `json:"status" form:"status"`
	typ.PageReq
}

type ListAccount struct {
	Id        int64       `json:"id"`
	Nickname  string      `json:"nickname"`
	Username  string      `json:"username"`
	PwdWrong  int         `json:"pwdWrong"`
	LoginLock int64       `json:"loginLock"`
	Memo      string      `json:"memo"`
	Status    LimitStatus `json:"status"`
	Roles     []RoleBase  `json:"roles"`
	UpdatedAt int64       `json:"updatedAt"`
	CreatedAt int64       `json:"createdAt"`
}

type RoleBase struct {
	Id      int64       `json:"id"`
	Name    string      `json:"name"`
	AppId   string      `json:"appId"`
	AppName string      `json:"appName"`
	Status  LimitStatus `json:"status"`
}

type GetAccountReq struct {
	OID      string
	Username string
	UseCache bool
}

func (in *GetAccountReq) CacheKey() string {
	return fmt.Sprintf("GetAccountReq.OID%s.Username%s", in.OID, in.Username)
}

type GetAccountAppActivateReq struct {
	OID        string
	AccountOID string
	AppNo      string
	Token      string
	UseCache   bool
}

func (in *GetAccountAppActivateReq) CacheKey() string {
	return fmt.Sprintf("GetAccountAppActivateReq.OID%s.AccountOID%s.AppNo%s.Token%s",
		in.OID, in.AccountOID, in.AppNo, in.Token)
}

type FindAccountReq struct {
	Status LimitStatus
}

type FindAccountAppActivateReq struct {
	AccountOID string
	AppOID     string
}

type CreateAccountReq struct {
	Nickname string `json:"nickname" binding:"required,lte=64"`
	Username string `json:"username" binding:"required,lte=64"`
	Password string `json:"password" binding:"required,len=32"`
	Memo     string `json:"memo" binding:"omitempty,lte=128"`
}

type UpdateAccountReq struct {
	typ.IdReq
	Nickname string      `json:"nickname"`
	Memo     string      `json:"memo"`
	Status   LimitStatus `json:"status"`
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
	OIDReq
	Password string `json:"password" binding:"required,len=32"`
}

type UpdateAccountRoleReq struct {
	OIDReq
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
