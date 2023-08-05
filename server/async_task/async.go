package async_task

import (
	"sync"
	"time"

	"github.com/zerozwt/octant/server/db"
	"github.com/zerozwt/octant/server/utils"
	"github.com/zerozwt/swe"
)

type Scheduler interface {
	AddTask(ctx *swe.Context, name, param string, ts int64, cb func(id int64)) error
}

type TaskContext interface {
	ID() int64
	Name() string
	Param() string
	ChangeSchedule(ts int64)
}

type Handler func(ctx *swe.Context, taskCtx TaskContext) error

var hMap map[string]Handler

func RegisterHandler(name string, handler Handler) {
	hMap[name] = handler
}

//---------------------------------------------------------------------------------

type scheduler struct {
	lock sync.Mutex

	queue *utils.PriorQueue[db.AsyncTask]
}

var gs *scheduler = &scheduler{}

func GetScheduler() Scheduler { return gs }

func Init() error { return gs.init() }

type context struct {
	task  *db.AsyncTask
	sched int64
}

func (ctx *context) ID() int64               { return ctx.task.ID }
func (ctx *context) Name() string            { return ctx.task.Handler }
func (ctx *context) Param() string           { return ctx.task.Param }
func (ctx *context) ChangeSchedule(ts int64) { ctx.sched = ts }

func (s *scheduler) AddTask(ctx *swe.Context, name, param string, ts int64, cb func(id int64)) error {
	task := &db.AsyncTask{
		ID:       utils.GenerateID(),
		Handler:  name,
		Param:    param,
		Status:   db.ASYNC_TASK_IDLE,
		Schedule: ts,
	}

	if err := db.GetAsyncTaskDAL().Put(ctx, task); err != nil {
		return err
	}

	if cb != nil {
		cb(task.ID)
	}

	s.lock.Lock()
	s.queue.Put(task)
	s.lock.Unlock()

	return nil
}

func (s *scheduler) init() error {
	// load all tasks from db
	tasks, err := db.GetAsyncTaskDAL().All(nil)
	if err != nil {
		return err
	}

	cmp := func(a, b *db.AsyncTask) bool { return a.Schedule < b.Schedule }
	s.queue = utils.PriorityQueue(cmp)

	for _, item := range tasks {
		s.queue.Put(item)
	}

	// start tick thread
	go s.tick()

	return nil
}

func (s *scheduler) tick() {
	for {
		<-time.After(time.Second)
		s.run()
	}
}

func (s *scheduler) run() {
	s.lock.Lock()
	defer s.lock.Unlock()

	now := time.Now().Unix()

	for {
		task := s.queue.Head()
		if task == nil || task.Schedule > now {
			return
		}

		ctx := &swe.Context{}
		swe.AssignLogID(ctx)
		logger := swe.CtxLogger(ctx)
		logger.SetRenderer(utils.LogRenderer())

		task = s.queue.Pop()
		logger.Info("start async task %d handler: %s", task.ID, task.Handler)

		taskCtx := context{task: task}

		handler, ok := hMap[task.Handler]
		if !ok {
			logger.Error("run async task %d failed: handler not found for %s", task.ID, task.Handler)
			continue
		}

		go s.runTask(handler, ctx, taskCtx)
	}
}

func (s *scheduler) runTask(handler Handler, ctx *swe.Context, taskCtx context) {
	err := handler(ctx, &taskCtx)
	logger := swe.CtxLogger(ctx)

	if taskCtx.sched > 0 {
		logger.Info("async task %d not done, scheduled at %d", taskCtx.task.ID, taskCtx.sched)
		taskCtx.task.Schedule = taskCtx.sched
		taskCtx.task.Status = db.ASYNC_TASK_IDLE
	} else if err != nil {
		logger.Error("async task %d failed: %v", taskCtx.task.ID, err)
		taskCtx.task.Status = db.ASYNC_TASK_FAILED
	} else {
		logger.Info("async task %d sucess", taskCtx.task.ID)
		taskCtx.task.Status = db.ASYNC_TASK_DONE
	}

	err = db.GetAsyncTaskDAL().Put(ctx, taskCtx.task)
	if err != nil {
		logger.Error("update async task %d to db failed: %v", taskCtx.task.ID, err)
	}

	if taskCtx.sched > 0 {
		s.lock.Lock()
		s.queue.Put(taskCtx.task)
		s.lock.Unlock()
	}
}
