package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GiftRecord struct {
	BatchID    string `gorm:"type:string;size:1024;column:batch_id;primaryKey"`
	RoomID     int64  `gorm:"index:idx_gift_room_time"`
	SendTime   int64  `gorm:"index:idx_gift_room_time"`
	SenderUID  int64
	SenderName string `gorm:"type:string;size:256"`
	GiftID     int64
	GiftName   string `gorm:"type:string;size:256"`
	GiftPrice  int64
	GiftCount  int64 `gorm:"column:gift_count"`
}

func (s GiftRecord) TableName() string { return "t_gift" }

type GiftInfo struct {
	GiftID    int64  `gorm:"primaryKey;column:gift_id"`
	GiftName  string `gorm:"size:256"`
	GiftPrice int64
}

func (s GiftInfo) TableName() string { return "t_gift_info" }

func init() {
	registerModel(&GiftRecord{})
	registerModel(&GiftInfo{})
}

type GiftDAL struct{}

func GetGiftDAL() GiftDAL { return GiftDAL{} }

func (dal GiftDAL) Insert(gift *GiftRecord) error {
	return gDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "batch_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"gift_count": gorm.Expr("gift_count+?", gift.GiftCount)}),
	}).Create(gift).Error
}

func (dal GiftDAL) UpdateGiftInfo(id int64, name string, price int64) error {
	return gDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&GiftInfo{
		GiftID:    id,
		GiftName:  name,
		GiftPrice: price,
	}).Error
}
