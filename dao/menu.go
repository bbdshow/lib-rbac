package dao

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/gin-rabc/pkg/model"
	"github.com/bbdshow/gin-rabc/pkg/types"
	"xorm.io/builder"
	"xorm.io/xorm"
)

func (d *Dao) FindMenuConfig(ctx context.Context, in *model.FindMenuConfigReq) ([]*model.MenuConfig, error) {
	conds := make([]builder.Cond, 0)
	if in.ParentId > 0 {
		conds = append(conds, builder.Eq{"parent_id": in.ParentId})
	}
	if len(in.AppId) > 0 {
		conds = append(conds, builder.Eq{"app_id": in.AppId})
	}
	if in.ActionId > 0 {
		conds = append(conds, builder.Like{"actions", "%" + fmt.Sprintf("%d", in.ActionId) + "%"})
	}
	if in.Status > 0 {
		conds = append(conds, builder.Eq{"status": in.Status})
	}

	if len(conds) == 0 {
		return nil, errc.ErrParamInvalid.MultiMsg("condition required")
	}

	sess := d.mysql.Context(ctx).Where("1 = 1")
	for _, c := range conds {
		sess.And(c)
	}

	records := make([]*model.MenuConfig, 0)
	err := sess.Find(&records)
	return records, errc.WithStack(err)
}

func (d *Dao) GetMenuConfig(ctx context.Context, in *model.GetMenuConfigReq) (bool, *model.MenuConfig, error) {
	conds := make([]builder.Cond, 0)
	if in.Id > 0 {
		conds = append(conds, builder.Eq{"id": in.Id})
	}
	if len(in.AppId) > 0 {
		conds = append(conds, builder.Eq{"app_id": in.AppId})
	}

	if len(conds) == 0 {
		return false, nil, errc.ErrParamInvalid.MultiMsg("condition required")
	}
	sess := d.mysql.Context(ctx).Where("1 = 1")
	for _, c := range conds {
		sess.And(c)
	}
	r := &model.MenuConfig{}
	exists, err := sess.Get(r)
	return exists, r, errc.WithStack(err)
}

func (d *Dao) CreateMenuConfig(ctx context.Context, in *model.MenuConfig) error {
	_, err := d.mysql.Context(ctx).InsertOne(in)
	return errc.WithStack(err)
}

func (d *Dao) UpdateMenuConfig(ctx context.Context, in *model.MenuConfig, cols []string) error {
	_, err := d.mysql.Context(ctx).ID(in.Id).Cols(cols...).Update(in)
	return errc.WithStack(err)
}

func (d *Dao) DelMenuConfig(ctx context.Context, id int64) error {
	err := d.mysql.Transaction(func(sess *xorm.Session) error {
		_, err := sess.Context(ctx).ID(id).Delete(&model.MenuConfig{})
		if err != nil {
			return errc.WithStack(err)
		}
		_, err = sess.Context(ctx).Where("menu_id = ?", id).Delete(&model.RoleMenuAction{})
		if err != nil {
			return errc.WithStack(err)
		}
		return nil
	})
	return err
}

func (d *Dao) ListActionConfig(ctx context.Context, in *model.ListActionConfigReq) (int64, []*model.ActionConfig, error) {
	sess := d.mysql.Context(ctx).Where("1 = 1")
	if len(in.Name) > 0 {
		sess.And("name like ?", "%"+in.Name+"%")
	}
	if len(in.Path) > 0 {
		sess.And("path like ?", "%"+in.Path+"%")
	}
	if len(in.Method) > 0 {
		sess.And("method = ?", in.Method)
	}
	if len(in.AppId) > 0 {
		sess.And("app_id = ?", in.AppId)
	}
	if in.Id > 0 {
		sess.And("id = ?", in.Id)
	}

	records := make([]*model.ActionConfig, 0, in.Size)
	c, err := sess.OrderBy("updated_at DESC").Limit(in.LimitStart()).FindAndCount(&records)
	return c, records, errc.WithStack(err)
}

func (d *Dao) FindActionConfig(ctx context.Context, in *model.FindActionConfigReq) ([]*model.ActionConfig, error) {
	sess := d.mysql.Context(ctx).Where("app_id = ?", in.AppId)
	if len(in.ActionId) > 0 {
		sess.In("id", in.ActionId)
	}
	records := make([]*model.ActionConfig, 0)
	err := sess.Find(&records)
	return records, errc.WithStack(err)
}

func (d *Dao) GetActionConfig(ctx context.Context, in *model.GetActionConfigReq) (bool, *model.ActionConfig, error) {
	conds := make([]builder.Cond, 0)
	if in.Id > 0 {
		conds = append(conds, builder.Eq{"id": in.Id})
	}
	if in.AppId != "" {
		conds = append(conds, builder.Eq{"app_id": in.AppId})
	}

	if in.Path != "" {
		conds = append(conds, builder.Eq{"path": in.Path})
	}

	if in.Method != "" {
		conds = append(conds, builder.Eq{"method": in.Method})
	}

	if len(conds) == 0 {
		return false, nil, errc.ErrParamInvalid.MultiMsg("condition required")
	}
	sess := d.mysql.Context(ctx).Where("1 = 1")
	for _, c := range conds {
		sess.And(c)
	}
	r := &model.ActionConfig{}
	exists, err := sess.Get(r)
	return exists, r, errc.WithStack(err)
}

func (d *Dao) CreateActionConfig(ctx context.Context, in *model.ActionConfig) error {
	_, err := d.mysql.Context(ctx).InsertOne(in)
	return errc.WithStack(err)
}

func (d *Dao) UpdateActionConfig(ctx context.Context, in *model.ActionConfig, cols []string) error {
	_, err := d.mysql.Context(ctx).ID(in.Id).Cols(cols...).Update(in)
	return errc.WithStack(err)
}

func (d *Dao) ImportActionConfig(ctx context.Context, in *model.ImportActionConfigReq) error {
	err := d.mysql.Transaction(func(sess *xorm.Session) error {
		ac := &model.ActionConfig{}
		exists, err := sess.Context(ctx).Where("app_id = ?", in.AppId).And("path = ?", in.Path).
			And("method = ?", in.Method).Get(ac)
		if err != nil {
			return errc.WithStack(err)
		}
		if !exists {
			r := &model.ActionConfig{
				AppId:  in.AppId,
				Name:   in.Name,
				Path:   in.Path,
				Method: in.Method,
				Status: types.LimitNormal,
			}
			if _, err := sess.Context(ctx).InsertOne(r); err != nil {
				return errc.WithStack(err)
			}
			return nil
		}

		if in.Name == ac.Name && in.Status == ac.Status {
			return nil
		}
		ac.Name = in.Name
		ac.Status = in.Status
		if _, err := sess.Context(ctx).ID(ac.Id).Cols("name", "status").Update(ac); err != nil {
			return errc.WithStack(err)
		}
		return nil
	})
	return errc.WithStack(err)
}

func (d *Dao) DelActionConfig(ctx context.Context, id int64) error {
	_, err := d.mysql.Context(ctx).ID(id).Delete(&model.ActionConfig{})
	return errc.WithStack(err)
}

func (d *Dao) FindMenuConfigAsRootOrChildren(ctx context.Context, menuId []int64) (root, children []*model.MenuConfig, err error) {
	menus := make([]*model.MenuConfig, 0)
	if err := d.mysql.Context(ctx).In("id", menuId).Find(&menus); err != nil {
		return nil, nil, errc.WithStack(err)
	}
	root = make([]*model.MenuConfig, 0)
	children = make([]*model.MenuConfig, 0)
	for _, v := range menus {
		if v.ParentId <= 0 {
			root = append(root, v)
		} else {
			children = append(children, v)
		}
	}

	for _, v := range children {
		var childrenFindRoot func(ctx context.Context, parentId int64) error
		childrenFindRoot = func(ctx context.Context, parentId int64) error {
			exists, m, err := d.GetMenuConfig(ctx, &model.GetMenuConfigReq{Id: parentId})
			if err != nil {
				return errc.WithStack(err)
			}
			if !exists {
				return nil
			}
			if m.ParentId <= 0 {
				// 如果root不重复在放进去
				isHit := false
				for _, vRoot := range root {
					if vRoot.Id == m.Id {
						isHit = true
						break
					}
				}
				if !isHit {
					root = append(root, m)
				}
				return nil
			} else {
				isHit := false
				for _, vChildren := range children {
					if vChildren.Id == m.Id {
						isHit = true
						break
					}
				}
				if !isHit {
					children = append(children, m)
				}
			}

			return childrenFindRoot(ctx, m.ParentId)
		}

		if err := childrenFindRoot(ctx, v.ParentId); err != nil {
			return nil, nil, err
		}
	}
	return root, children, nil
}
