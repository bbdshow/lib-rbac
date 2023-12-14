package model

import (
	"fmt"
	"github.com/bbdshow/bkit/typ"
)

type ListAppConfigReq struct {
	Name   string `form:"name"`
	Status int    `form:"status"`
	typ.PageReq
}

type ListAppConfig struct {
	OID       string      `json:"oid"`
	AppNo     string      `json:"app_no"`
	Name      string      `json:"name"`       // APP名
	AccessKey string      `json:"access_key"` // 访问KEY
	SecretKey string      `json:"secret_key"` // 加密KEY
	Status    LimitStatus `json:"status"`     // 状态 1-正常 2-限制
	Memo      string      `json:"memo"`       // 备注
	UpdatedAt int64       `json:"updated_at"`
	CreatedAt int64       `json:"created_at"`
}

type SelectAppConfigReq struct {
	Name string `form:"name"`
	typ.PageReq
}

type SelectAppConfig struct {
	OID   string `json:"oid"`
	AppNo string `json:"app_no"`
	Name  string `json:"name"` // APP名
	Memo  string `json:"memo"` // 备注
}

type GetAppConfigReq struct {
	AppNo     string
	AccessKey string
	UseCache  bool
}

func (in *GetAppConfigReq) CacheKey() string {
	return fmt.Sprintf("AppConfig_appNo_%s_accessKey_%s", in.AppNo, in.AccessKey)
}

type GetAppConfigResp struct {
	OID       string      `json:"oid"`
	AppNo     string      `json:"app_no"`
	Name      string      `json:"name"`
	AccessKey string      `json:"access_key"` // 访问KEY
	SecretKey string      `json:"secret_key"` // 加密KEY
	Status    LimitStatus `json:"status"`     // 状态 1-正常 2-限制
	Memo      string      `json:"memo"`       // 备注
}

type CreateAppConfigReq struct {
	Name string `json:"name" binding:"required,gte=1,lte=128"`
	Memo string `json:"memo" binding:"required,lte=128"`
}

type UpdateAppConfigReq struct {
	OIDReq
	IsSecretKey int         `json:"is_secret_key"` // 1 = 重置加密KEY
	Name        string      `json:"name"`
	Memo        string      `json:"memo"`
	Status      LimitStatus `json:"status"` //状态 1-正常 2-限制
}

type DelAppConfigReq struct {
	OIDReq
}
