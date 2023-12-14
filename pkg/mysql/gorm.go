package mysql

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewGORMDBWithMysql 初始化 Mysql GORM DB
func NewGORMDBWithMysql(c *Conf) (*gorm.DB, error) {
	sqlDB, err := sql.Open("mysql", c.GenDSN())
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Hour)
	conf := &gorm.Config{}
	if c.Level > 0 {
		conf.Logger = logger.Default.LogMode(logger.LogLevel(c.Level))
	}
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn: sqlDB,
		}), conf,
	)
	return db, err
}

// CloseGORMDB -
func CloseGORMDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("close gorm db error: ", err)
		return err
	}
	return sqlDB.Close()
}
