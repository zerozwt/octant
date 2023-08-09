package handler

import (
	"sync/atomic"

	jsoniter "github.com/json-iterator/go"
	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
)

type eventUserData struct {
	uid    int64
	sendTs int64
	name   string

	gift   []*db.GiftRecord
	sc     []*db.SuperChatRecord
	member []*db.MembershipRecord

	stripped atomic.Bool
}

func newEventUser(uid int64) *eventUserData {
	return &eventUserData{
		uid: uid,
	}
}

func (user *eventUserData) strip(strip *eventUserStrip) *eventUserData {
	if user.stripped.CompareAndSwap(false, true) {
		user.gift = stripArray(user.gift, strip.gift)
		user.sc = stripArray(user.sc, strip.sc)
		user.member = stripArray(user.member, strip.member)

		for _, item := range user.gift {
			if user.sendTs == 0 || user.sendTs > item.SendTime {
				user.sendTs = item.SendTime
			}
			user.name = item.SenderName
		}
		for _, item := range user.sc {
			if user.sendTs == 0 || user.sendTs > item.SendTime {
				user.sendTs = item.SendTime
			}
			user.name = item.SenderName
		}
		for _, item := range user.member {
			if user.sendTs == 0 || user.sendTs > item.SendTime {
				user.sendTs = item.SendTime
			}
			user.name = item.SenderName
		}
	}
	return user
}

func (user *eventUserData) column() string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := map[string]any{}
	if len(user.gift) > 0 {
		data["gift"] = user.gift
	}
	if len(user.sc) > 0 {
		data["sc"] = user.sc
	}
	if len(user.member) > 0 {
		data["member"] = user.member
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

type eventUserStrip struct {
	gift   map[int]bool
	sc     map[int]bool
	member map[int]bool
}

func (s *eventUserStrip) combine(strip *eventUserStrip) {
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

func newEventUserStrip() *eventUserStrip {
	return &eventUserStrip{
		gift:   map[int]bool{},
		sc:     map[int]bool{},
		member: map[int]bool{},
	}
}

type eventFilter interface {
	OK(user *eventUserData, strip *eventUserStrip) bool
}

type eventFilterAnd struct {
	filters []eventFilter
}

func (f *eventFilterAnd) OK(user *eventUserData, strip *eventUserStrip) bool {
	ss := newEventUserStrip()
	for _, filter := range f.filters {
		if !filter.OK(user, ss) {
			return false
		}
	}
	strip.combine(ss)
	return true
}

type eventFilterOr struct {
	filters []eventFilter
}

func (f *eventFilterOr) OK(user *eventUserData, strip *eventUserStrip) bool {
	ss := newEventUserStrip()
	ret := false
	for _, filter := range f.filters {
		if filter.OK(user, ss) {
			ret = true
		}
	}
	if ret {
		strip.combine(ss)
	}
	return ret
}

type eventFilterGift struct {
	id      int64
	count   int64
	startTs int64
	endTs   int64
	total   bool
}

func (f *eventFilterGift) OK(user *eventUserData, strip *eventUserStrip) bool {
	value := int64(0)
	ss := newEventUserStrip()
	for idx := range user.gift {
		if user.gift[idx].GiftID != f.id || user.gift[idx].SendTime < f.startTs || user.gift[idx].SendTime > f.endTs {
			continue
		}
		if f.total {
			ss.gift[idx] = true
			value += user.gift[idx].GiftCount
		} else if user.gift[idx].GiftCount >= f.count {
			ss.gift[idx] = true
			value = user.gift[idx].GiftCount
		}
	}
	if value >= f.count {
		strip.combine(ss)
		return true
	}
	return false
}

type eventFilterSC struct {
	count   int64
	startTs int64
	endTs   int64
	total   bool
}

func (f *eventFilterSC) OK(user *eventUserData, strip *eventUserStrip) bool {
	value := int64(0)
	ss := newEventUserStrip()
	for idx := range user.sc {
		if user.sc[idx].SendTime < f.startTs || user.sc[idx].SendTime > f.endTs {
			continue
		}
		if f.total {
			value += user.sc[idx].Price
			ss.gift[idx] = true
		} else if user.sc[idx].Price >= f.count {
			ss.gift[idx] = true
			value = user.sc[idx].Price
		}
	}
	if value >= f.count {
		strip.combine(ss)
		return true
	}
	return false
}

type eventFilterMember struct {
	startTs   int64
	endTs     int64
	count     int64
	levelMask int
	total     bool
}

func (f *eventFilterMember) OK(user *eventUserData, strip *eventUserStrip) bool {
	value := int64(0)
	ss := newEventUserStrip()
	for idx := range user.member {
		if user.member[idx].SendTime < f.startTs || user.member[idx].SendTime > f.endTs {
			continue
		}
		if (1<<user.member[idx].GuardLevel)&f.levelMask == 0 {
			continue
		}
		if f.total {
			ss.gift[idx] = true
			value += int64(user.member[idx].Count)
		} else if int64(user.member[idx].Count) >= f.count {
			ss.gift[idx] = true
			value = int64(user.member[idx].Count)
		}
	}
	if value >= f.count {
		strip.combine(ss)
		return true
	}
	return false
}

func buildEventFilter(cond *bs.EventCondition) eventFilter {
	if cond.IsAnd() {
		ret := &eventFilterAnd{}
		for idx := range cond.Subs {
			ret.filters = append(ret.filters, buildEventFilter(&cond.Subs[idx]))
		}
		return ret
	}
	if cond.IsOr() {
		ret := &eventFilterOr{}
		for idx := range cond.Subs {
			ret.filters = append(ret.filters, buildEventFilter(&cond.Subs[idx]))
		}
		return ret
	}
	if cond.IsGift() {
		return &eventFilterGift{
			id:      cond.GiftID,
			count:   cond.Count,
			startTs: cond.StartTs(),
			endTs:   cond.EndTs(),
			total:   cond.IsTotal(),
		}
	}
	if cond.IsSuperChat() {
		return &eventFilterSC{
			count:   cond.Count,
			startTs: cond.StartTs(),
			endTs:   cond.EndTs(),
			total:   cond.IsTotal(),
		}
	}
	// member
	mask := 0
	for _, level := range cond.GuardLevel {
		mask = mask | (1 << level)
	}
	return &eventFilterMember{
		count:     cond.Count,
		startTs:   cond.StartTs(),
		endTs:     cond.EndTs(),
		total:     cond.IsTotal(),
		levelMask: mask,
	}
}
