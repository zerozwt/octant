package batch_dm

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	dm "github.com/zerozwt/BLiveDanmaku"
	"github.com/zerozwt/octant/server/bs"
	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

var ErrSenderNotSet error = fmt.Errorf("sender info not set")
var ErrTaskNotFound error = fmt.Errorf("task not found")

type Manager interface {
	StartTask(ctx *swe.Context, id int64) error
	StopTask(id int64) error
	SetSenderInfo(id int64, sender *bs.DMSenderInfo) error
}

func GetManager() Manager {
	return gm
}

type manager struct {
	taskSender map[int64]*bs.DMSenderInfo
	tasks      map[int64]*executor

	lock sync.Mutex
}

var gm *manager = &manager{
	taskSender: map[int64]*bs.DMSenderInfo{},
	tasks:      map[int64]*executor{},
}

func (m *manager) SetSenderInfo(id int64, sender *bs.DMSenderInfo) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.taskSender[id] = &bs.DMSenderInfo{
		UID:      sender.UID,
		SessData: sender.SessData,
		JCT:      sender.JCT,
	}

	return nil
}

func (m *manager) StartTask(ctx *swe.Context, id int64) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.tasks[id]; ok {
		return nil
	}

	if info, ok := m.taskSender[id]; !ok || info.Validate(ctx) != nil {
		return ErrSenderNotSet
	}

	devID, err := dm.GetDMDeviceID()
	if err != nil {
		return err
	}

	exe := &executor{
		taskID: id,
		sender: m.taskSender[id],
		devID:  devID,
	}

	if err := exe.start(ctx); err != nil {
		return err
	}

	m.tasks[id] = exe

	return nil
}

func (m *manager) StopTask(id int64) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if exe, ok := m.tasks[id]; ok {
		exe.stop()
		delete(m.tasks, id)
	}

	return nil
}

type executor struct {
	taskID int64
	sender *bs.DMSenderInfo
	task   *db.DMTask
	devID  string

	stopped atomic.Bool
}

func (e *executor) start(ctx *swe.Context) error {
	logger := swe.CtxLogger(ctx)

	// load task
	task, err := db.GetDirectMsgDAL().Get(ctx, e.taskID)
	if err != nil {
		logger.Error("load direct msg task %d failed: %v", e.taskID, err)
		return err
	}
	if task == nil {
		logger.Error("load direct msg task %d failed: not found", e.taskID)
		return ErrTaskNotFound
	}

	// change status
	if err = db.GetDirectMsgDAL().UpdateStatus(ctx, e.taskID, db.DM_TASK_STATUS_RUNNING); err != nil {
		logger.Error("update direct msg task %d status to running failed: %v", e.taskID, err)
		return err
	}

	// create new context
	exeCtx := &swe.Context{}
	swe.AssignLogID(exeCtx)
	swe.CtxLogger(exeCtx).SetRenderer(utils.LogRenderer())

	// start exec
	logger.Info("start direct msg task %d log id %s", e.taskID, swe.CtxLogID(exeCtx))
	e.task = task
	go e.exec(exeCtx)

	return nil
}

func (e *executor) stop() error {
	e.stopped.Store(true)
	return nil
}

func (e *executor) exec(ctx *swe.Context) {
	defer gm.StopTask(e.taskID)

	logger := swe.CtxLogger(ctx)

	// load details
	details, err := db.GetDirectMsgDAL().LoadUnsentDetails(ctx, e.taskID, e.task.BatchMax)
	if err != nil {
		logger.Error("load details for task %d failed: %v", e.taskID, err)
		db.GetDirectMsgDAL().UpdateStatus(ctx, e.taskID, db.DM_TASK_STATUS_PAUSED)
		return
	}
	logger.Info("load %d details for task %d", len(details), e.taskID)

	failCount := 0

	// exec details
	for idx, item := range details {
		if idx > 0 {
			e.randomSleep()
		}

		// check if failed several times
		if failCount >= 5 {
			logger.Error("direct msg task %d failed too much, stop", e.taskID)
			db.GetDirectMsgDAL().UpdateStatus(ctx, e.taskID, db.DM_TASK_STATUS_PAUSED)
			return
		}

		// check if stopped from outside
		if e.stopped.Load() {
			logger.Info("direct msg task %d interruped from manager, %d details processed", e.taskID, idx)
			db.GetDirectMsgDAL().UpdateStatus(ctx, e.taskID, db.DM_TASK_STATUS_PAUSED)
			return
		}

		// send dm
		logger.Info("[%d/%d] sending direct message to uid %d ...", idx+1, len(details), item.RecieverUID)
		rsp, err := dm.SendDirectMsg(e.sender.UID, item.RecieverUID, item.Content, e.devID, e.sender.SessData, e.sender.JCT)
		err = e.processDirectMsgRsp(rsp, err)
		if err != nil {
			logger.Error("send dm to %d failed: %v", item.RecieverUID, err)
			failCount += 1
			db.GetDirectMsgDAL().UpdateDetailStatus(ctx, e.taskID, item.RecieverUID, db.DM_DETAIL_STATUS_FAIL, err.Error())
			continue
		}

		// update detail status
		failCount = 0
		logger.Info("[%d/%d] sending direct message to uid %d succeed", idx+1, len(details), item.RecieverUID)
		db.GetDirectMsgDAL().UpdateDetailStatus(ctx, e.taskID, item.RecieverUID, db.DM_DETAIL_STATUS_DONE, "")
	}

	// get stat and update task status
	stat, err := db.GetDirectMsgDAL().Stats(ctx, e.taskID)
	if err != nil {
		logger.Error("direct msg task %d stat failed: %v", e.taskID, err)
		db.GetDirectMsgDAL().UpdateStatus(ctx, e.taskID, db.DM_TASK_STATUS_PAUSED)
		return
	}

	if stat.NotSend == 0 && stat.Fail == 0 {
		logger.Info("direct msg task %d all done", e.taskID)
		db.GetDirectMsgDAL().UpdateStatus(ctx, e.taskID, db.DM_TASK_STATUS_DONE)
	} else {
		logger.Info("direct msg task %d partially done: done[%d] fail[%d] not_sent[%d]",
			e.taskID, stat.Done, stat.Fail, stat.NotSend)
		db.GetDirectMsgDAL().UpdateStatus(ctx, e.taskID, db.DM_TASK_STATUS_PAUSED)
	}
}

func (e *executor) randomSleep() {
	lower := time.Second * time.Duration(e.task.IntervalMin)
	upper := time.Second * time.Duration(e.task.IntervalMax)
	if lower < upper {
		lower += time.Duration(rand.Int63n(int64(upper - lower)))
	}
	time.Sleep(lower)
}

func (e *executor) processDirectMsgRsp(rsp *dm.SendDirectMsgRsp, err error) error {
	if err != nil {
		return err
	}
	if rsp.Code != 0 {
		return fmt.Errorf("[%d] %s", rsp.Code, rsp.Message)
	}
	return nil
}
