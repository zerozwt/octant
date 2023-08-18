package db

import (
	"strings"
	"time"

	"github.com/zerozwt/swe"
)

type DMTask struct {
	ID          int64  `gorm:"primaryKey;column:id"`
	RoomID      int64  `gorm:"column:room_id;index:idx_dmtask_room"`
	EventID     int64  `gorm:"column:event_id"`
	TaskName    string `gorm:"column:task_name;type:string;size:256"`
	MsgType     int    `gorm:"column:msg_type"`
	Content     string `gorm:"column:content;type:string;size:4096"`
	BatchMax    int    `gorm:"column:batch_max"`
	Status      int    `gorm:"column:status"`
	IntervalMin int    `gorm:"column:interval_min"`
	IntervalMax int    `gorm:"column:interval_max"`
	CreateTime  int64  `gorm:"create_time"`
}

func (s DMTask) TableName() string { return "t_dm_task" }

const (
	DM_TASK_STATUS_PAUSED = iota + 1
	DM_TASK_STATUS_RUNNING
	DM_TASK_STATUS_DONE
)

type DMDetail struct {
	TaskID      int64  `gorm:"column:task_id;index:idx_dm_task"`
	RecieverUID int64  `gorm:"column:uid;index:idx_dm_task"`
	Content     string `gorm:"column:content;type:TEXT"`
	Status      int    `gorm:"column:status"`
	SendTime    int64  `gorm:"column:send_ts"`
	FailReason  string `gorm:"column:fail_reason;type:string;size:1024"`
}

func (s DMDetail) TableName() string { return "t_dm_detail" }

type DMTaskStat struct {
	NotSend int
	Fail    int
	Done    int
}

const (
	DM_DETAIL_STATUS_NOT_SEND = iota + 1
	DM_DETAIL_STATUS_FAIL
	DM_DETAIL_STATUS_DONE
)

func init() {
	registerModel(&DMTask{})
	registerModel(&DMDetail{})
}

type DirectMsgDAL struct{}

func GetDirectMsgDAL() DirectMsgDAL { return DirectMsgDAL{} }

func (dal DirectMsgDAL) BatchCreateDetails(ctx *swe.Context, details []DMDetail) error {
	return getInstance(ctx).CreateInBatches(details, 500).Error
}

func (dal DirectMsgDAL) Put(ctx *swe.Context, task *DMTask) error {
	return getInstance(ctx).Create(task).Error
}

func (dal DirectMsgDAL) Page(ctx *swe.Context, roomID int64, offset, limit int) (int, []DMTask, error) {
	count := 0
	ret := []DMTask{}
	tx := getInstance(ctx).Where("room_id = ?", roomID)

	err := newDBSession(ctx, tx).Select("count(*)").Scan(&count).Error
	if err != nil {
		return 0, nil, err
	}

	err = tx.Offset(offset).Limit(limit).Order("create_time desc").Find(&ret).Error

	return count, ret, err
}

func (dal DirectMsgDAL) GetByRoomID(ctx *swe.Context, id, roomID int64) (*DMTask, error) {
	ret := []DMTask{}
	err := getInstance(ctx).Where("id = ? and room_id = ?", id, roomID).Find(&ret).Error
	if len(ret) == 0 {
		return nil, err
	}
	return &ret[0], err
}

func (dal DirectMsgDAL) Get(ctx *swe.Context, id int64) (*DMTask, error) {
	ret := []DMTask{}
	err := getInstance(ctx).Where("id = ?", id).Find(&ret).Error
	if len(ret) == 0 {
		return nil, err
	}
	return &ret[0], err
}

func (dal DirectMsgDAL) Stats(ctx *swe.Context, id int64) (DMTaskStat, error) {
	ret := DMTaskStat{}

	tmp := []struct {
		Status int `gorm:"column:status"`
		Count  int `gorm:"column:ct"`
	}{}

	tx := getInstance(ctx).Table((DMDetail{}).TableName())
	tx = tx.Select("status, count(*) as ct")
	tx = tx.Where("task_id = ?", id)
	tx = tx.Group("status")
	err := tx.Find(&tmp).Error

	for _, item := range tmp {
		switch item.Status {
		case DM_DETAIL_STATUS_NOT_SEND:
			ret.NotSend = item.Count
		case DM_DETAIL_STATUS_DONE:
			ret.Done = item.Count
		case DM_DETAIL_STATUS_FAIL:
			ret.Fail = item.Count
		}
	}

	return ret, err
}

func (dal DirectMsgDAL) UpdateStatus(ctx *swe.Context, id int64, status int) error {
	return getInstance(ctx).Exec("update t_dm_task set status = ? where id = ?", status, id).Error
}

func (dal DirectMsgDAL) LoadUnsentDetails(ctx *swe.Context, id int64, limit int) ([]DMDetail, error) {
	ret := []DMDetail{}

	tx := getInstance(ctx).Where("task_id = ?", id)
	tx = tx.Where("status <> ?", DM_DETAIL_STATUS_NOT_SEND)
	tx = tx.Order("uid")
	if limit > 0 {
		tx = tx.Limit(limit)
	}

	err := tx.Find(&ret).Error

	return ret, err
}

func (dal DirectMsgDAL) UpdateDetailStatus(ctx *swe.Context, id, uid int64, status int, reason string) error {
	fields := []string{"status", "send_ts"}
	params := []any{status, time.Now().Unix()}
	if len(reason) > 0 {
		fields = append(fields, "fail_reason")
		params = append(params, reason)
	}
	params = append(params, id, uid)

	builder := strings.Builder{}
	builder.WriteString("update t_dm_detail set")
	for idx, col := range fields {
		if idx > 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune(' ')
		builder.WriteString(col)
		builder.WriteString(" = ?")
	}
	builder.WriteString(" where task_id = ? and uid = ?")

	return getInstance(ctx).Exec(builder.String(), params...).Error
}
