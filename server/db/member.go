package db

import "github.com/zerozwt/swe"

type MembershipRecord struct {
	RoomID     int64  `gorm:"index:idx_member_room_time;column:room_id"`
	SendTime   int64  `gorm:"index:idx_member_room_time;column:send_time"`
	SenderUID  int64  `gorm:"column:sender_uid"`
	SenderName string `gorm:"type:string;size:256;column:sender_name"`
	GuardLevel int    `gorm:"column:level"`
	Count      int
}

func (s MembershipRecord) TableName() string { return "t_member" }

func init() {
	registerModel(&MembershipRecord{})
}

type MemberDAL struct{}

func GetMemberDal() MemberDAL { return MemberDAL{} }

func (dal MemberDAL) Insert(ctx *swe.Context, roomID int64, ts int64, uid int64, name string, level, count int) error {
	item := MembershipRecord{
		RoomID:     roomID,
		SendTime:   ts,
		SenderUID:  uid,
		SenderName: name,
		GuardLevel: level,
		Count:      count,
	}
	return getInstance(ctx).Create(&item).Error
}

func (dal MemberDAL) Page(ctx *swe.Context, roomID, tsBegin, tsEnd int64, offset, limit int,
	uid int64, name string, level []int) (int, []MembershipRecord, error) {
	tx := getInstance(ctx).Table("t_member").Where("room_id = ? and (send_time between ? and ?)", roomID, tsBegin, tsEnd)
	if uid > 0 {
		tx = tx.Where("sender_uid = ?", uid)
	}
	if len(name) > 0 {
		tx = tx.Where("sender_name like %?%", name)
	}
	if len(level) > 0 {
		tx = tx.Where("level in ?", level)
	}

	count := 0
	err := newDBSession(ctx, tx).Select("count(*)").Scan(&count).Error
	if err != nil {
		return 0, nil, err
	}

	tx = tx.Offset(offset).Limit(limit).Order("send_time desc")
	ret := []MembershipRecord{}
	err = tx.Find(&ret).Error

	return count, ret, err
}

func (dal MemberDAL) Range(ctx *swe.Context, roomID, tsBegin, tsEnd int64) ([]*MembershipRecord, error) {
	ret := []*MembershipRecord{}
	tx := getInstance(ctx)
	tx = tx.Where("room_id = ?", roomID)
	tx = tx.Where("send_time between ? and ?", tsBegin, tsEnd)
	tx = tx.Order("send_time")
	err := tx.Find(&ret).Error
	return ret, err
}
