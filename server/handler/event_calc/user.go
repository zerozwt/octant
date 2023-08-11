package event_calc

import (
	"sync/atomic"

	jsoniter "github.com/json-iterator/go"
	"github.com/zerozwt/octant/server/db"
)

type UserData struct {
	UID    int64
	SendTs int64
	Name   string

	Gift   []*db.GiftRecord
	SC     []*db.SuperChatRecord
	Member []*db.MembershipRecord

	stripped atomic.Bool
}

func NewEventUser(uid int64) *UserData {
	return &UserData{
		UID: uid,
	}
}

func EventUserfromDB(item *db.RewardUser) (*UserData, error) {
	ret := &UserData{
		UID:    item.UID,
		SendTs: item.Time,
		Name:   item.UserName,
	}
	ret.stripped.Store(true)

	tmp := struct {
		Gift   []*db.GiftRecord       `json:"gift"`
		SC     []*db.SuperChatRecord  `json:"sc"`
		Member []*db.MembershipRecord `json:"member"`
	}{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.UnmarshalFromString(item.Columns, &tmp)

	ret.Gift = tmp.Gift
	ret.SC = tmp.SC
	ret.Member = tmp.Member

	return ret, err
}

func (user *UserData) Strip(strip *UserStrip) *UserData {
	if user.stripped.CompareAndSwap(false, true) {
		user.Gift = stripArray(user.Gift, strip.gift)
		user.SC = stripArray(user.SC, strip.sc)
		user.Member = stripArray(user.Member, strip.member)

		for _, item := range user.Gift {
			if user.SendTs == 0 || user.SendTs > item.SendTime {
				user.SendTs = item.SendTime
			}
			user.Name = item.SenderName
		}
		for _, item := range user.SC {
			if user.SendTs == 0 || user.SendTs > item.SendTime {
				user.SendTs = item.SendTime
			}
			user.Name = item.SenderName
		}
		for _, item := range user.Member {
			if user.SendTs == 0 || user.SendTs > item.SendTime {
				user.SendTs = item.SendTime
			}
			user.Name = item.SenderName
		}
	}
	return user
}

func (user *UserData) Column() string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := map[string]any{}
	if len(user.Gift) > 0 {
		data["gift"] = user.Gift
	}
	if len(user.SC) > 0 {
		data["sc"] = user.SC
	}
	if len(user.Member) > 0 {
		data["member"] = user.Member
	}
	ret, _ := json.MarshalToString(data)
	return ret
}

func stripArray[T any](array []T, strip map[int]bool) []T {
	ret := []T{}
	for idx := range strip {
		ret = append(ret, array[idx])
	}
	return ret
}

type UserStrip struct {
	gift   map[int]bool
	sc     map[int]bool
	member map[int]bool
}

func (s *UserStrip) combine(strip *UserStrip) {
	for k, v := range strip.gift {
		s.gift[k] = v
	}
	for k, v := range strip.sc {
		s.sc[k] = v
	}
	for k, v := range strip.member {
		s.member[k] = v
	}
}

func NewEventUserStrip() *UserStrip {
	return &UserStrip{
		gift:   map[int]bool{},
		sc:     map[int]bool{},
		member: map[int]bool{},
	}
}
