package db

import (
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

func (dal StreamerDAL) All() ([]Streamer, error) {
	var ret []Streamer
	err := gDB.Select("room_id").Find(&ret).Error
	return ret, err
}

func (dal StreamerDAL) Page(offset, limit int) (int, []Streamer, error) {
	count := 0
	err := gDB.Raw("select count(*) from t_streamer").Scan(&count).Error

	if err != nil {
		return count, nil, err
	}

	var ret []Streamer
	err = gDB.Offset(offset).Limit(limit).Order("room_id").Find(&ret).Error

	return count, ret, err
}

func (dal StreamerDAL) Insert(data *Streamer, upsert bool) (int64, error) {
	cc := clause.OnConflict{DoNothing: true}
	if upsert {
		cc = clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}, {Name: "account_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"private_key", "public_key"}),
		}
	}
	result := gDB.Clauses(cc).Create(data)
	return result.RowsAffected, result.Error
}

func (dal StreamerDAL) Delete(id int64) (int64, error) {
	result := gDB.Where("room_id = ?", id).Delete(&Streamer{})
	return result.RowsAffected, result.Error
}

func (dal StreamerDAL) Find(id int64) (*Streamer, error) {
	ret := []*Streamer{}

	err := gDB.Where("room_id = ?", id).Find(&ret).Error

	if len(ret) > 0 {
		return ret[0], nil
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return nil, err
}

func (dal StreamerDAL) FindByAccount(account string) (*Streamer, error) {
	ret := []*Streamer{}

	err := gDB.Where("account_name = ?", account).Find(&ret).Error

	if len(ret) > 0 {
		return ret[0], nil
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	return nil, err
}

func (dal StreamerDAL) UpdatePrivateKey(id int64, b64EncPriKey string) error {
	return gDB.Exec("update t_streamer set private_key = ? where room_id = ?", b64EncPriKey, id).Error
}
