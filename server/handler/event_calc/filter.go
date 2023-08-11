package event_calc

import "github.com/zerozwt/octant/server/bs"

type Filter interface {
	OK(user *UserData, strip *UserStrip) bool
}

type eventFilterAnd struct {
	filters []Filter
}

func (f *eventFilterAnd) OK(user *UserData, strip *UserStrip) bool {
	ss := NewEventUserStrip()
	for _, filter := range f.filters {
		if !filter.OK(user, ss) {
			return false
		}
	}
	strip.combine(ss)
	return true
}

type eventFilterOr struct {
	filters []Filter
}

func (f *eventFilterOr) OK(user *UserData, strip *UserStrip) bool {
	ss := NewEventUserStrip()
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

func (f *eventFilterGift) OK(user *UserData, strip *UserStrip) bool {
	value := int64(0)
	ss := NewEventUserStrip()
	for idx := range user.Gift {
		if user.Gift[idx].GiftID != f.id || user.Gift[idx].SendTime < f.startTs || user.Gift[idx].SendTime > f.endTs {
			continue
		}
		if f.total {
			ss.gift[idx] = true
			value += user.Gift[idx].GiftCount
		} else if user.Gift[idx].GiftCount >= f.count {
			ss.gift[idx] = true
			value = user.Gift[idx].GiftCount
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

func (f *eventFilterSC) OK(user *UserData, strip *UserStrip) bool {
	value := int64(0)
	ss := NewEventUserStrip()
	for idx := range user.SC {
		if user.SC[idx].SendTime < f.startTs || user.SC[idx].SendTime > f.endTs {
			continue
		}
		if f.total {
			value += user.SC[idx].Price
			ss.gift[idx] = true
		} else if user.SC[idx].Price >= f.count {
			ss.gift[idx] = true
			value = user.SC[idx].Price
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

func (f *eventFilterMember) OK(user *UserData, strip *UserStrip) bool {
	value := int64(0)
	ss := NewEventUserStrip()
	for idx := range user.Member {
		if user.Member[idx].SendTime < f.startTs || user.Member[idx].SendTime > f.endTs {
			continue
		}
		if (1<<user.Member[idx].GuardLevel)&f.levelMask == 0 {
			continue
		}
		if f.total {
			ss.gift[idx] = true
			value += int64(user.Member[idx].Count)
		} else if int64(user.Member[idx].Count) >= f.count {
			ss.gift[idx] = true
			value = int64(user.Member[idx].Count)
		}
	}
	if value >= f.count {
		strip.combine(ss)
		return true
	}
	return false
}

func BuildFilter(cond *bs.EventCondition) Filter {
	if cond.IsAnd() {
		ret := &eventFilterAnd{}
		for idx := range cond.Subs {
			ret.filters = append(ret.filters, BuildFilter(&cond.Subs[idx]))
		}
		return ret
	}
	if cond.IsOr() {
		ret := &eventFilterOr{}
		for idx := range cond.Subs {
			ret.filters = append(ret.filters, BuildFilter(&cond.Subs[idx]))
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
