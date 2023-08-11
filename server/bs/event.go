package bs

import (
	"fmt"

	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

type TimeRange [2]int64

func (t TimeRange) Start() int64 { return t[0] }
func (t TimeRange) End() int64   { return t[1] }
func (t TimeRange) Combine(value TimeRange) TimeRange {
	ret := TimeRange{t[0], t[1]}
	if ret[0] == 0 || ret[0] > value[0] {
		ret[0] = value[0]
	}
	if ret[1] < value[1] {
		ret[1] = value[1]
	}
	return ret
}

type ConditionTimeRange struct {
	Range map[string]TimeRange
}

type EventCondition struct {
	Subs []EventCondition `json:"sub_conditions"`

	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Type      string `json:"type"`

	Mode  string `json:"mode"`
	Count int64  `json:"count"`

	GiftID     int64 `json:"gift_id"`
	GuardLevel []int `json:"guard_levels"`

	startTime int64
	endTime   int64
}

func (c *EventCondition) Validate(ctx *swe.Context) error {
	if !c.IsValidType() {
		return fmt.Errorf("invalid type %s", c.Type)
	}

	if c.IsMulti() {
		if len(c.Subs) == 0 {
			return fmt.Errorf("condition group have no sub conditions")
		}
		for idx := range c.Subs {
			if err := c.Subs[idx].Validate(ctx); err != nil {
				return err
			}
		}
		return nil
	}

	if !utils.IsValidTimeString(c.StartTime) {
		return fmt.Errorf("start time invalid: %s", c.StartTime)
	}
	if !utils.IsValidTimeString(c.EndTime) {
		return fmt.Errorf("end time invalid: %s", c.EndTime)
	}

	var err error
	c.startTime, err = utils.TimeStringToUTC(c.StartTime)
	if err != nil {
		return err
	}
	c.endTime, err = utils.TimeStringToUTC(c.EndTime)
	if err != nil {
		return err
	}

	if !c.IsOnce() && !c.IsTotal() {
		return fmt.Errorf("invalid mode: %s", c.Mode)
	}
	if c.IsGift() && (c.Count < 1 || c.GiftID < 1) {
		return fmt.Errorf("invalid gift param: count %d id %d", c.Count, c.GiftID)
	}
	if c.IsSuperChat() && c.Count < 1 {
		return fmt.Errorf("invalid super chat price")
	}
	if c.IsMember() {
		if len(c.GuardLevel) == 0 {
			return fmt.Errorf("empty guard level")
		}
		for _, level := range c.GuardLevel {
			if level < 1 || level > 3 {
				return fmt.Errorf("invalid guard level %d", level)
			}
		}
	}
	return nil
}

func (c *EventCondition) ScheduleTime() int64 {
	if c.IsMulti() {
		ret := int64(0)
		for idx := range c.Subs {
			ts := c.Subs[idx].ScheduleTime()
			if ts > ret {
				ret = ts
			}
		}
		return ret
	}
	return c.endTime
}

func (c *EventCondition) IsValidType() bool {
	return c.IsMulti() || c.IsGift() || c.IsSuperChat() || c.IsMember()
}

func (c *EventCondition) IsAnd() bool       { return c.Type == "and" }
func (c *EventCondition) IsOr() bool        { return c.Type == "or" }
func (c *EventCondition) IsMulti() bool     { return c.IsAnd() || c.IsOr() }
func (c *EventCondition) IsGift() bool      { return c.Type == "gift" }
func (c *EventCondition) IsSuperChat() bool { return c.Type == "sc" }
func (c *EventCondition) IsMember() bool    { return c.Type == "member" }
func (c *EventCondition) IsOnce() bool      { return c.Mode == "once" }
func (c *EventCondition) IsTotal() bool     { return c.Mode == "total" }

func (c *EventCondition) StartTs() int64 { return c.startTime }
func (c *EventCondition) EndTs() int64   { return c.endTime }

func (c *EventCondition) CalculateRange(value *ConditionTimeRange) {
	if c.IsMulti() {
		for idx := range c.Subs {
			c.Subs[idx].CalculateRange(value)
		}
		return
	}
	value.Range[c.Type] = value.Range[c.Type].Combine(TimeRange{c.startTime, c.endTime})
}

type EventAddReq struct {
	Name      string         `json:"name"`
	Reward    string         `json:"reward"`
	Condition EventCondition `json:"conditions"`
}

func (req *EventAddReq) Validate(ctx *swe.Context) error {
	if len(req.Name) == 0 {
		return fmt.Errorf("empty event name")
	}
	return req.Condition.Validate(ctx)
}

type EventModifyReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Reward string `json:"reward"`
}

func (req *EventModifyReq) Validate(ctx *swe.Context) error {
	if req.ID < 1 {
		return fmt.Errorf("invalid id %d", req.ID)
	}
	if len(req.Name) == 0 {
		return fmt.Errorf("no event name")
	}

	return nil
}

type EventDetailRsp struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Reward    string         `json:"reward"`
	Condition EventCondition `json:"conditions"`
	Status    int            `json:"status"`
}

type EventUserListReq struct {
	EventID int64 `json:"event_id"`
	Page    int   `json:"page"`
	Size    int   `json:"size"`
}

func (req EventUserListReq) Validate(ctx *swe.Context) error {
	if req.Page < 1 {
		return fmt.Errorf("invalid page %d", req.Page)
	}
	if req.Size < 1 || req.Size > 100 {
		return fmt.Errorf("invalid size %d", req.Size)
	}
	if req.EventID < 0 {
		return fmt.Errorf("invalid event id %d", req.EventID)
	}
	return nil
}

type EventUserListItem struct {
	UID   int64          `json:"uid"`
	Name  string         `json:"name"`
	Time  string         `json:"time"`
	Cols  map[string]any `json:"cols"`
	Block bool           `json:"block"`
}

type EventUIDReq struct {
	EventID int64 `json:"event_id"`
	UID     int64 `json:"uid"`
}
