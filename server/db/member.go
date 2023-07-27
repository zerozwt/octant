package db

type MembershipRecord struct {
	RoomID     int64 `gorm:"index:idx_member_room_time"`
	SendTime   int64 `gorm:"index:idx_member_room_time"`
	SenderUID  int64
	SenderName string `gorm:"type:string;size:256"`
	GuardLevel int
	Count      int
}

func (s MembershipRecord) TableName() string { return "t_member" }

func init() {
	registerModel(&MembershipRecord{})
}

type MemberDAL struct{}

func GetMemberDal() MemberDAL { return MemberDAL{} }

func (dal MemberDAL) Insert(roomID int64, ts int64, uid int64, name string, level, count int) error {
	item := MembershipRecord{
		RoomID:     roomID,
		SendTime:   ts,
		SenderUID:  uid,
		SenderName: name,
		GuardLevel: level,
		Count:      count,
	}
	return gDB.Create(&item).Error
}
