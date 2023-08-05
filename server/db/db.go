package db

import (
	"context"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/zerozwt/swe"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{Logger: &gormLogger{level: 99}})
	if err != nil {
		return err
	}
	gDB = db
	return gDB.AutoMigrate(models...)
}

func InitMySQL(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: &gormLogger{level: 99}})
	if err != nil {
		return err
	}
	gDB = db
	return gDB.AutoMigrate(models...)
}

func getInstance(ctx *swe.Context) *gorm.DB {
	return newDBSession(ctx, gDB)
}

func newDBSession(ctx *swe.Context, db *gorm.DB) *gorm.DB {
	return db.Session(&gorm.Session{
		Logger:  &gormLogger{level: 99},
		Context: context.WithValue(context.Background(), GORM_CTX_KEY, ctx),
	})
}

type gormCtxKey string

const (
	GORM_CTX_KEY gormCtxKey = "swe_ctx"
)

type gormLogger struct {
	level logger.LogLevel
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLogger{level: level}
}

func (l *gormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.level >= logger.Info {
		l.sweLogger(ctx).Info(format, args...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.level >= logger.Warn {
		l.sweLogger(ctx).Info(format, args...)
	}
}

func (l *gormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.level >= logger.Error {
		l.sweLogger(ctx).Info(format, args...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, ra := fc()
	l.sweLogger(ctx).Debug("gorm trace sql: %s rows affceted: %d err: %v", sql, ra, err)
}

func (l *gormLogger) sweLogger(ctx context.Context) *swe.Logger {
	raw := ctx.Value(GORM_CTX_KEY)
	var tmp *swe.Context = nil
	if raw != nil {
		if v, ok := raw.(*swe.Context); ok {
			tmp = v
		}
	}
	return swe.CtxLogger(tmp)
}
