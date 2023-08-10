package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"

	"github.com/zerozwt/swe"
	"gorm.io/gorm/clause"
)

type DDInfo struct {
	UID        int64  `gorm:"primaryKey;column:uid"`
	UserName   string `gorm:"type:string;size:256"`
	AccessCode string `gorm:"type:string;size:64;index:idx_dd_ac,unique;column:access_code"`
	PrivateKey string `gorm:"type:string;size:256;column:private_key"`
	PublicKey  string `gorm:"type:string;size:256;column:public_key"`
}

func (s DDInfo) TableName() string { return "t_dd" }

func init() {
	registerModel(&DDInfo{})

	var buf [16]byte
	rand.Read(buf[:])
	ddDal.enc, _ = aes.NewCipher(buf[:])
}

type DDInfoDAL struct {
	enc cipher.Block
}

var ddDal DDInfoDAL

func GetDDInfoDAL() *DDInfoDAL { return &ddDal }

func (dal DDInfoDAL) BatchCreate(ctx *swe.Context, dd []DDInfo) error {
	return getInstance(ctx).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(dd, 500).Error
}

func (dal *DDInfoDAL) GenerateAccessCode(ts, eventID, uid int64) string {
	var buf [24]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(ts))
	binary.BigEndian.PutUint64(buf[8:], uint64(eventID))
	binary.BigEndian.PutUint64(buf[16:], uint64(uid))
	dal.enc.Encrypt(buf[8:], buf[8:])
	return base64.RawURLEncoding.EncodeToString(buf[:])
}
