package bridge

type Publisher interface {
	AddRoom(roomID int64) error
	DelRoom(roomID int64) error
}

type Receiver interface {
	OnAddRoom(roomID int64)
	OnDelRoom(roomID int64)
}

type Bridge interface {
	Publisher
	SetReceiver(Receiver)
	Start() error
	Stop() error
}

var gBridge Bridge = nil

func SetBridge(b Bridge) { gBridge = b }
func GetBridge() Bridge  { return gBridge }
