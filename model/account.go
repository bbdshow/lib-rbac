package model

import (
	"fmt"
	"github.com/bbdshow/bkit/gen/str"
	"time"
)

const (
	RBACAccountTableName            = "rbac_account"
	RBACAccountAppActivateTableName = "rbac_account_app_activate"
)

type Account struct {
	TableID
	Nickname  string      `bson:"nickname" gorm:"NOT NULL;type:VARCHAR(64);comment:昵称"`
	Username  string      `bson:"username" gorm:"NOT NULL;type:VARCHAR(64);index:username,unique;comment:账号名"`
	Password  string      `bson:"password" gorm:"NOT NULL;type:VARCHAR(64);comment:密码"`
	Salt      string      `bson:"salt" gorm:"NOT NULL;type:VARCHAR(6);comment:盐"`
	PwdWrong  int         `bson:"pwd_wrong" gorm:"NOT NULL;type:TINYINT(4);comment:密码错误次数"`
	LoginLock int64       `bson:"login_lock" gorm:"NOT NULL;type:BIGINT(20);comment:登录锁定时间"`
	Memo      string      `bson:"memo" gorm:"NOT NULL;type:VARCHAR(128);comment:备注"`
	Status    LimitStatus `bson:"status" gorm:"NOT NULL;type:TINYINT(2);comment:状态 1-正常 2-锁定"`
	TableAt
}

func (*Account) TableName() string {
	return RBACAccountTableName
}

// AccountAppActivate 账户APP激活
type AccountAppActivate struct {
	TableID
	AccountOID   int64       `gorm:"NOT NULL;type:BIGINT(20);index:account_app_oid,unique;comment:账户OID"`
	AppOID       string      `gorm:"NOT NULL;type:VARCHAR(6);index:account_app_oid,unique;comment:APP分组"`
	Token        string      `gorm:"NOT NULL;type:VARCHAR(32);index:token,unique;comment:Token"`
	TokenExpired int64       `gorm:"NOT NULL;type:BIGINT(20);comment:Token过期时间"`
	Roles        OIDSliceStr `gorm:"NOT NULL;type:TEXT;comment:角色ID"`
	TableAt
}

func (*AccountAppActivate) TableName() string {
	return RBACAccountAppActivateTableName
}

func (in *AccountAppActivate) GenToken() string {
	return str.Md5String(str.RandAlphaNumString(12), fmt.Sprintf("%d_%s_%d", in.AccountOID, in.AppOID, time.Now().UnixNano()))
}

func (in *AccountAppActivate) GenTokenExpiredAt() int64 {
	return time.Now().AddDate(0, 0, 1).Unix()
}
