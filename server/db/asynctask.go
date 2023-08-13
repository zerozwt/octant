package db

import (
	"github.com/zerozwt/swe"
	"gorm.io/gorm/clause"
)

type AsyncTask struct {
	ID       int64  `gorm:"primaryKey;column:id"`
	Handler  string `gorm:"type:string;size:128;column:handler"`
	Param    string `gorm:"type:string;size:4096;column:param"`
	Status   int    `gorm:"column:status;index:idx_ast_status"`
	Schedule int64  `gorm:"column:schedule"`
}

const (
	ASYNC_TASK_IDLE = iota + 1
	ASYNC_TASK_RUNNING
	ASYNC_TASK_FAILED
	ASYNC_TASK_DONE
)

func (s AsyncTask) TableName() string { return "t_async_task" }

func init() {
	registerModel(&AsyncTask{})
}

type AsyncTaskDAL struct{}

func GetAsyncTaskDAL() AsyncTaskDAL { return AsyncTaskDAL{} }

func (dal AsyncTaskDAL) All(ctx *swe.Context) ([]*AsyncTask, error) {
	tmp := []AsyncTask{}
	tx := getInstance(ctx).Where("status in ?", []int{ASYNC_TASK_IDLE, ASYNC_TASK_RUNNING})
	err := tx.Find(&tmp).Error
	ret := make([]*AsyncTask, 0, len(tmp))
	for idx := range tmp {
		ret = append(ret, &tmp[idx])
	}
	return ret, err
}

func (dal AsyncTaskDAL) Put(ctx *swe.Context, value *AsyncTask) error {
	return getInstance(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "schedule"}),
	}).Create(value).Error
}

func (dal AsyncTaskDAL) Get(ctx *swe.Context, id int64) (*AsyncTask, error) {
	ret := []AsyncTask{}
	err := getInstance(ctx).Where("id = ?", id).Find(&ret).Error
	if len(ret) == 0 {
		return nil, err
	}
	return &ret[0], nil
}
