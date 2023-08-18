package bs

import (
	"fmt"

	"github.com/zerozwt/swe"
)

type DMSenderInfo struct {
	UID      int64  `json:"uid"`
	SessData string `json:"sess"`
	JCT      string `json:"jct"`
}

func (req *DMSenderInfo) Validate(ctx *swe.Context) error {
	if req.UID <= 0 {
		return fmt.Errorf("invalid sender uid: %d", req.UID)
	}
	if len(req.SessData) == 0 || len(req.JCT) == 0 {
		return fmt.Errorf("no sender cookie info")
	}
	return nil
}

type DMCreateReq struct {
	EventID     int64        `json:"event_id"`
	Name        string       `json:"name"`
	Content     string       `json:"content"`
	BatchMax    int          `json:"batch_max"`
	TLS         bool         `json:"tls"`
	IntervalMin int          `json:"interval_min"`
	IntervalMax int          `json:"interval_max"`
	RunTask     bool         `json:"run_task"`
	Sender      DMSenderInfo `json:"sender"`
}

func (req *DMCreateReq) Validate(ctx *swe.Context) error {
	if len(req.Name) == 0 {
		return fmt.Errorf("dm task no name")
	}
	if len(req.Content) == 0 {
		return fmt.Errorf("no dm content")
	}
	if len(req.Content) > 4096 {
		return fmt.Errorf("content too long")
	}
	if req.IntervalMin > req.IntervalMax || req.IntervalMin < 0 || req.IntervalMax < 0 {
		return fmt.Errorf("invalid interval range")
	}
	if req.RunTask {
		return req.Sender.Validate(ctx)
	}
	return nil
}

type DMTaskListItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	BatchMax    int    `json:"batch_max"`
	IntervalMin int    `json:"interval_min"`
	IntervalMax int    `json:"interval_max"`
	Status      int    `json:"status"`
	Event       struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"event"`
}

type DMTaskDetail struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	BatchMax    int    `json:"batch_max"`
	IntervalMin int    `json:"interval_min"`
	IntervalMax int    `json:"interval_max"`
	Status      int    `json:"status"`
	Succ        int    `json:"succ"`
	Fail        int    `json:"fail"`
	Total       int    `json:"total"`
}

type DMSetSenderReq struct {
	TaskID  int64        `json:"task_id"`
	RunTask bool         `json:"run_task"`
	Sender  DMSenderInfo `json:"sender"`
}

func (req *DMSetSenderReq) Validate(ctx *swe.Context) error {
	if req.RunTask {
		return req.Sender.Validate(ctx)
	}
	return nil
}

type DMSwitchReq struct {
	TaskID  int64 `json:"task_id"`
	RunTask bool  `json:"run_task"`
}
