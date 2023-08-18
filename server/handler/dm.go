package handler

import (
	"fmt"
	"time"

	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/handler/batch_dm"
	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

func init() {
	registerHandler(POST, "/dm/task", dmsg.create, session.CheckStreamer)
	registerHandler(GET, "/dm/tasks", dmsg.page, session.CheckStreamer)
	registerHandler(GET, "/dm/task/detail", dmsg.detail, session.CheckStreamer)
	registerHandler(POST, "/dm/sender", dmsg.setSender, session.CheckStreamer)
	registerHandler(POST, "/dm/siwtch", dmsg.setState, session.CheckStreamer)
}

type dmHandler struct{}

var dmsg dmHandler

func (ins dmHandler) create(ctx *swe.Context, req *bs.DMCreateReq) (*bs.Nothing, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	st, _ := session.GetStreamerSession(ctx)

	// load event
	event, err := db.GetRewardEventDAL().GetByRoomID(ctx, req.EventID, st.RoomID)
	if err != nil {
		logger.Error("query db for event %d error %v", req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if event == nil {
		logger.Error("query db for event %d room id %d not found", req.EventID, st.RoomID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("event not found"))
	}
	if event.Status == db.EVENT_READY {
		logger.Error("query db for event %d room id %d not ready", req.EventID, st.RoomID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("event not ready"))
	}

	// load event records
	records, err := db.GetRewardEventDAL().Users(ctx, req.EventID)
	if err != nil {
		logger.Error("query users for event %d error %v", req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	uids := make([]int64, 0, len(records))
	for _, item := range records {
		uids = append(uids, item.UID)
	}

	// load dd info
	userMap, err := db.GetDDInfoDAL().BatchGet(ctx, uids)
	if err != nil {
		logger.Error("query dd info for event %d error %v", req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	// build content converter
	proto := "http://"
	if req.TLS {
		proto = "https://"
	}
	builder := batch_dm.MakeBuilder(req.Content)
	buildCtx := batch_dm.BuildCtx{
		StreamerName: st.StreamerName,
		Event:        event,
		InviteLink:   proto + ctx.Request.Host + "/dd/access?access_code=",
		InfoMap:      userMap,
	}

	// save task to db
	task := db.DMTask{
		ID:          utils.GenerateID(),
		RoomID:      st.RoomID,
		EventID:     event.ID,
		TaskName:    req.Name,
		MsgType:     1,
		Content:     req.Content,
		BatchMax:    req.BatchMax,
		Status:      db.DM_TASK_STATUS_PAUSED,
		IntervalMin: req.IntervalMin,
		IntervalMax: req.IntervalMax,
		CreateTime:  time.Now().Unix(),
	}

	if err := db.GetDirectMsgDAL().Put(ctx, &task); err != nil {
		logger.Error("save task info error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	// generate DM details
	details := make([]db.DMDetail, 0, len(userMap))
	for _, item := range records {
		content, err := builder.BuildContent(ctx, &item, &buildCtx)
		if err != nil {
			logger.Error("generate msg content for uid %d failed: %v", item.UID, err)
		} else {
			details = append(details, db.DMDetail{
				TaskID:      task.ID,
				RecieverUID: item.UID,
				Content:     content,
				Status:      db.DM_DETAIL_STATUS_NOT_SEND,
			})
		}
	}

	// write to db
	if err := db.GetDirectMsgDAL().BatchCreateDetails(ctx, details); err != nil {
		logger.Error("save task details error %v", err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	// start dm job
	if req.RunTask {
		if err = batch_dm.GetManager().SetSenderInfo(task.ID, &req.Sender); err != nil {
			logger.Error("set sender info for task %d failed: %v", task.ID, err)
		} else {
			if err = batch_dm.GetManager().StartTask(ctx, task.ID); err != nil {
				logger.Error("start dm task %d failed: %v", task.ID, err)
			} else {
				logger.Info("start dm task %d", task.ID)
			}
		}
	}

	return &bs.Nothing{}, nil
}

func (ins dmHandler) page(ctx *swe.Context, req *bs.PageReq) (*bs.PageRsp, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	st, _ := session.GetStreamerSession(ctx)

	// query tasks
	count, tasks, err := db.GetDirectMsgDAL().Page(ctx, st.RoomID, (req.Page-1)*req.Size, req.Size)
	if err != nil {
		logger.Error("query tasks for room %d error %v", st.RoomID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	rsp := bs.PageRsp{Count: count, List: []any{}}

	// query events
	eids := make([]int64, 0, len(tasks))
	for _, item := range tasks {
		eids = append(eids, item.EventID)
	}
	events, err := db.GetRewardEventDAL().GetByIDs(ctx, eids)
	if err != nil {
		logger.Error("query events for tasks of room %d error %v", st.RoomID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	eventMap := map[int64]*db.RewardEvent{}
	for idx := range events {
		item := &events[idx]
		eventMap[item.ID] = item
	}

	for _, item := range tasks {
		tmp := bs.DMTaskListItem{
			ID:          item.ID,
			Name:        item.TaskName,
			Content:     item.Content,
			BatchMax:    item.BatchMax,
			IntervalMin: item.IntervalMin,
			IntervalMax: item.IntervalMax,
			Status:      item.Status,
		}
		event, ok := eventMap[item.EventID]
		if !ok {
			logger.Error("task %d event %d not found", item.ID, item.EventID)
			continue
		}
		tmp.Event.ID = event.ID
		tmp.Event.Name = event.EventName
		rsp.List = append(rsp.List, tmp)
	}

	return &rsp, nil
}

func (ins dmHandler) detail(ctx *swe.Context, req *bs.IDReq) (*bs.DMTaskDetail, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	st, _ := session.GetStreamerSession(ctx)

	task, err := db.GetDirectMsgDAL().GetByRoomID(ctx, req.ID, st.RoomID)
	if err != nil {
		logger.Error("query task %d for room %d error %v", req.ID, st.RoomID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if task == nil {
		logger.Error("query task %d for room %d not found", req.ID, st.RoomID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	ret := bs.DMTaskDetail{
		ID:          task.ID,
		Name:        task.TaskName,
		Content:     task.Content,
		BatchMax:    task.BatchMax,
		IntervalMin: task.IntervalMin,
		IntervalMax: task.IntervalMax,
		Status:      task.Status,
	}

	stats, err := db.GetDirectMsgDAL().Stats(ctx, req.ID)
	if err != nil {
		logger.Error("query stat for task %d error %v", req.ID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	ret.Total = stats.Fail + stats.Done + stats.NotSend
	ret.Fail = stats.Fail
	ret.Succ = stats.Done

	return &ret, nil
}

func (ins dmHandler) setSender(ctx *swe.Context, req *bs.DMSetSenderReq) (*bs.Nothing, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	st, _ := session.GetStreamerSession(ctx)

	task, err := db.GetDirectMsgDAL().GetByRoomID(ctx, req.TaskID, st.RoomID)
	if err != nil {
		logger.Error("query task %d for room %d error %v", req.TaskID, st.RoomID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if task == nil {
		logger.Error("query task %d for room %d not found", req.TaskID, st.RoomID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	err = batch_dm.GetManager().SetSenderInfo(req.TaskID, &req.Sender)
	if err != nil {
		logger.Error("set sender info for dm task %d failed: %v", req.TaskID, err)
		return nil, swe.Error(EC_DM_SET_SENDER_FAIL, err)
	}

	if req.RunTask {
		err = batch_dm.GetManager().StartTask(ctx, req.TaskID)
		if err != nil {
			logger.Error("start dm task %d failed: %v", req.TaskID, err)
			return nil, swe.Error(EC_DM_START_FAIL, err)
		}
	}

	return &bs.Nothing{}, nil
}

func (ins dmHandler) setState(ctx *swe.Context, req *bs.DMSwitchReq) (*bs.Nothing, swe.SweError) {
	logger := swe.CtxLogger(ctx)
	st, _ := session.GetStreamerSession(ctx)

	task, err := db.GetDirectMsgDAL().GetByRoomID(ctx, req.TaskID, st.RoomID)
	if err != nil {
		logger.Error("query task %d for room %d error %v", req.TaskID, st.RoomID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if task == nil {
		logger.Error("query task %d for room %d not found", req.TaskID, st.RoomID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}

	if req.RunTask {
		err = batch_dm.GetManager().StartTask(ctx, req.TaskID)
	} else {
		err = batch_dm.GetManager().StopTask(req.TaskID)
	}
	if err != nil {
		logger.Error("change dm task %d status failed: %v", req.TaskID, err)
		return nil, swe.Error(EC_DM_START_FAIL, err)
	}

	return &bs.Nothing{}, nil
}
