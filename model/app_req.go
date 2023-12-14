package model

import (
	"github.com/bbdshow/bkit/typ"
)

type ListAppConfigReq struct {
	Name   string `form:"name"`
	Status int    `form:"status"`
	typ.PageReq
}

type ListAppConfig struct {
	Id        int64       `json:"id"`
	AppId     string      `json:"appId"`
	Name      string      `json:"name"`      // APP名
	AccessKey string      `json:"accessKey"` // 访问KEY
	SecretKey string      `json:"secretKey"` // 加密KEY
	Status    LimitStatus `json:"status"`    // 状态 1-正常 2-限制
	Memo      string      `json:"memo"`      // 备注
	UpdatedAt int64       `json:"updatedAt"`
	CreatedAt int64       `json:"createdAt"`
}

type SelectAppConfigReq struct {
	Name string `form:"name"`
	typ.PageReq
}

type SelectAppConfig struct {
	Id    int64  `json:"id"`
	AppId string `json:"appId"`
	Name  string `json:"name"` // APP名
	Memo  string `json:"memo"` // 备注
}

type GetAppConfigReq struct {
	AppId     string
	AccessKey string
}

type GetAppConfigResp struct {
	Id        int64       `json:"id"`
	AppId     string      `json:"appId"`
	Name      string      `json:"name"`
	AccessKey string      `json:"accessKey"` // 访问KEY
	SecretKey string      `json:"secretKey"` // 加密KEY
	Status    LimitStatus `json:"status"`    // 状态 1-正常 2-限制
	Memo      string      `json:"memo"`      // 备注
}

type CreateAppConfigReq struct {
	Name string `json:"name" binding:"required,gte=1,lte=128"`
	Memo string `json:"memo" binding:"required,lte=128"`
}

type UpdateAppConfigReq struct {
	typ.IdReq
	IsSecretKey int         `json:"isSecretKey"` // 1 = 重置加密KEY
	Name        string      `json:"name"`
	Memo        string      `json:"memo"`
	Status      LimitStatus `json:"status"` //状态 1-正常 2-限制
}

type DelAppConfigReq struct {
	typ.IdReq
}
