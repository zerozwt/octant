package db

type SuperChatRecord struct {
	RoomID     int64  `gorm:"index:idx_sc_room_time;column:room_id"`
	SendTime   int64  `gorm:"index:idx_sc_room_time;column:send_time"`
	SenderUID  int64  `gorm:"column:sender_uid"`
	SenderName string `gorm:"type:string;size:256;column:sender_name"`
	Price      int64
	Content    string `gorm:"type:string;size:1024;column:content"`
	BgColor    string `gorm:"type:string;size:32"`
	FontColor  string `gorm:"type:string;size:32"`
}

func (s SuperChatRecord) TableName() string { return "t_super_chat" }

func init() {
	registerModel(&SuperChatRecord{})
}

type SCDal struct{}

func GetSCDal() SCDal { return SCDal{} }

func (dal SCDal) Insert(roomID int64, ts int64, uid int64, name string, price int64, content string, bgColor, fontColor string) error {
	sc := SuperChatRecord{
		RoomID:     roomID,
		SendTime:   ts,
		SenderUID:  uid,
		SenderName: name,
		Price:      price,
		Content:    content,
		BgColor:    bgColor,
		FontColor:  fontColor,
	}
	return gDB.Create(&sc).Error
}

func (dal SCDal) Page(roomID, tsBegin, tsEnd int64, offset, limit int,
	uid int64, name, content string) (int, []SuperChatRecord, error) {
	tx := gDB.Table("t_super_chat").Where("room_id = ? and (send_time between ? and ?)", roomID, tsBegin, tsEnd)
	if uid > 0 {
		tx = tx.Where("sender_uid = ?", uid)
	}
	if len(name) > 0 {
		tx = tx.Where("sender_name like %?%", name)
	}
	if len(content) > 0 {
		tx = tx.Where("content like %?%", content)
	}

	count := 0
	err := tx.Select("count(*)").Scan(&count).Error
	if err != nil {
		return 0, nil, err
	}

	tx = tx.Offset(offset).Limit(limit).Order("send_time")
	ret := []SuperChatRecord{}
	err = tx.Find(&ret).Error

	return count, ret, err
}
