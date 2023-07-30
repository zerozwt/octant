package db

type Streamer struct {
	RoomID       int64  `gorm:"primaryKey"`
	StreamerName string `gorm:"type:string;size:256"`
	AccountName  string `gorm:"type:string;size:256"`
	PrivateKey   string `gorm:"type:string;size:256"`
	PublicKey    string `gorm:"type:string;size:256"`
}

func (s Streamer) TableName() string { return "t_streamer" }

func init() {
	registerModel(&Streamer{})
}

type StreamerDAL struct{}

func GetStreamerDAL() StreamerDAL { return StreamerDAL{} }

func (dal StreamerDAL) All() ([]Streamer, error) {
	var ret []Streamer
	err := gDB.Find(&ret).Error
	return ret, err
}
