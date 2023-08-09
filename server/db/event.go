package db

import (
	"github.com/zerozwt/swe"
	"gorm.io/gorm/clause"
)

type RewardEvent struct {
	ID            int64  `gorm:"primaryKey;column:id"`
	RoomID        int64  `gorm:"index:idx_re_room;column:room_id"`
	EventName     string `gorm:"type:string;size:256;column:name"`
	RewardContent string `gorm:"type:TEXT"`
	Conditions    string `gorm:"type:TEXT"`
	CreateTime    int64  `gorm:"index:idx_re_room;column:create_time"`
	Status        int    `gorm:"column:status"`
}

func (s RewardEvent) TableName() string { return "t_event" }

type RewardUser struct {
	EventID     int64  `gorm:"index:idx_ru_evt;column:event_id"`
	UID         int64  `gorm:"index:idx_ru_uid;column:uid"`
	UserName    string `gorm:"type:string;size:256"`
	Time        int64  `gorm:"column:ts"`
	Blocked     int    `gorm:"column:block"`
	Columns     string `gorm:"type:TEXT"`
	AddressInfo string `gorm:"type:string;size:4096"`
}

func (s RewardUser) TableName() string { return "t_event_user" }

const (
	EVENT_IDLE = iota + 1
	EVENT_CALCULATING
	EVENT_ERROR
	EVENT_READY
)

func init() {
	registerModel(&RewardEvent{})
	registerModel(&RewardUser{})
}

type rewardEventDAL struct{}

func GetRewardEventDAL() rewardEventDAL { return rewardEventDAL{} }

func (dal rewardEventDAL) Page(ctx *swe.Context, roomID int64, offset, limit int) (int, []RewardEvent, error) {
	ret := []RewardEvent{}
	count := 0

	tx := getInstance(ctx).Table("t_event").Where("room_id = ?", roomID)
	err := newDBSession(ctx, tx).Select("count(*)").Scan(&count).Error
	if err != nil {
		return 0, nil, err
	}

	tx = tx.Offset(offset).Limit(limit)
	tx = tx.Select("id", "name", "status").Order("create_time desc")
	err = tx.Find(&ret).Error
	return count, ret, err
}

func (dal rewardEventDAL) Put(ctx *swe.Context, item *RewardEvent) error {
	return getInstance(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error
}

func (dal rewardEventDAL) Get(ctx *swe.Context, id int64) (*RewardEvent, error) {
	ret := []RewardEvent{}
	err := getInstance(ctx).Where("id = ?", id).Find(&ret).Error
	if len(ret) == 0 {
		return nil, err
	}
	return &ret[0], nil
}

func (dal rewardEventDAL) SetStatus(ctx *swe.Context, id int64, status int) error {
	return getInstance(ctx).Exec("update t_event set status = ? where id = ?", status, id).Error
}

func (dal rewardEventDAL) ClearUsers(ctx *swe.Context, eventID int64) error {
	return getInstance(ctx).Exec("delete from t_event_user where event_id = ?", eventID).Error
}

func (dal rewardEventDAL) PutUsers(ctx *swe.Context, users []RewardUser) error {
	return getInstance(ctx).CreateInBatches(users, 500).Error
}
