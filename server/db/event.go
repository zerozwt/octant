package db

import (
	"github.com/zerozwt/swe"
	"gorm.io/gorm/clause"
)

type RewardEvent struct {
	ID            int64  `gorm:"primaryKey;column:id"`
	RoomID        int64  `gorm:"index:idx_re_room;column:room_id"`
	EventName     string `gorm:"type:string;size:256;column:name"`
	RewardContent string `gorm:"type:TEXT;column:content"`
	Conditions    string `gorm:"type:TEXT;column:conditions"`
	Hidden        int    `gorm:"column:hidden"`
	CreateTime    int64  `gorm:"index:idx_re_room;column:create_time"`
	Status        int    `gorm:"column:status"`
}

func (s RewardEvent) TableName() string { return "t_event" }

type RewardUser struct {
	EventID     int64  `gorm:"index:idx_ru_evt;column:event_id"`
	UID         int64  `gorm:"index:idx_ru_uid;column:uid"`
	UserName    string `gorm:"type:string;size:256;column:user_name"`
	Time        int64  `gorm:"column:ts"`
	Blocked     int    `gorm:"column:block"`
	Columns     string `gorm:"type:TEXT;column:cols"`
	AddressInfo string `gorm:"type:string;size:4096;column:address_info"`
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

type RewardEventDAL struct{}

func GetRewardEventDAL() RewardEventDAL { return RewardEventDAL{} }

func (dal RewardEventDAL) Page(ctx *swe.Context, roomID int64, offset, limit int) (int, []RewardEvent, error) {
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

func (dal RewardEventDAL) Put(ctx *swe.Context, item *RewardEvent) error {
	return getInstance(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(item).Error
}

func (dal RewardEventDAL) Get(ctx *swe.Context, id int64) (*RewardEvent, error) {
	ret := []RewardEvent{}
	err := getInstance(ctx).Where("id = ?", id).Find(&ret).Error
	if len(ret) == 0 {
		return nil, err
	}
	return &ret[0], nil
}

func (dal RewardEventDAL) SetStatus(ctx *swe.Context, id int64, status int) error {
	return getInstance(ctx).Exec("update t_event set status = ? where id = ?", status, id).Error
}

func (dal RewardEventDAL) ClearUsers(ctx *swe.Context, eventID int64) error {
	return getInstance(ctx).Exec("delete from t_event_user where event_id = ?", eventID).Error
}

func (dal RewardEventDAL) PutUsers(ctx *swe.Context, users []RewardUser) error {
	return getInstance(ctx).CreateInBatches(users, 500).Error
}

func (dal RewardEventDAL) UpdateEventInfo(ctx *swe.Context, id, roomID int64, name, content string, hidden int) error {
	return getInstance(ctx).Exec("update t_event set name = ? and content = ? where id = ? and room_id = ? and hidden = ?",
		name, content, id, roomID, hidden).Error
}

func (dal RewardEventDAL) GetByRoomID(ctx *swe.Context, id, roomID int64) (*RewardEvent, error) {
	ret := []RewardEvent{}
	err := getInstance(ctx).Where("id = ? and room_id = ?", id, roomID).Find(&ret).Error
	if len(ret) == 0 {
		return nil, err
	}
	return &ret[0], nil
}

func (dal RewardEventDAL) Delete(ctx *swe.Context, id, roomID int64) (int64, error) {
	result := getInstance(ctx).Exec("delete from t_event where id = ? and room_id = ?", id, roomID)
	return result.RowsAffected, result.Error
}

func (dal RewardEventDAL) Exist(ctx *swe.Context, id, roomID int64) (bool, error) {
	ret := []RewardEvent{}
	err := getInstance(ctx).Where("id = ? and room_id = ?", id, roomID).Select("id", "name").Find(&ret).Error

	return len(ret) > 0, err
}

func (dal RewardEventDAL) UserPage(ctx *swe.Context, eventID int64, offset, limit int) (int, []RewardUser, error) {
	ret := []RewardUser{}
	tx := getInstance(ctx)
	tx = tx.Where("event_id = ?", eventID)

	count := 0
	if err := newDBSession(ctx, tx).Select("count(*)").Scan(&count).Error; err != nil {
		return 0, nil, err
	}

	tx = tx.Offset(offset).Limit(limit).Order("ts")
	err := tx.Find(&ret).Error
	return count, ret, err
}

func (dal RewardEventDAL) BlockUser(ctx *swe.Context, eventID, uid int64, block bool) (bool, error) {
	blockValue := 0
	if block {
		blockValue = 1
	}
	result := getInstance(ctx).Exec("update t_event_user set block = ? where event_id = ? and uid = ?", blockValue, eventID, uid)
	return result.RowsAffected > 0, result.Error
}

func (dal RewardEventDAL) Users(ctx *swe.Context, eventID int64) ([]RewardUser, error) {
	ret := []RewardUser{}
	tx := getInstance(ctx)
	tx = tx.Where("event_id = ?", eventID)
	err := tx.Find(&ret).Error
	return ret, err
}

func (dal RewardEventDAL) UserRecords(ctx *swe.Context, uid int64, offset, limit int) (int, []RewardUser, error) {
	ret := []RewardUser{}
	tx := getInstance(ctx).Where("uid = ?", uid)

	count := 0
	err := newDBSession(ctx, tx).Select("count(*)").Scan(&count).Error
	if err != nil {
		return 0, nil, err
	}

	tx = tx.Offset(offset).Limit(limit).Order("ts desc")
	err = tx.Find(&ret).Error
	return count, ret, err
}

func (dal RewardEventDAL) GetByIDs(ctx *swe.Context, eventIDs []int64) ([]RewardEvent, error) {
	ret := []RewardEvent{}
	tx := getInstance(ctx).Where("id in ?", eventIDs)
	tx = tx.Select("id", "room_id", "name", "content")
	err := tx.Find(&ret).Error
	return ret, err
}

func (dal RewardEventDAL) UserInfoByUIDAndEventID(ctx *swe.Context, uid, eventID int64) ([]RewardUser, error) {
	ret := []RewardUser{}
	err := getInstance(ctx).Where("uid = ? and event_id = ?", uid, eventID).Find(&ret).Error
	return ret, err
}

func (dal RewardEventDAL) UpdateAddrInfo(ctx *swe.Context, uid, eventID int64, addr string) (int64, error) {
	tx := getInstance(ctx).Exec("update t_event_user set address_info = ? where uid = ? and event_id = ?",
		addr, uid, eventID)
	return tx.RowsAffected, tx.Error
}
