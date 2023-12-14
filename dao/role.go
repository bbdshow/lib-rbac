package dao

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/gin-rabc/pkg/model"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
)

func (d *Dao) ListRoleConfig(ctx context.Context, in *model.ListRoleConfigReq) (int64, []*model.RoleConfig, error) {
	sess := d.mysql.Context(ctx).Where("1 = 1")
	if len(in.Name) > 0 {
		sess.And("name like ?", "%"+in.Name+"%")
	}
	if len(in.AppId) > 0 {
		sess.And("app_id = ?", in.AppId)
	}
	records := make([]*model.RoleConfig, 0, in.Size)
	c, err := sess.OrderBy("id DESC").Limit(in.LimitStart()).FindAndCount(&records)
	return c, records, errc.WithStack(err)
}

func (d *Dao) FindRoleConfig(ctx context.Context, in *model.FindRoleConfigReq) ([]*model.RoleConfig, error) {
	conds := make([]builder.Cond, 0)
	if len(in.RoleId) > 0 {
		conds = append(conds, builder.In("id", in.RoleId))
	}

	if len(conds) == 0 {
		return nil, errc.ErrParamInvalid.MultiMsg("condition required")
	}

	sess := d.mysql.Context(ctx).Where("1 = 1")
	for _, c := range conds {
		sess.And(c)
	}
	records := make([]*model.RoleConfig, 0)
	err := sess.Find(&records)
	return records, errc.WithStack(err)
}

func (d *Dao) GroupRolesMenuId(ctx context.Context, roleId []int64) ([]int64, error) {
	records := make([]*model.RoleMenuAction, 0)
	err := d.mysql.Context(ctx).In("role_id", roleId).GroupBy("menu_id").Find(&records)
	if err != nil {
		return nil, errc.WithStack(err)
	}
	menuId := make([]int64, 0, len(records))
	for _, v := range records {
		menuId = append(menuId, v.MenuId)
	}
	return menuId, nil
}

func (d *Dao) GetRoleConfig(ctx context.Context, in *model.GetRoleConfigReq) (bool, *model.RoleConfig, error) {
	conds := make([]builder.Cond, 0)
	if in.Id > 0 {
		conds = append(conds, builder.Eq{"id": in.Id})
	}
	if in.Name != "" {
		conds = append(conds, builder.Eq{"name": in.Name})
	}
	if len(conds) == 0 {
		return false, nil, errc.ErrParamInvalid.MultiMsg("condition required")
	}
	sess := d.mysql.Context(ctx).Where("1 = 1")
	for _, c := range conds {
		sess.And(c)
	}

	r := &model.RoleConfig{}
	exists, err := sess.Get(r)
	return exists, r, errc.WithStack(err)
}

func (d *Dao) CreateRoleConfig(ctx context.Context, in *model.RoleConfig) error {
	_, err := d.mysql.Context(ctx).InsertOne(in)
	return errc.WithStack(err)
}

func (d *Dao) UpdateRoleConfig(ctx context.Context, in *model.RoleConfig, cols []string) error {
	_, err := d.mysql.Context(ctx).ID(in.Id).Cols(cols...).Update(in)
	return errc.WithStack(err)
}

func (d *Dao) DelRoleConfig(ctx context.Context, id int64) error {
	err := d.mysql.Transaction(func(sess *xorm.Session) error {
		_, err := sess.Context(ctx).ID(id).Delete(&model.RoleConfig{})
		if err != nil {
			return errc.WithStack(err)
		}
		_, err = sess.Context(ctx).Where("role_id = ?", id).Delete(&model.RoleMenuAction{})
		if err != nil {
			return errc.WithStack(err)
		}
		return nil
	})
	return err
}

func (d *Dao) GetRoleConfigFromCache(ctx context.Context, in *model.GetRoleConfigReq) (bool, *model.RoleConfig, error) {
	key := fmt.Sprintf("RoleConfig_id_%d", in.Id)
	v, err := d.memCache.Get(key)
	if err == nil {
		c, ok := v.(*model.RoleConfig)
		if ok {
			return true, c, nil
		}
	}
	exists, c, err := d.GetRoleConfig(ctx, in)
	if err != nil {
		return false, nil, errc.WithStack(err)
	}
	if !exists {
		return false, nil, nil
	}
	_ = d.memCache.SetWithTTL(key, c, 5*time.Minute)
	return true, c, nil
}

func (d *Dao) FindRoleAllMenuAction(ctx context.Context, roleId int64) ([]*model.RoleMenuAction, error) {
	records := make([]*model.RoleMenuAction, 0)
	err := d.mysql.Context(ctx).Where("role_id = ?", roleId).Find(&records)
	return records, errc.WithStack(err)
}

func (d *Dao) UpdateRoleMenuAction(ctx context.Context, add []*model.RoleMenuAction, del []int64) error {
	err := d.mysql.Transaction(func(sess *xorm.Session) error {
		if len(add) > 0 {
			_, err := sess.Context(ctx).Insert(add)
			if err != nil {
				return errc.WithStack(err)
			}
		}

		if len(del) > 0 {
			_, err := sess.Context(ctx).In("id", del).Delete(&model.RoleMenuAction{})
			if err != nil {
				return errc.WithStack(err)
			}
		}
		return nil
	})

	return errc.WithStack(err)
}

func (d *Dao) FindAllRole(ctx context.Context) (model.Roles, error) {
	roles := make(model.Roles, 0)
	roleRecords := make([]*model.RoleConfig, 0)
	if err := d.mysql.Context(ctx).Where("status = 1").Find(&roleRecords); err != nil {
		return nil, errc.WithStack(err)
	}
	actionsMap := map[string][]*model.ActionConfig{}

	for _, v := range roleRecords {
		_, ok := actionsMap[v.AppId]
		if !ok {
			actions, err := d.FindActionConfig(ctx, &model.FindActionConfigReq{
				AppId: v.AppId,
			})
			if err != nil {
				return nil, err
			}
			actionsMap[v.AppId] = actions
		}

		role := model.Role{
			RoleId:  v.Id,
			Actions: make(model.Actions, 0),
		}

		if v.IsRoot == 1 {
			// 如果是超级角色，则角色拥有此 appId 下的所有功能权限
			actions, ok := actionsMap[v.AppId]
			if !ok {
				continue
			}
			for _, act := range actions {
				role.Actions = append(role.Actions, &model.Action{
					Id:     act.Id,
					AppId:  act.AppId,
					Name:   act.Name,
					Path:   act.Path,
					Method: act.Method,
					Status: act.Status,
				})
			}
			roles = append(roles, role)
			continue
		}

		// 普通用户
		menuActions, err := d.FindRoleAllMenuAction(ctx, v.Id)
		if err != nil {
			return nil, err
		}

		actions, ok := actionsMap[v.AppId]
		if !ok {
			continue
		}
		for _, ma := range menuActions {
			for _, act := range actions {
				if ma.ActionId == act.Id {
					role.Actions = append(role.Actions, &model.Action{
						Id:     act.Id,
						AppId:  act.AppId,
						Name:   act.Name,
						Path:   act.Path,
						Method: act.Method,
						Status: act.Status,
					})
					break
				}
			}
		}
		roles = append(roles, role)
	}
	return roles, nil
}
