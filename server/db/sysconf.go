package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (dal SysConfigDAL) GetConfig(key string) (string, error) {
	tmp := []SysConfig{}
	err := gDB.Where("key = ?", key).Find(&tmp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}

	if err != nil || len(tmp) == 0 {
		return "", nil
	}

	return tmp[0].Value, nil
}

func (dal SysConfigDAL) SetConfig(key, value string) error {
	tmp := SysConfig{Key: key, Value: value}
	return gDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&tmp).Error
}

func (dal SysConfigDAL) EncodeAdminPassword(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	block, _ := aes.NewCipher(hash[:])
	aead, _ := cipher.NewGCM(block)
	enc := aead.Seal(nil, []byte(`123456789012`), []byte(pass), nil)
	return hex.EncodeToString(enc)
}
