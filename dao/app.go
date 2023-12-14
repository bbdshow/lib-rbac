package dao

import (
	"context"
	"github.com/bbdshow/bkit/caches"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/lib-rbac/model"
	"gorm.io/gorm"
	"time"
)

type AppOperation interface {
	CreateAppConfig(ctx context.Context, in *model.AppConfig) error
	UpdateAppConfig(ctx context.Context, in *model.AppConfig, cols []string) error
	DelAppConfig(ctx context.Context, id int64) error
	ListAppConfig(ctx context.Context, in *model.ListAppConfigReq) (int64, []*model.AppConfig, error)
	GetAppConfig(ctx context.Context, in *model.GetAppConfigReq) (bool, *model.AppConfig, error)
}

type AppOperationMysqlImpl struct {
	db       *gorm.DB
	memCache caches.Cacher
}

func (op *AppOperationMysqlImpl) CreateAppConfig(ctx context.Context, in *model.AppConfig) error {
	err := op.db.WithContext(ctx).Create(in).Error
	return errc.WithStack(err)
}

func (op *AppOperationMysqlImpl) UpdateAppConfig(ctx context.Context, in *model.AppConfig, cols ...interface{}) error {
	sess := op.db.WithContext(ctx).Model(in)
	if len(cols) > 0 {
		sess.Select(cols[0], cols[1:]...)
	}
	if err := sess.Updates(in).Error; err != nil {
		return errc.WithStack(err)
	}
	return nil
}

func (op *AppOperationMysqlImpl) DelAppConfig(ctx context.Context, oid int64) error {
	sess := op.db.WithContext(ctx).Where("oid = ?", oid)
	if err := sess.Limit(1).Delete(&model.AppConfig{}).Error; err != nil {
		return errc.WithStack(err)
	}
	return nil
}

func (op *AppOperationMysqlImpl) ListAppConfig(ctx context.Context, in *model.ListAppConfigReq) (int64, []*model.AppConfig, error) {
	sess := op.db.WithContext(ctx).Model(&model.AppConfig{})
	if len(in.Name) > 0 {
		sess.Where("name like ?", "%"+in.Name+"%")
	}

	if in.Status > 0 {
		sess.Where("status = ?", in.Status)
	}

	var c int64
	if err := sess.Count(&c).Error; err != nil {
		return 0, nil, errc.WithStack(err)
	}

	sess.Order("id DESC").Offset(int(in.Skip())).Limit(in.Size)

	out := make([]*model.AppConfig, 0, in.Size)
	if err := sess.Find(&out).Error; err != nil {
		return 0, nil, errc.WithStack(err)
	}
	return c, out, nil
}

func (op *AppOperationMysqlImpl) GetAppConfig(ctx context.Context, in *model.GetAppConfigReq) (bool, *model.AppConfig, error) {

	key := in.CacheKey()
	if in.UseCache {
		v, err := op.memCache.Get(key)
		if err == nil {
			c, ok := v.(*model.AppConfig)
			if ok {
				return true, c, nil
			}
		}
	}

	sess := op.db.WithContext(ctx).Model(&model.AppConfig{})
	if in.AppNo != "" {
		sess.Where("app_no = ?", in.AppNo)
	}
	if in.AccessKey != "" {
		sess.Where("access_key = ?", in.AccessKey)
	}

	out := &model.AppConfig{}
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

type AppOperationMongoImpl struct {
	db *mongo.Database
}

//func (op *AppOperationMongoImpl) CreateAppConfig(ctx context.Context, in *model.AppConfig) error {
//	_, err := d.mysql.Context(ctx).InsertOne(in)
//	return errc.WithStack(err)
//}
//
//func (op *AppOperationMongoImpl) UpdateAppConfig(ctx context.Context, in *model.AppConfig, cols []string) error {
//	_, err := d.mysql.Context(ctx).ID(in.Id).Cols(cols...).Update(in)
//	return errc.WithStack(err)
//}
//
//func (op *AppOperationMongoImpl) DelAppConfig(ctx context.Context, id int64) error {
//	_, err := d.mysql.Context(ctx).ID(id).Delete(&model.AppConfig{})
//	return errc.WithStack(err)
//}
//
//func (op *AppOperationMongoImpl) ListAppConfig(ctx context.Context, in *model.ListAppConfigReq) (int64, []*model.AppConfig, error) {
//	sess := d.mysql.Context(ctx).Where("1 = 1")
//	if len(in.Name) > 0 {
//		sess.And("name like ?", "%"+in.Name+"%")
//	}
//
//	if in.Status > 0 {
//		sess.And("status = ?", in.Status)
//	}
//
//	records := make([]*model.AppConfig, 0, in.Size)
//	c, err := sess.OrderBy("id DESC").Limit(in.LimitStart()).FindAndCount(&records)
//	return c, records, errc.WithStack(err)
//}
//
//func (op *AppOperationMongoImpl) GetAppConfig(ctx context.Context, in *model.GetAppConfigReq) (bool, *model.AppConfig, error) {
//	conds := make([]builder.Cond, 0)
//	if len(in.AppId) > 0 {
//		conds = append(conds, builder.Eq{"app_id": in.AppId})
//	}
//	if len(in.AccessKey) > 0 {
//		conds = append(conds, builder.Eq{"access_key": in.AccessKey})
//	}
//	if len(conds) == 0 {
//		return false, nil, errc.ErrParamInvalid.MultiMsg("condition required")
//	}
//	sess := d.mysql.Context(ctx).Where("1 = 1")
//	for _, c := range conds {
//		sess.And(c)
//	}
//
//	r := &model.AppConfig{}
//	exists, err := sess.Get(r)
//	return exists, r, errc.WithStack(err)
//}
//
//func (op *AppOperationMongoImpl) GetAppConfigFromCache(ctx context.Context, in *model.GetAppConfigReq) (bool, *model.AppConfig, error) {
//	key := fmt.Sprintf("AppConfig_appId_%s_accessKey_%s", in.AppId, in.AccessKey)
//
//	v, err := d.memCache.Get(key)
//	if err == nil {
//		c, ok := v.(*model.AppConfig)
//		if ok {
//			return true, c, nil
//		}
//	}
//	exists, c, err := d.GetAppConfig(ctx, in)
//	if err != nil {
//		return false, nil, errc.WithStack(err)
//	}
//	if !exists {
//		return false, nil, nil
//	}
//	_ = d.memCache.SetWithTTL(key, c, 5*time.Minute)
//	return true, c, nil
//}
