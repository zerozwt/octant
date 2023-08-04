package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GiftRecord struct {
	BatchID    string `gorm:"type:string;size:1024;column:batch_id;primaryKey"`
	RoomID     int64  `gorm:"index:idx_gift_room_time;column:room_id"`
	SendTime   int64  `gorm:"index:idx_gift_room_time;column:send_time"`
	SenderUID  int64  `gorm:"column:sender_uid"`
	SenderName string `gorm:"type:string;size:256;column:sender_name"`
	GiftID     int64  `gorm:"column:gift_id"`
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

func (dal GiftDAL) Page(roomID, tsBegin, tsEnd int64, offset, limit int,
	uid int64, name string, giftID int64) (int, []GiftRecord, error) {
	tx := gDB.Table("t_gift").Where("room_id = ? and (send_time between ? and ?)", roomID, tsBegin, tsEnd)
	if uid > 0 {
		tx = tx.Where("sender_uid = ?", uid)
	}
	if len(name) > 0 {
		tx = tx.Where("sender_name like %?%", name)
	}
	if giftID > 0 {
		tx = tx.Where("gift_id = ?", giftID)
	}

	count := 0
	err := tx.Select("count(*)").Scan(&count).Error
	if err != nil {
		return 0, nil, err
	}

	tx = tx.Offset(offset).Limit(limit).Order("send_time")
	ret := []GiftRecord{}
	err = tx.Find(&ret).Error

	return count, ret, err
}

func (dal GiftDAL) Infos() (ret []GiftInfo, err error) {
	ret = []GiftInfo{}
	err = gDB.Find(&ret).Error
	return
}
