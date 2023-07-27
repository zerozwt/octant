package db

import (
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gDB *gorm.DB = nil

var models []any = []any{}
var modelsLock sync.Mutex

func registerModel(model any) {
	modelsLock.Lock()
	defer modelsLock.Unlock()
	models = append(models, model)
}

func InitSQLite(file string) error {
	db, err := gorm.Open(sqlite.Open(file))
	if err != nil {
		return err
	}
	gDB = db
	return gDB.AutoMigrate(models...)
}

func InitMySQL(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}
	gDB = db
	return gDB.AutoMigrate(models...)
}
