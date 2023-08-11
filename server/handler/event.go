package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/zerozwt/octant/server/async_task"
	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/handler/event_calc"
	"github.com/zerozwt/octant/server/session"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

func init() {
	registerHandler(GET, "/event/list", event.list, session.CheckStreamer)
	registerHandler(POST, "/event/add", event.add, session.CheckStreamer)
	registerHandler(POST, "/event/modify", event.modify, session.CheckStreamer)
	registerHandler(GET, "/event/detail", event.detail, session.CheckStreamer)
	registerHandler(POST, "/event/delete", event.delete, session.CheckStreamer)

	async_task.RegisterHandler(asyncTaskCalculateEventList, event.calculate)

	registerHandler(POST, "/event/user/list", event.userList, session.CheckStreamer)
	registerHandler(POST, "/event/user/block", event.blockUser, session.CheckStreamer)
	registerHandler(POST, "/event/user/unblock", event.unblockUser, session.CheckStreamer)

	registerRawHandler(GET, "/api/event/user/dl", event.download, session.CheckStreamer)
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
	if event == nil {
		logger.Error("find event %d from db failed: not found", evtID)
		return fmt.Errorf("event not found")
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
	users := map[int64]*event_calc.UserData{}

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
				users[item.SenderUID] = event_calc.NewEventUser(item.SenderUID)
			}
			users[item.SenderUID].Gift = append(users[item.SenderUID].Gift, item)
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
				users[item.SenderUID] = event_calc.NewEventUser(item.SenderUID)
			}
			users[item.SenderUID].SC = append(users[item.SenderUID].SC, item)
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
				users[item.SenderUID] = event_calc.NewEventUser(item.SenderUID)
			}
			users[item.SenderUID].Member = append(users[item.SenderUID].Member, item)
		}
	}

	logger.Info("filtering data for event %d", evtID)

	// filter sender
	filter := event_calc.BuildFilter(&cond)
	tmp := users
	users = map[int64]*event_calc.UserData{}
	for uid, data := range tmp {
		strip := event_calc.NewEventUserStrip()
		if filter.OK(data, strip) {
			users[uid] = data.Strip(strip)
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
		return users[uids[i]].SendTs < users[uids[j]].SendTs
	})

	data := make([]db.RewardUser, 0, len(users))
	for _, uid := range uids {
		data = append(data, db.RewardUser{
			EventID:  evtID,
			UID:      uid,
			UserName: users[uid].Name,
			Time:     users[uid].SendTs,
			Columns:  users[uid].Column(),
		})
	}

	if err = db.GetRewardEventDAL().PutUsers(ctx, data); err != nil {
		db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
		logger.Error("write new list for event %d failed: %v", evtID, err)
		return err
	}

	// create dd accounts
	nowTs := time.Now().Unix()
	dd := make([]db.DDInfo, 0, len(users))
	for _, uid := range uids {
		dd = append(dd, db.DDInfo{
			UID:        uid,
			UserName:   users[uid].Name,
			AccessCode: db.GetDDInfoDAL().GenerateAccessCode(nowTs, evtID, uid),
		})
	}
	if err = db.GetDDInfoDAL().BatchCreate(ctx, dd); err != nil {
		db.GetRewardEventDAL().SetStatus(ctx, evtID, db.EVENT_ERROR)
		logger.Error("create dd accounts for event %d failed: %v", evtID, err)
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

func (ins eventHandler) modify(ctx *swe.Context, req *bs.EventModifyReq) (*bs.Nothing, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	if err := db.GetRewardEventDAL().UpdateEventInfo(ctx, req.ID, st.RoomID, req.Name, req.Reward); err != nil {
		swe.CtxLogger(ctx).Error("update info for event %d error %v", req.ID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	return &bs.Nothing{}, nil
}

func (ins eventHandler) detail(ctx *swe.Context, req *bs.IDReq) (*bs.EventDetailRsp, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	event, err := db.GetRewardEventDAL().GetByRoomID(ctx, req.ID, st.RoomID)
	if err != nil {
		swe.CtxLogger(ctx).Error("query db for event %d error %v", req.ID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if event == nil {
		swe.CtxLogger(ctx).Error("query db for event %d room id %d not found", req.ID, st.RoomID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("event not found"))
	}

	ret := &bs.EventDetailRsp{
		ID:     req.ID,
		Name:   event.EventName,
		Reward: event.RewardContent,
		Status: event.Status,
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.UnmarshalFromString(event.Conditions, &ret.Condition)
	if err != nil {
		swe.CtxLogger(ctx).Error("decode condition for event %d error %v", req.ID, err)
		return nil, swe.Error(EC_EVT_COND_DECODE_FAIL, err)
	}

	return ret, nil
}

func (ins eventHandler) delete(ctx *swe.Context, req *bs.IDReq) (*bs.Nothing, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	rows, err := db.GetRewardEventDAL().Delete(ctx, req.ID, st.RoomID)
	if err != nil {
		swe.CtxLogger(ctx).Error("delete event %d error %v", req.ID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if rows > 0 {
		err = db.GetRewardEventDAL().ClearUsers(ctx, req.ID)
		if err != nil {
			swe.CtxLogger(ctx).Error("clear user list for event %d error %v", req.ID, err)
		}
	}
	return &bs.Nothing{}, nil
}

func (ins eventHandler) userList(ctx *swe.Context, req *bs.EventUserListReq) (*bs.PageRsp, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	exist, err := db.GetRewardEventDAL().Exist(ctx, req.EventID, st.RoomID)
	if err != nil {
		swe.CtxLogger(ctx).Error("query event %d error %v", req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if !exist {
		swe.CtxLogger(ctx).Error("query event %d not exist", req.EventID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("event not found"))
	}

	count, users, err := db.GetRewardEventDAL().UserPage(ctx, req.EventID, (req.Page-1)*req.Size, req.Size)
	if err != nil {
		swe.CtxLogger(ctx).Error("query user list for event %d error %v", req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	ret := &bs.PageRsp{Count: count, List: []any{}}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	for _, item := range users {
		user := bs.EventUserListItem{
			UID:   item.UID,
			Name:  item.UserName,
			Time:  utils.TimeToCSTString(item.Time),
			Cols:  map[string]any{},
			Block: item.Blocked != 0,
		}
		json.UnmarshalFromString(item.Columns, &user.Cols)
		ret.List = append(ret.List, user)
	}

	return ret, nil
}

func (ins eventHandler) blockUser(ctx *swe.Context, req *bs.EventUIDReq) (*bs.Nothing, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	exist, err := db.GetRewardEventDAL().Exist(ctx, req.EventID, st.RoomID)
	if err != nil {
		swe.CtxLogger(ctx).Error("query event %d error %v", req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if !exist {
		swe.CtxLogger(ctx).Error("query event %d not exist", req.EventID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("event not found"))
	}

	ok, err := db.GetRewardEventDAL().BlockUser(ctx, req.EventID, req.UID, true)
	if err != nil {
		swe.CtxLogger(ctx).Error("block user %d for %d error %v", req.UID, req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if !ok {
		swe.CtxLogger(ctx).Error("block user %d for %d failed: user not found", req.UID, req.EventID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("user not found"))
	}

	return &bs.Nothing{}, nil
}

func (ins eventHandler) unblockUser(ctx *swe.Context, req *bs.EventUIDReq) (*bs.Nothing, swe.SweError) {
	st, _ := session.GetStreamerSession(ctx)
	exist, err := db.GetRewardEventDAL().Exist(ctx, req.EventID, st.RoomID)
	if err != nil {
		swe.CtxLogger(ctx).Error("query event %d error %v", req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if !exist {
		swe.CtxLogger(ctx).Error("query event %d not exist", req.EventID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("event not found"))
	}

	ok, err := db.GetRewardEventDAL().BlockUser(ctx, req.EventID, req.UID, false)
	if err != nil {
		swe.CtxLogger(ctx).Error("block user %d for %d error %v", req.UID, req.EventID, err)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, err)
	}
	if !ok {
		swe.CtxLogger(ctx).Error("block user %d for %d failed: user not found", req.UID, req.EventID)
		return nil, swe.Error(EC_GENERIC_DB_FAIL, fmt.Errorf("user not found"))
	}

	return &bs.Nothing{}, nil
}

func (ins eventHandler) download(ctx *swe.Context) {
	logger := swe.CtxLogger(ctx)
	req := bs.IDReq{}
	if err := swe.DecodeForm(ctx.Request, &req); err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		ctx.Response.Write([]byte(err.Error()))
		return
	}

	// query event
	st, _ := session.GetStreamerSession(ctx)
	event, err := db.GetRewardEventDAL().GetByRoomID(ctx, req.ID, st.RoomID)
	if err != nil {
		logger.Error("query event %d error %v", req.ID, err)
		http.NotFound(ctx.Response, ctx.Request)
		return
	}
	if event == nil {
		logger.Error("query event %d not exist", req.ID)
		http.NotFound(ctx.Response, ctx.Request)
		return
	}
	if event.Status != db.EVENT_READY {
		logger.Error("event %d not ready", req.ID)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		ctx.Response.Write([]byte(`event not ready`))
		return
	}

	// get users
	users, err := db.GetRewardEventDAL().Users(ctx, req.ID)
	if err != nil {
		logger.Error("load users for event %d failed: %v", req.ID, err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.Response.Write([]byte(err.Error()))
		return
	}

	// get user public keys
	uids := make([]int64, 0, len(users))
	for _, item := range users {
		uids = append(uids, item.UID)
	}
	keyMap, err := db.GetDDInfoDAL().GetPublicKeys(ctx, uids)
	if err != nil {
		logger.Error("query public keys for users in event %d failed: %v", req.ID, err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.Response.Write([]byte(err.Error()))
		return
	}

	// convert users & decrypt user address
	userDatas := make([]*event_calc.UserData, 0, len(users))
	addrs := event_calc.AddrMap{}
	for _, item := range users {
		data, err := event_calc.EventUserfromDB(&item)
		if err != nil {
			logger.Error("convert user %d for event %d failed: %v, skip user ...", item.UID, req.ID, err)
			continue
		}
		userDatas = append(userDatas, data)
		if key, ok := keyMap[item.UID]; ok {
			addr, err := db.DecryptUserAddress(ctx, key, item.AddressInfo)
			if err != nil {
				logger.Error("decrypt address info for user %d in event %d failed: %v", item.UID, req.ID, err)
			} else {
				addrs[item.UID] = addr
			}
		}
	}
	ctx.Put(event_calc.CTX_KEY_ADDR, addrs)

	// decode condition
	cond := bs.EventCondition{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.UnmarshalFromString(event.Conditions, &cond)
	if err != nil {
		logger.Error("decode condition failed: %v", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.Response.Write([]byte(err.Error()))
		return
	}
	err = cond.Validate(ctx)
	if err != nil {
		logger.Error("validate condition failed: %v", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.Response.Write([]byte(err.Error()))
		return
	}

	// generate csv lines
	lines := event_calc.Table(ctx, userDatas, event_calc.BuildPickers(ctx, &cond))
	csvData := bytes.Buffer{}
	csvData.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF8 BOM
	err = csv.NewWriter(&csvData).WriteAll(lines)
	if err != nil {
		logger.Error("write csv data failed: %v", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.Response.Write([]byte(err.Error()))
		return
	}

	// generate file name
	fileName := filterFileName(strings.Join([]string{st.StreamerName, event.EventName}, "_"))

	// set header & write csv data
	ctx.Response.Header().Set("Content-Type", "application/octet-stream")
	ctx.Response.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.csv"`, fileName))
	ctx.Response.Write(csvData.Bytes())
}

func filterFileName(value string) string {
	ret := strings.Builder{}

	for _, ch := range value {
		switch ch {
		case '<', '>', ':', '"', '\'', '/', '\\', '|', '?', '*':
			ret.WriteByte('_')
		default:
			ret.WriteRune(ch)
		}
	}

	return ret.String()
}
