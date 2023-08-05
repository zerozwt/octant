package collector

import (
	"sync"
	"sync/atomic"
	"time"

	dm "github.com/zerozwt/BLiveDanmaku"
	"github.com/zerozwt/BLiveDanmaku/cmds"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/swe"
)

type room struct {
	id     int64
	client *dm.Client
	lock   sync.Mutex

	stopped int32
}

type Collector struct {
	rooms map[int64]*room
	lock  sync.Mutex
	stop  atomic.Bool
}

var cc *Collector = &Collector{rooms: map[int64]*room{}}

func GetCollector() *Collector {
	return cc
}

func (c *Collector) OnAddRoom(roomID int64) {
	if c.stop.Load() {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	if c.stop.Load() {
		return
	}

	if _, ok := c.rooms[roomID]; !ok {
		r := newRoom(roomID)
		r.Start()
		c.rooms[roomID] = r
	}
}

func (c *Collector) OnDelRoom(roomID int64) {
	if c.stop.Load() {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	if c.stop.Load() {
		return
	}

	if r, ok := c.rooms[roomID]; ok {
		r.Stop()
		delete(c.rooms, roomID)
	}
}

func (c *Collector) Stop() {
	c.stop.Store(true)
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, r := range c.rooms {
		r.Stop()
	}
	c.rooms = map[int64]*room{}
}

func newRoom(roomID int64) *room {
	return &room{
		id: roomID,
	}
}

func (r *room) Start() {
	logger := swe.CtxLogger(nil)
	logger.Info("connecting to live room %d", r.id)
	if err := r.connect(); err != nil {
		r.onDisconnect(nil, err)
	}
}

func (r *room) Stop() {
	if atomic.CompareAndSwapInt32(&r.stopped, 0, 1) {
		swe.CtxLogger(nil).Info("stop monitoring live room %d", r.id)
		r.lock.Lock()
		defer r.lock.Unlock()
		if r.client != nil {
			r.client.Close()
		}
	}
}

func (r *room) connect() error {
	if atomic.LoadInt32(&r.stopped) != 0 {
		return nil
	}

	conf := &dm.ClientConf{
		OnNetError:         r.onDisconnect,
		OnServerDisconnect: r.onDisconnect,
	}
	conf.AddCmdHandler(dm.CMD_SUPER_CHAT_MESSAGE, r.onSuperChat)
	conf.AddCmdHandler(dm.CMD_GUARD_BUY, r.onGuardBuy)
	conf.AddCmdHandler(dm.CMD_SEND_GIFT, r.onGift)

	tmp, err := dm.Dial(r.id, conf)
	if err != nil {
		return err
	}

	swe.CtxLogger(nil).Info("successfully connectoed to live room %d", r.id)

	r.lock.Lock()
	defer r.lock.Unlock()
	if atomic.LoadInt32(&r.stopped) == 0 {
		r.client = tmp
	} else {
		r.client = nil
		tmp.Close()
	}

	return nil
}

func (r *room) onDisconnect(client *dm.Client, err error) {
	if atomic.LoadInt32(&r.stopped) != 0 {
		return
	}

	logger := swe.CtxLogger(nil)
	logger.Error("coneection to live room %d interrupted: %v", r.id, err)

	go func() {
		waitTime := time.Second
		maxWait := time.Second * 30

		for {
			if atomic.LoadInt32(&r.stopped) != 0 {
				return
			}

			logger.Info("try reconnect to live room %d", r.id)
			err := r.connect()
			if err != nil {
				logger.Error("reconnect to live room %d failed: %v, retry after %v", r.id, err, waitTime)
				time.Sleep(waitTime)
				waitTime *= 2
				if waitTime > maxWait {
					waitTime = maxWait
				}
			} else {
				return
			}
		}
	}()
}

func (r *room) onSuperChat(client *dm.Client, cmd string, data []byte) bool {
	if atomic.LoadInt32(&r.stopped) != 0 {
		return true
	}
	msg := cmds.SuperChatMessage{}
	if err := msg.Decode(data); err != nil {
		swe.CtxLogger(nil).Error("decode superchat message from live room %d failed: %v", r.id, err)
		return true
	}

	err := db.GetSCDal().Insert(nil, r.id, msg.Timestamp, msg.UID, msg.User.UserName, int64(msg.Price),
		msg.Message, msg.BackgroundColor, msg.MessageFontColor)
	if err != nil {
		swe.CtxLogger(nil).Error("insert sc to room %d failed: %v, original: %s", r.id, err, string(data))
	}
	return false
}

func (r *room) onGuardBuy(client *dm.Client, cmd string, data []byte) bool {
	if atomic.LoadInt32(&r.stopped) != 0 {
		return true
	}
	msg := cmds.GuardBuy{}
	if err := msg.Decode(data); err != nil {
		swe.CtxLogger(nil).Error("decode member buy message from live room %d failed: %v", r.id, err)
		return true
	}

	err := db.GetMemberDal().Insert(nil, r.id, msg.StartTime, msg.UID, msg.UserName, msg.GuardLevel, msg.Num)
	if err != nil {
		swe.CtxLogger(nil).Error("insert member to room %d failed: %v, original: %s", r.id, err, string(data))
	}
	return false
}

func (r *room) onGift(client *dm.Client, cmd string, data []byte) bool {
	if atomic.LoadInt32(&r.stopped) != 0 {
		return true
	}
	msg := cmds.SendGift{}
	if err := msg.Decode(data); err != nil {
		swe.CtxLogger(nil).Error("decode gift message from live room %d failed: %v", r.id, err)
		return true
	}

	if msg.CoinType != "gold" {
		return false
	}

	gift := db.GiftRecord{
		BatchID:    msg.BatchComboID,
		RoomID:     r.id,
		SendTime:   msg.Timestamp,
		SenderUID:  msg.UID,
		SenderName: msg.UserName,
		GiftID:     msg.GiftID,
		GiftName:   msg.GiftName,
		GiftPrice:  msg.Price,
		GiftCount:  int64(msg.Num),
	}

	if err := db.GetGiftDAL().Insert(nil, &gift); err != nil {
		swe.CtxLogger(nil).Error("insert gift to room %d failed: %v, original: %s", r.id, err, string(data))
	}

	db.GetGiftDAL().UpdateGiftInfo(nil, msg.GiftID, msg.GiftName, msg.Price)

	return false
}
