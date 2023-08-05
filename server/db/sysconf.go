package db

import (
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	DB_SYSCONF_ADMIN_PASS = "admin_pass"
)

type SysConfig struct {
	Key   string `gorm:"type:string;size:1024;column:key;primaryKey"`
	Value string `gorm:"type:string;size:1024;column:value"`
}

func (s SysConfig) TableName() string { return "t_sys_config" }

func init() {
	registerModel(&SysConfig{})
}

type SysConfigDAL struct{}

func GetSysConfigDAL() SysConfigDAL { return SysConfigDAL{} }

func (dal SysConfigDAL) GetConfig(ctx *swe.Context, key string) (string, error) {
	tmp := []SysConfig{}
	err := getInstance(ctx).Where("key = ?", key).Find(&tmp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if err != nil || len(tmp) == 0 {
		return "", nil
	}

	return tmp[0].Value, nil
}

func (dal SysConfigDAL) SetConfig(ctx *swe.Context, key, value string) error {
	tmp := SysConfig{Key: key, Value: value}
	return getInstance(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&tmp).Error
}

func (dal SysConfigDAL) EncodeAdminPassword(pass string) string {
	return utils.EncryptByPass(pass, []byte(pass))
}
