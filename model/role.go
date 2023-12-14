package model

const (
	RBACRoleConfigTableName     = "rbac_role_config"
	RBACRoleMenuActionTableName = "rbac_role_menu_action"
)

type RoleConfig struct {
	TableID
	AppId  string      `xorm:"not null VARCHAR(6) comment('APP分组')"`
	Name   string      `xorm:"not null VARCHAR(128) comment('名称')"`
	IsRoot int32       `xorm:"not null TINYINT(2) comment('ROOT 1-ROOT')"`
	Status LimitStatus `xorm:"not null TINYINT(2) comment('状态 1-正常 2-锁定')"`
	Memo   string      `xorm:"not null VARCHAR(128) comment('备注')"`
	TableAt
}

func (*RoleConfig) TableName() string {
	return RBACRoleConfigTableName
}

type RoleMenuAction struct {
	TableID
	RoleId   int64 `xorm:"not null BIGINT(20) comment('角色ID')"`
	MenuId   int64 `xorm:"not null BIGINT(20) comment('主菜单ID')"`
	ActionId int64 `xorm:"not null BIGINT(20) comment('功能ID')"`
}

func (*RoleMenuAction) TableName() string {
	return RBACRoleMenuActionTableName
}
