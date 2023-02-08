package trc

import (
	"context"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

const (
	StateCanceled = "Canceled"
	StateStarted  = "Started"
)

// 分时量比控制
type Trc[T ITask] struct {
	cache      *redis.Client // 缓存
	hostname   string        // 主机名
	lastTasks  map[string]T  // 上次任务集合
	prefix     string        // 前缀
	timeshares [24]int       // 分时集合
}

// 构造函数
func NewTrc[T ITask](cache *redis.Client, prefix string, timeshares [24]int) *Trc[T] {
	hostname, _ := os.Hostname()
	stateKey := fmt.Sprintf("%s:*:state", prefix)
	stateKeys := cache.Keys(context.Background(), stateKey).Val()
	for _, k := range stateKeys {
		v := cache.Get(context.Background(), k).Val()
		if strings.HasPrefix(v, hostname) {
			cache.Del(context.Background(), k)
		}
	}

	return &Trc[T]{
		cache:      cache,
		prefix:     prefix,
		timeshares: timeshares,
		hostname:   hostname,
	}
}

// 获取状态
func (t *Trc[T]) GetState(task T) string {
	stateKey := fmt.Sprintf("%s:%s:state", t.prefix, task.GetCode())
	stateValue := t.cache.Get(context.Background(), stateKey).Val()
	state := strings.TrimPrefix(stateValue, t.hostname)

	return state
}

// 运行
func (t *Trc[T]) Run(tasks []T, process func(T) error) error {
	// 停止任务
	thisTasks := make(map[string]T)
	for _, task := range tasks {
		thisTasks[task.GetCode()] = task
	}
	for _, lt := range t.lastTasks {
		_, ok := thisTasks[lt.GetCode()]
		if !ok {
			t.cancel(lt)
		}
	}
	// 匹配任务
	matchedTasks := make(map[string]T)
	for _, tt := range thisTasks {
		lt, ok := t.lastTasks[tt.GetCode()]
		if ok && (reflect.DeepEqual(lt, tt) || t.cancel(lt)) {
			continue
		}
		if t.isAvailable(tt) {
			matchedTasks[tt.GetCode()] = tt
		}
	}
	// 运行任务
	for _, mt := range matchedTasks {
		go func(task T) {
			started := t.start(task)
			logrus.WithField("code", task.GetCode()).Debug("Start ", started, ". ", task)
			if started {
				defer func() {
					stopped := t.stop(task)
					logrus.WithField("code", task.GetCode()).Debug("Stop ", stopped, ". ", task)
				}()

				err := process(task)
				if err != nil {
					logrus.WithField("code", task.GetCode()).Error(err)
				}
			}
		}(mt)
	}
	// 归档任务
	t.lastTasks = matchedTasks

	return nil
}

// 计算百分比
func (t *Trc[T]) calcPercent(beginTime, endTime time.Time) float64 {
	percent := 0.0
	timeshares := append(t.timeshares[:], t.timeshares[:]...)
	hourCount := int(math.Ceil(endTime.Sub(beginTime).Hours()))
	hours := timeshares[beginTime.Hour() : beginTime.Hour()+hourCount]
	beginMinute, endMinute := float64(beginTime.Minute()), float64(endTime.Minute())
	for i, j := 0, hourCount-1; i <= j; i++ {
		hourPercent := float64(hours[i])
		if i == 0 {
			if i < j {
				percent += hourPercent * (60 - beginMinute) / 60
			} else {
				percent += hourPercent * (endMinute - beginMinute) / 60
			}
		} else if i < j {
			percent += hourPercent
		} else {
			percent += hourPercent * endMinute / 60
		}
	}

	return percent
}

// 计算比率
func (t *Trc[T]) calcRate(beginTime, currentTime, endTime time.Time) float64 {
	currentPercent := t.calcPercent(beginTime, currentTime)
	totalPercent := t.calcPercent(beginTime, endTime)
	rate := cast.ToFloat64(fmt.Sprintf("%.4f", currentPercent/totalPercent))

	return rate
}

// 取消
func (t *Trc[T]) cancel(task T) bool {
	stateKey := fmt.Sprintf("%s:%s:state", t.prefix, task.GetCode())
	stateValue := fmt.Sprintf("%s.%s", t.hostname, StateCanceled)
	canceled := t.cache.SetXX(context.Background(), stateKey, stateValue, 24*time.Hour).Val()

	return canceled
}

// 删除状态
func (t *Trc[T]) delState(task T) bool {
	stateKey := fmt.Sprintf("%s:%s:state", t.prefix, task.GetCode())
	deleted := t.cache.Del(context.Background(), stateKey).Val() > 0

	return deleted
}

// 是否可用
func (t *Trc[T]) isAvailable(task T) bool {
	expectedRate := t.calcRate(task.GetBeginTime(), time.Now(), task.GetEndTime())
	expectedCount := int64(float64(task.GetCount()) * expectedRate)
	timesharesKey := fmt.Sprintf("%s:%s:%s:timeshares", t.prefix, task.GetBeginTime().Format("20060102"), task.GetCode())
	expectedKey := "expected"
	t.cache.HSet(context.Background(), timesharesKey, expectedKey, expectedCount)
	actualKey := "actual"
	actualCount := cast.ToInt64(t.cache.HGet(context.Background(), timesharesKey, actualKey).Val())
	available := actualCount < expectedCount

	return available
}

// 启动
func (t *Trc[T]) start(task T) bool {
	stateKey := fmt.Sprintf("%s:%s:state", t.prefix, task.GetCode())
	stateValue := fmt.Sprintf("%s.%s", t.hostname, StateStarted)
	if !t.cache.SetNX(context.Background(), stateKey, stateValue, 24*time.Hour).Val() {
		return false
	}
	timesharesKey := fmt.Sprintf("%s:%s:%s:timeshares", t.prefix, task.GetBeginTime().Format("20060102"), task.GetCode())
	expectedKey := "expected"
	expectedCount := cast.ToInt64(t.cache.HGet(context.Background(), timesharesKey, expectedKey).Val())
	actualKey := "actual"
	if t.cache.HIncrBy(context.Background(), timesharesKey, actualKey, 1).Val() > expectedCount {
		t.cache.HIncrBy(context.Background(), timesharesKey, actualKey, -1)
		return false
	}

	return true
}

// 停止
func (t *Trc[T]) stop(task T) bool {
	state := t.GetState(task)
	t.delState(task)
	timesharesKey := fmt.Sprintf("%s:%s:%s:timeshares", t.prefix, task.GetBeginTime().Format("20060102"), task.GetCode())
	field := "finished"
	if state == StateCanceled {
		field = "canceled"
	}
	t.cache.HIncrBy(context.Background(), timesharesKey, field, 1)

	return true
}
