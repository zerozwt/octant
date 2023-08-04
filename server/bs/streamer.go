package bs

import (
	"fmt"

	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

type StreamerCheckLoginRsp struct {
	ID      int64  `json:"room_id"`
	Name    string `json:"name"`
	Account string `json:"account_name"`
}

type StreamerLoginReq struct {
	Account  string `json:"name"`
	Password string `json:"password"`
}

type SimpleSearchReq struct {
	DataSource string `json:"datasource"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Filter     struct {
		UID              int64  `json:"sender_uid"`
		Name             string `json:"sender_name"`
		GiftID           int64  `json:"gift_id"`
		SuperchatContent string `json:"sc_content"`
		GuardLevel       []int  `json:"guard_level"`
	} `json:"filter"`

	startTs int64
	endTs   int64
}

func (req SimpleSearchReq) Validate(ctx *swe.Context) error {
	if req.Page < 1 {
		return fmt.Errorf("invalid page %d", req.Page)
	}
	if req.Size < 1 || req.Size > 100 {
		return fmt.Errorf("invalid size %d", req.Size)
	}
	if !utils.IsValidTimeString(req.StartTime) {
		return fmt.Errorf("start time format invalid: %s", req.StartTime)
	}
	if !utils.IsValidTimeString(req.EndTime) {
		return fmt.Errorf("end time format invalid: %s", req.EndTime)
	}
	if req.DataSource != "sc" && req.DataSource != "gift" && req.DataSource != "member" {
		return fmt.Errorf("invalid datasource: %s", req.DataSource)
	}
	return nil
}

func (req *SimpleSearchReq) ParseTimeRange() (err error) {
	req.startTs, err = utils.TimeStringToUTC(req.StartTime)
	if err != nil {
		return err
	}
	req.endTs, err = utils.TimeStringToUTC(req.EndTime)
	return
}

func (req SimpleSearchReq) StartTs() int64    { return req.startTs }
func (req SimpleSearchReq) EndTs() int64      { return req.endTs }
func (req SimpleSearchReq) IsSuperChat() bool { return req.DataSource == "sc" }
func (req SimpleSearchReq) IsGift() bool      { return req.DataSource == "gift" }
func (req SimpleSearchReq) IsMember() bool    { return req.DataSource == "member" }

type SimpleSearchItem struct {
	UID  int64  `json:"uid"`
	Name string `json:"name"`
	Time string `json:"time"`
	Gift struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Price int64  `json:"price"`
		Count int64  `json:"count"`
	} `json:"gift"`
	SuperChat struct {
		Price     int64  `json:"price"`
		Content   string `json:"content"`
		BgColor   string `json:"bg_color"`
		FontColor string `json:"font_color"`
	} `json:"sc"`
	Member struct {
		Level int `json:"level"`
		Count int `json:"count"`
	} `json:"guard"`
}
