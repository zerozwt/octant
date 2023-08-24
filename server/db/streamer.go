package db

import (
	"github.com/zerozwt/swe"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Streamer struct {
	RoomID       int64  `gorm:"primaryKey;column:room_id"`
	StreamerName string `gorm:"type:string;size:256"`
	AccountName  string `gorm:"type:string;size:256;index:idx_accname,unique;column:account_name"`
	PrivateKey   string `gorm:"type:string;size:256;column:private_key"`
	PublicKey    string `gorm:"type:string;size:256;column:public_key"`
}

func (s Streamer) TableName() string { return "t_streamer" }

func init() {
	registerModel(&Streamer{})
}

type StreamerDAL struct{}

func GetStreamerDAL() StreamerDAL { return StreamerDAL{} }

func (dal StreamerDAL) All(ctx *swe.Context) ([]Streamer, error) {
	var ret []Streamer
	err := getInstance(ctx).Select("room_id").Find(&ret).Error
	return ret, err
}

func (dal StreamerDAL) Page(ctx *swe.Context, offset, limit int) (int, []Streamer, error) {
	count := 0
	err := getInstance(ctx).Raw("select count(*) from t_streamer").Scan(&count).Error

	if err != nil {
		return count, nil, err
	}

	var ret []Streamer
	err = getInstance(ctx).Offset(offset).Limit(limit).Order("room_id").Find(&ret).Error

	return count, ret, err
}

func (dal StreamerDAL) Insert(ctx *swe.Context, data *Streamer, upsert bool) (int64, error) {
	cc := clause.OnConflict{DoNothing: true}
	if upsert {
		cc = clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"private_key", "public_key"}),
		}
	}
	result := getInstance(ctx).Clauses(cc).Create(data)
	return result.RowsAffected, result.Error
}

func (dal StreamerDAL) Delete(ctx *swe.Context, id int64) (int64, error) {
	result := getInstance(ctx).Where("room_id = ?", id).Delete(&Streamer{})
	return result.RowsAffected, result.Error
}

func (dal StreamerDAL) Find(ctx *swe.Context, id int64) (*Streamer, error) {
	ret := []*Streamer{}

	err := getInstance(ctx).Where("room_id = ?", id).Find(&ret).Error

	if len(ret) > 0 {
		return ret[0], nil
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return nil, err
}

func (dal StreamerDAL) FindByAccount(ctx *swe.Context, account string) (*Streamer, error) {
	ret := []*Streamer{}

	err := getInstance(ctx).Where("account_name = ?", account).Find(&ret).Error

	if len(ret) > 0 {
		return ret[0], nil
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return nil, err
}

func (dal StreamerDAL) UpdatePrivateKey(ctx *swe.Context, id int64, b64EncPriKey string) error {
	return getInstance(ctx).Exec("update t_streamer set private_key = ? where room_id = ?", b64EncPriKey, id).Error
}

func (dal StreamerDAL) FindByRoomIDs(ctx *swe.Context, roomIDs []int64) ([]Streamer, error) {
	ret := []Streamer{}
	err := getInstance(ctx).Where("room_id in ?", roomIDs).Find(&ret).Error
	return ret, err
}
