package db

type SuperChatRecord struct {
	RoomID     int64 `gorm:"index:idx_sc_room_time"`
	SendTime   int64 `gorm:"index:idx_sc_room_time"`
	SenderUID  int64
	SenderName string `gorm:"type:string;size:256"`
	Price      int64
	Content    string `gorm:"type:string;size:1024"`
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
