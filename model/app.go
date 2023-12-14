package model

const (
	RBACAppConfigTableName = "rbac_app_config"
)

type AppConfig struct {
	TableID
	AppNo     string      `bson:"app_no" gorm:"type:VARCHAR(6);index:app_no,unique;COMMENT:APP编码"`
	Name      string      `bson:"name" gorm:"type:VARCHAR(128);COMMENT:名称"`
	AccessKey string      `bson:"access_key" gorm:"type:VARCHAR(16);COMMENT:AccessKey"`
	SecretKey string      `bson:"secret_key" gorm:"type:VARCHAR(32);COMMENT:SecretKey"`
	Memo      string      `bson:"memo" gorm:"type:VARCHAR(128);COMMENT:备注"`
	Status    LimitStatus `bson:"status" gorm:"type:TINYINT(2);COMMENT:状态 1-正常 2-锁定"`
	TableAt
}

func (*AppConfig) TableName() string {
	return RBACAppConfigTableName
}
