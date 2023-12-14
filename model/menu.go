package model

const (
	RBACMenuConfigTableName   = "rbac_menu_config"
	RBACActionConfigTableName = "rbac_action_config"
)

type MenuConfig struct {
	TableID
	AppOID   string      `gorm:"NOT NULL;type:VARCHAR(6);index:account_app_oid,unique;comment:APP分组"`
	Name     string      `gorm:"NOT NULL;type:VARCHAR(128);comment:名称"`
	Typ      int         `gorm:"NOT NULL;type:TINYINT(2);comment:分类 1-菜单 2-分组"`
	Memo     string      `gorm:"NOT NULL;type:VARCHAR(128);comment:备注"`
	ParentId int64       `gorm:"NOT NULL;type:BIGINT(20);comment:父ID"`
	Status   LimitStatus `gorm:"NOT NULL;type:TINYINT(2);comment:状态 1-正常 2-锁定"`
	Sequence int         `gorm:"NOT NULL;type:INT(11);default:0;comment:序号"`
	Path     string      `gorm:"NOT NULL;type:VARCHAR(255);comment:路径"`
	Actions  OIDSliceStr `gorm:"NOT NULL;type:TEXT;comment:功能ID"`
	TableAt
}

func (*MenuConfig) TableName() string {
	return RBACMenuConfigTableName
}

type ActionConfig struct {
	TableID
	AppOID string      `gorm:"NOT NULL;type:VARCHAR(6);index:app_oid_path_method,priority:1,unique;comment:APP分组"`
	Name   string      `xorm:"NOT NULL;type:VARCHAR(128);index:name;comment:名称"`
	Path   string      `xorm:"NOT NULL;type:VARCHAR(255);index:app_oid_path_method,priority:1,unique;comment:访问路径"`
	Method string      `xorm:"NOT NULL;type:VARCHAR(10);index:app_oid_path_method,priority:1,unique;comment:GET POST PUT DELETE"`
	Status LimitStatus `xorm:"NOT NULL;type:TINYINT(2);comment:状态 1-正常 2-锁定"`
	TableAt
}

func (*ActionConfig) TableName() string {
	return RBACMenuConfigTableName
}
