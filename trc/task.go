package trc

import (
	"time"
)

// 任务接口
type ITask interface {
	// 获取代码
	GetCode() string
	// 获取开始时间
	GetBeginTime() time.Time
	// 获取结束时间
	GetEndTime() time.Time
	// 获取数量
	GetCount() int
}
