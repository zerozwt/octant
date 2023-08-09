package handler

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/zerozwt/octant/server/async_task"
	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

func init() {
	registerHandler(GET, "/event/list", event.list, session.CheckStreamer)
	registerHandler(POST, "/event/add", event.add, session.CheckStreamer)

	async_task.RegisterHandler(asyncTaskCalculateEventList, event.calculate)
}

type eventHandler struct{}

const (
	asyncTaskCalculateEventList = "CalculateEventList"
)

var event eventHandler

func (ins eventHandler) list(ctx *swe.Context, req *bs.PageReq) (*bs.PageRsp, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	count, list, err := db.GetRewardEventDAL().Page(ctx, st.RoomID, (req.Page-1)*req.Size, req.Size)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	ret := &bs.PageRsp{Count: count, List: []any{}}
	for _, item := range list {
		ret.List = append(ret.List, map[string]any{
			"id":     item.ID,
			"name":   item.EventName,
			"status": item.Status,
		})
	}

	return ret, nil
}

func (ins eventHandler) add(ctx *swe.Context, req *bs.EventAddReq) (*bs.Nothing, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	event := &db.RewardEvent{
		ID:            utils.GenerateID(),
		RoomID:        st.RoomID,
		EventName:     req.Name,
		RewardContent: req.Reward,
		CreateTime:    time.Now().Unix(),
		Status:        db.EVENT_IDLE,
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	event.Conditions, _ = json.MarshalToString(req.Condition)

	err := db.GetRewardEventDAL().Put(ctx, event)
	if err != nil {
		swe.CtxLogger(ctx).Error("write event to db error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	err = async_task.GetScheduler().AddTask(ctx, asyncTaskCalculateEventList, fmt.Sprint(event.ID), req.Condition.ScheduleTime()+60, nil)
	if err != nil {
		swe.CtxLogger(ctx).Error("create async task error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	return &bs.Nothing{}, nil
}

func (ins eventHandler) calculate(ctx *swe.Context, taskCtx async_task.TaskContext) error {
	logger := swe.CtxLogger(ctx)

	// get event id
	evtID, err := strconv.ParseInt(taskCtx.Param(), 10, 64)
	if err != nil {
		logger.Error("parse event id failed: %v param: %s", err, taskCtx.Param())
		return err
	}

	logger.Info("start calculating event list for event %d", evtID)

	// load event
	event, err := db.GetRewardEventDAL().Get(ctx, evtID)
	if err != nil {
		logger.Error("find event from db failed: %v", err)
		return err
	}
	err = db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_CALCULATING)
	if err != nil {
		logger.Error("set event status to calculating failed: %v", err)
		return err
	}

	// decode condition
	cond := bs.EventCondition{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.UnmarshalFromString(event.Conditions, &cond)
	if err != nil {
		db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
		logger.Error("decode condition failed: %v", err)
		return err
	}
	err = cond.Validate(ctx)
	if err != nil {
		db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
		logger.Error("validate condition failed: %v", err)
		return err
	}

	logger.Info("loading needed data for event %d", evtID)

	// load needed data
	timeRange := bs.ConditionTimeRange{}
	cond.CalculateRange(&timeRange)
	users := map[int64]*eventUserData{}

	if tr, ok := timeRange.Range["gift"]; ok {
		rec, err := db.GetGiftDAL().Range(ctx, event.RoomID, tr.Start(), tr.End())
		if err != nil {
			db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
			logger.Error("load gift record failed: %v", err)
			return err
		}
		logger.Info("%d records of gift loaded", len(rec))
		for _, item := range rec {
			if _, ok := users[item.SenderUID]; !ok {
				users[item.SenderUID] = newEventUser(item.SenderUID)
			}
			users[item.SenderUID].gift = append(users[item.SenderUID].gift, item)
		}
	}
	if tr, ok := timeRange.Range["sc"]; ok {
		rec, err := db.GetSCDal().Range(ctx, event.RoomID, tr.Start(), tr.End())
		if err != nil {
			db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
			logger.Error("load sc record failed: %v", err)
			return err
		}
		logger.Info("%d records of super chat loaded", len(rec))
		for _, item := range rec {
			if _, ok := users[item.SenderUID]; !ok {
				users[item.SenderUID] = newEventUser(item.SenderUID)
			}
			users[item.SenderUID].sc = append(users[item.SenderUID].sc, item)
		}
	}
	if tr, ok := timeRange.Range["member"]; ok {
		rec, err := db.GetMemberDal().Range(ctx, event.RoomID, tr.Start(), tr.End())
		if err != nil {
			db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
			logger.Error("load membership record failed: %v", err)
			return err
		}
		logger.Info("%d records of membership loaded", len(rec))
		for _, item := range rec {
			if _, ok := users[item.SenderUID]; !ok {
				users[item.SenderUID] = newEventUser(item.SenderUID)
			}
			users[item.SenderUID].member = append(users[item.SenderUID].member, item)
		}
	}

	logger.Info("filtering data for event %d", evtID)

	// filter sender
	filter := buildEventFilter(&cond)
	tmp := users
	users = map[int64]*eventUserData{}
	for uid, data := range tmp {
		strip := newEventUserStrip()
		if filter.OK(data, strip) {
			users[uid] = data.strip(strip)
		}
	}

	logger.Info("%d users after filter, event %d", len(users), evtID)

	// delete older list record
	if err = db.GetRewardEventDAL().ClearUsers(ctx, evtID); err != nil {
		db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
		logger.Error("clear old list for event %d failed: %v", evtID, err)
		return err
	}

	// insert new list record
	uids := make([]int64, 0, len(users))
	for uid := range users {
		uids = append(uids, uid)
	}
	sort.Slice(uids, func(i, j int) bool {
		return users[uids[i]].sendTs < users[uids[j]].sendTs
	})

	data := make([]db.RewardUser, 0, len(users))
	for _, uid := range uids {
		data = append(data, db.RewardUser{
			EventID:  evtID,
			UID:      uid,
			UserName: users[uid].name,
			Time:     users[uid].sendTs,
			Columns:  users[uid].column(),
		})
	}

	if err = db.GetRewardEventDAL().PutUsers(ctx, data); err != nil {
		db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
		logger.Error("write new list for event %d failed: %v", evtID, err)
		return err
	}

	// set status to ready
	if err = db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_READY); err != nil {
		logger.Error("set event status to ready failed: %v ", err)
		return err
	}

	logger.Info("list calculation for event %d done", evtID)
	return nil
}
