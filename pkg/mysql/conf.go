// Package orm provides gorm
package mysql

import "fmt"

// Conf -
type Conf struct {
	Database    string `yaml:"database"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
	Level       int    `yaml:"level"` // Silent = 1, Error = 2, Warn = 3, Info = 4
}

// Validate -
func (c *Conf) Validate() error {
	if c.Database == "" {
		return fmt.Errorf("database required")
	}
	if c.Host == "" {
		return fmt.Errorf("host required")
	}
	if c.Port == "" {
		return fmt.Errorf("port required")
	}
	if c.User == "" {
		return fmt.Errorf("user required")
	}
	if c.Password == "" {
		return fmt.Errorf("password required")
	}
	return nil
}

// GenDSN -
func (c *Conf) GenDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
}
