package scheduler

import (
	"context"
	"time"

	"github.com/yb2020/odoc/pkg/logging"
)

// SimpleJob 是一个简单的任务案例，它使用cron表达式来调度任务，并使用锁来防止多个实例同时执行任务。
type SimpleJob struct {
	logger   logging.Logger
	spec     string        // 任务的cron表达式，6字段标准cron表达式
	key      string        // 任务的锁key，必须是唯一的unique-job-key
	expiry   time.Duration // 任务的锁过期时间
	lockOpts *LockOptions  // 任务的锁选项
}

func NewSimpleJob(logger logging.Logger, spec string, key string, expiry time.Duration) *SimpleJob {
	lockOpts := &LockOptions{
		Key:    key,
		Expiry: expiry,
	}
	return &SimpleJob{logger: logger, spec: spec, key: key, expiry: expiry, lockOpts: lockOpts}
}

func (j *SimpleJob) NewUserContext(ctx context.Context, userId int64) context.Context {
	return ctx
}

// Spec 获取任务的cron表达式
func (j *SimpleJob) Spec() string {
	return j.spec
}

// LockOpts 获取任务的锁选项
func (j *SimpleJob) LockOpts() *LockOptions {
	return j.lockOpts
}

// Run 执行任务，在执行任务前会获取锁，执行任务后会释放锁
// 这是必须实现的接口，否则实现ScheduledJob接口的调度器无法正确调度任务
func (j *SimpleJob) Run() {
	j.logger.Info("msg", "Running simple job")
}
