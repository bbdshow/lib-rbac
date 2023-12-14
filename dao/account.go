package dao

import (
	"context"
	"github.com/bbdshow/bkit/caches"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/lib-rbac/model"
	"gorm.io/gorm"

	"time"
)

type AccountOperation interface {
	ListAccount(ctx context.Context, in *model.ListAccountReq) (int64, []*model.Account, error)
	GetAccount(ctx context.Context, in *model.GetAccountReq) (bool, *model.Account, error)
	GetAccountAppActivate(ctx context.Context, in *model.GetAccountAppActivateReq) (bool, *model.AccountAppActivate, error)
	CreateAccount(ctx context.Context, in *model.Account) error
	UpdateAccount(ctx context.Context, in *model.Account, cols []string) error
	CreateAccountAppActivate(ctx context.Context, in *model.AccountAppActivate) error
	UpdateAccountAppActivate(ctx context.Context, in *model.AccountAppActivate, cols []string)
	FindAccount(ctx context.Context, in *model.FindAccountReq) ([]*model.Account, error)
	FindAccountAppActivate(ctx context.Context, in *model.FindAccountAppActivateReq) ([]*model.AccountAppActivate, error)
	DelAccount(ctx context.Context, id int64) error
}

type AccountOperationMysqlImpl struct {
	db       *gorm.DB
	memCache caches.Cacher
}

func (op *AccountOperationMysqlImpl) ListAccount(ctx context.Context, in *model.ListAccountReq) (int64, []*model.Account, error) {
	sess := op.db.WithContext(ctx).Model(&model.Account{})
	if in.Nickname != "" {
		sess.Where("nickname like ?", "%"+in.Nickname+"%")
	}
	if in.Username != "" {
		sess.Where("username like ?", "%"+in.Username+"%")
	}
	if in.Status > 0 {
		sess.Where("status = ?", in.Status)
	}

	var c int64
	if err := sess.Count(&c).Error; err != nil {
		return 0, nil, errc.WithStack(err)
	}

	sess.Order("id DESC").Offset(int(in.Skip())).Limit(in.Size)

	out := make([]*model.Account, 0, in.Size)
	if err := sess.Find(&out).Error; err != nil {
		return 0, nil, errc.WithStack(err)
	}
	return c, out, nil
}

func (op *AccountOperationMysqlImpl) GetAccount(ctx context.Context, in *model.GetAccountReq) (bool, *model.Account, error) {
	key := in.CacheKey()
	if in.UseCache {
		v, err := op.memCache.Get(key)
		if err == nil {
			vv, ok := v.(*model.Account)
			if ok {
				return true, vv, nil
			}
		}
	}

	sess := op.db.WithContext(ctx).Model(&model.Account{})
	if in.OID != "" {
		sess.Where("oid = ?", in.OID)
	}
	if in.Username != "" {
		sess.Where("username = ?", in.Username)
	}

	out := &model.Account{}
	ret := sess.First(&out)
	if ret.Error != nil {
		if IsNotFoundErr(ret.Error) {
			return false, nil, nil
		}
		return false, nil, errc.WithStack(ret.Error)
	}

	if in.UseCache {
		_ = op.memCache.SetWithTTL(key, out, 5*time.Minute)
	}

	return true, out, nil
}

func (op *AccountOperationMysqlImpl) GetAccountAppActivate(ctx context.Context, in *model.GetAccountAppActivateReq) (bool, *model.AccountAppActivate, error) {

	key := in.CacheKey()
	if in.UseCache {
		v, err := op.memCache.Get(key)
		if err == nil {
			vv, ok := v.(*model.AccountAppActivate)
			if ok {
				return true, vv, nil
			}
		}
	}

	sess := op.db.WithContext(ctx).Model(&model.AccountAppActivate{})
	if in.OID != "" {
		sess.Where("oid = ?", in.OID)
	}
	if in.AccountOID != "" {
		sess.Where("account_oid = ?", in.AccountOID)
	}
	if in.AppNo != "" {
		sess.Where("app_no = ?", in.AppNo)
	}
	if in.Token != "" {
		sess.Where("token = ?", in.Token)
	}

	out := &model.AccountAppActivate{}
	ret := sess.First(&out)
	if ret.Error != nil {
		if IsNotFoundErr(ret.Error) {
			return false, nil, nil
		}
		return false, nil, errc.WithStack(ret.Error)
	}

	if in.UseCache {
		_ = op.memCache.SetWithTTL(key, out, 5*time.Minute)
	}

	return true, out, nil
}

func (op *AccountOperationMysqlImpl) CreateAccount(ctx context.Context, in *model.Account) error {
	err := op.db.WithContext(ctx).Create(in).Error
	return errc.WithStack(err)
}

func (op *AccountOperationMysqlImpl) UpdateAccount(ctx context.Context, in *model.Account, cols ...interface{}) error {
	sess := op.db.WithContext(ctx).Model(in)
	if len(cols) > 0 {
		sess.Select(cols[0], cols[1:]...)
	}
	if err := sess.Updates(in).Error; err != nil {
		return errc.WithStack(err)
	}
	return nil
}

func (op *AccountOperationMysqlImpl) CreateAccountAppActivate(ctx context.Context, in *model.AccountAppActivate) error {
	err := op.db.WithContext(ctx).Create(in).Error
	return errc.WithStack(err)
}

func (op *AccountOperationMysqlImpl) UpdateAccountAppActivate(ctx context.Context, in *model.AccountAppActivate, cols ...interface{}) error {
	sess := op.db.WithContext(ctx).Model(in)
	if len(cols) > 0 {
		sess.Select(cols[0], cols[1:]...)
	}
	if err := sess.Updates(in).Error; err != nil {
		return errc.WithStack(err)
	}
	return nil
}

func (op *AccountOperationMysqlImpl) UpsertAccountAppActivateRole(ctx context.Context, in []*model.AccountAppActivate) error {
	if len(in) <= 0 {
		return nil
	}
	err := op.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range in {
			sess := tx.Model(&model.AccountAppActivate{})
			c := int64(0)
			if err := tx.Where("account_oid = ? AND app_oid = ?", v.AccountOID, v.AppOID).Count(&c).Error; err != nil {
				return err
			}
			if c > 0 {
				if err := sess.Select("roles").Updates(v).Error; err != nil {
					return err
				}
				return nil
			}

			if err := sess.Create(in).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return errc.WithStack(err)
}

func (op *AccountOperationMysqlImpl) FindAccount(ctx context.Context, in *model.FindAccountReq) ([]*model.Account, error) {
	sess := op.db.WithContext(ctx).Model(&model.Account{})
	if in.Status > 0 {
		sess.Where("status = ?", in.Status)
	}

	out := make([]*model.Account, 0)
	ret := sess.Find(&out)
	if ret.Error != nil {
		return nil, errc.WithStack(ret.Error)
	}
	return out, nil
}

func (op *AccountOperationMysqlImpl) FindAccountAppActivate(ctx context.Context, in *model.FindAccountAppActivateReq) ([]*model.AccountAppActivate, error) {
	sess := op.db.WithContext(ctx).Model(&model.AccountAppActivate{})
	if in.AccountOID != "" {
		sess.Where("account_oid = ?", in.AccountOID)
	}
	if in.AppOID != "" {
		sess.Where("app_oid = ?", in.AppOID)
	}

	out := make([]*model.AccountAppActivate, 0)
	ret := sess.Find(&out)
	if ret.Error != nil {
		return nil, errc.WithStack(ret.Error)
	}
	return out, nil
}

func (op *AccountOperationMysqlImpl) DelAccount(ctx context.Context, oid int64) error {
	err := op.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("oid = ?", oid).Limit(1).
			Delete(&model.Account{}).Error; err != nil {
			return err
		}
		if err := tx.Where("account_oid = ?", oid).Limit(1).
			Delete(&model.AccountAppActivate{}).Error; err != nil {
			return err
		}
		return nil
	})
	return errc.WithStack(err)
}
