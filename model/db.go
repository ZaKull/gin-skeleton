package model

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/hyperjiang/gin-skeleton/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

// DBInstance is a singleton DB instance
type DBInstance struct {
	initializer func() interface{}
	instance    interface{}
	once        sync.Once
}

var (
	dbInstance *DBInstance
)

// Instance gets the singleton instance
func (i *DBInstance) Instance() interface{} {
	i.once.Do(func() {
		i.instance = i.initializer()
	})
	return i.instance
}

func dbInit() interface{} {
	db, err := gorm.Open(mysql.Open(config.Database.DSN), &gorm.Config{})

	if err != nil {
		glog.Fatalf("Cannot connect to database: %v", err)
	}

	// sql log
	if config.Server.Mode != gin.ReleaseMode {
		db.Logger.LogMode(1)

	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConns)

	return db
}

// DB returns the database instance
func DB() *gorm.DB {
	return dbInstance.Instance().(*gorm.DB)
}

func init() {
	dbInstance = &DBInstance{initializer: dbInit}
}
