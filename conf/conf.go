package conf

import (
	"github.com/bbdshow/bkit/conf"
	"github.com/bbdshow/bkit/gen/defval"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/lib-rbac/pkg/mysql"
)

var (
	Conf = &Config{}
)

type Config struct {
	Env    string `defval:"dev"`
	Admin  Admin
	Casbin Casbin
	Mysql  mysql.Conf
	Logger logs.Config
}

func InitConf(path ...string) error {
	return conf.ReadConfig(conf.FlagConfigPath(path...), Conf)
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) Release() bool {
	return c.Env == "release"
}

func (c *Config) EraseSensitive() Config {
	// 脱敏数据，可用于打印
	cc := *c
	_ = defval.InitialNullVal(&cc)
	return cc
}

type Admin struct {
	HttpListenAddr string `defval:"0.0.0.0:49000"`
	AuthEnable     bool   `defval:"true"`
	SignEnable     bool   `defval:"true"`
}

type Casbin struct {
	Enable          bool `defval:"true"`
	Debug           bool `defval:"true"`
	AutoIntervalSec int  `defval:"60"`
}
