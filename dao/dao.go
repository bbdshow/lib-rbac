package dao

import (
	"github.com/bbdshow/bkit/caches"
	"github.com/bbdshow/bkit/db/mysql"
	"github.com/bbdshow/gin-rabc/pkg/conf"
	"github.com/bbdshow/gin-rabc/pkg/model"
	"gorm.io/gorm"
	"strings"
)

// Operation 数据库操作
type Operation interface {
}

type Dao struct {
	cfg *conf.Config

	memCache caches.Cacher
}

func New(cfg *conf.Config) *Dao {
	d := &Dao{
		cfg:   cfg,
		mysql: mysql.NewXorm(cfg.Mysql),

		memCache: caches.NewLRUMemory(10000),
	}
	d.Sync2()

	return d
}

func (d *Dao) Close() {
	if d.mysql != nil {
		_ = d.mysql.Close()
	}
}

func (d *Dao) Sync2() {
	if !d.cfg.Release() {
		err := d.mysql.Sync2(
			new(model.AppConfig),
			new(model.MenuConfig),
			new(model.ActionConfig),
			new(model.RoleConfig),
			new(model.RoleMenuAction),
			new(model.Account),
			new(model.AccountAppActivate),
		)
		if err != nil {
			panic(err)
		}
	}
}

// IsNotFoundErr -
func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == gorm.ErrRecordNotFound.Error() {
		return true
	}
	return false
}

// IsDuplicateErr -
func IsDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "Duplicate") {
		return true
	}
	return false
}
