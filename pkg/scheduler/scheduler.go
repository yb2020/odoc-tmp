package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/yb2020/odoc/pkg/logging"
)

// Scheduler 调度器
// 调度器用于执行定时任务
// 调度器使用robfig/cron库来实现定时任务，并支持基于Redis的分布式锁
type Scheduler struct {
	logger  logging.Logger
	cron    *cron.Cron
	redis   *redis.Client
	redsync *redsync.Redsync
}

// NewScheduler 创建一个新的调度器
// 如果提供了Redis客户端，则调度器将支持分布式锁
// NewScheduler 创建一个新的调度器
// redisClient is assumed to be an initialized *redis.Client instance.
func NewScheduler(logger logging.Logger, redisClient *redis.Client) *Scheduler {
	return &Scheduler{
		logger: logger,
		redis:  redisClient,
	}
}

// Init 初始化调度器，包括cron实例和redsync(如果提供了redis)
func (s *Scheduler) Init() error {
	s.logger.Info("msg", "Initializing scheduler...")
	s.cron = cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			// 捕获job中的panic，防止整个调度器崩溃
			cron.Recover(cron.DefaultLogger),
		),
	)

	if s.redis != nil {
		pool := goredis.NewPool(s.redis)
		s.redsync = redsync.New(pool)
		s.logger.Info("msg", "Redsync for distributed locks initialized successfully.")
	} else {
		s.logger.Warn("msg", "Redis client not provided to scheduler, distributed lock feature will be disabled.")
	}

	return nil
}

// AddFunc 向Cron添加一个要在给定计划上运行的函数。
// 这是 AddJob 的一个方便的包装器。
func (s *Scheduler) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return s.AddJob(spec, cron.FuncJob(cmd), nil)
}

// LockOptions 定义作业的分布式锁参数。
type LockOptions struct {
	// Key 是锁的唯一键。
	Key string
	// Expiry 是锁应该被持有的时间。应大于任务的最大执行时间。
	Expiry time.Duration
}

// AddJob 向Cron添加一个要在给定计划上运行的Job。
// 如果提供了lockOptions，作业将被包装在分布式锁中。如果lockOptions为nil，作业将是一个普通的Job。
func (s *Scheduler) AddJob(spec string, job cron.Job, lockOptions *LockOptions) (cron.EntryID, error) {
	jobToSchedule := job
	logMessage := "Scheduled a job"
	logFields := []interface{}{"spec", spec}

	if lockOptions != nil {
		if s.redsync == nil {
			err := errors.New("redsync is not initialized, cannot add locked job")
			s.logger.Error("msg", "Failed to add locked job", "error", err)
			return 0, err
		}
		wrapper := &DistributedLockJobWrapper{
			LockKey:       lockOptions.Key,
			UnderlyingJob: job,
			Redsync:       s.redsync,
			Logger:        s.logger,
			LockExpiry:    lockOptions.Expiry,
		}
		jobToSchedule = wrapper
		logMessage = "Scheduled a locked job"
		logFields = append(logFields, "lockKey", lockOptions.Key)
	}

	entryID, err := s.cron.AddJob(spec, jobToSchedule)
	if err != nil {
		s.logger.Error("msg", "Failed to add job", "error", err, "spec", spec)
		return 0, err
	}

	logFields = append(logFields, "entryID", entryID)
	finalLogFields := append([]interface{}{"msg", logMessage}, logFields...)
	s.logger.Info(finalLogFields...)

	return entryID, nil
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	s.logger.Info("msg", "Starting scheduler...")
	s.cron.Start()
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() error {
	s.logger.Info("msg", "Stopping scheduler...")
	// Stop会平滑地停止调度器，它会等待所有正在运行的任务完成
	ctx := s.cron.Stop()
	select {
	case <-ctx.Done():
		s.logger.Info("msg", "Scheduler gracefully stopped.")
		return nil
	case <-time.After(15 * time.Second): // 设置一个最长等待时间
		s.logger.Error("msg", "Scheduler stop timed out, forcing shutdown.")
		return errors.New("scheduler stop timed out")
	}
}

// ScheduledJob 封装了一个待调度的任务及其所有配置。
// 这种结构便于以声明式的方式管理和注册多个任务。
type ScheduledJob interface {
	NewUserContext(ctx context.Context, userId string) context.Context
	Spec() string
	LockOpts() *LockOptions
	cron.Job
}

func (s *Scheduler) RegisterJobs(jobs ...ScheduledJob) {
	for _, sj := range jobs {
		if sj == nil {
			s.logger.Warn("msg", "Skipping nil or invalid scheduled job")
			continue
		}
		_, err := s.AddJob(sj.Spec(), sj, sj.LockOpts())
		if err != nil {
			s.logger.Error("msg", "Failed to register a scheduled job", "spec", sj.Spec(), "error", err)
		}
	}
}

// --- Distributed Lock Wrapper ---

// DistributedLockJobWrapper 是一个cron.Job的包装器，它在执行实际任务前获取分布式锁。
type DistributedLockJobWrapper struct {
	LockKey       string           // LockKey 是用于分布式锁的唯一键名。
	UnderlyingJob cron.Job         // UnderlyingJob 是被包装的实际要执行的 cron 任务。
	Redsync       *redsync.Redsync // Redsync 是 redsync 库的实例，用于处理 Redis 分布式锁的获取和释放。
	Logger        logging.Logger   // Logger 用于记录与锁相关的操作日志，例如获取锁、释放锁、错误等。
	LockExpiry    time.Duration    // LockExpiry 定义了锁的有效时间（租约期）。一旦超过这个时间，即使任务未完成，锁也会自动释放，以防止死锁。
}

// Run 在获取分布式锁后执行任务。
func (w *DistributedLockJobWrapper) Run() {
	lockName := "cronjob_lock:" + w.LockKey
	// 只尝试获取一次锁，如果失败则不重试，这对于cron job是合适的行为
	mutex := w.Redsync.NewMutex(lockName, redsync.WithExpiry(w.LockExpiry), redsync.WithTries(1))

	// w.Logger.Debug("msg", "Attempting to acquire distributed lock", "key", lockName)

	if err := mutex.Lock(); err != nil {
		if errors.Is(err, redsync.ErrFailed) {
			// 这是预期的行为，当另一个实例持有锁时
			w.Logger.Info("Could not acquire distributed lock, job skipped", "key", lockName, "reason", "lock held by another instance")
			return
		}
		// 获取锁时发生其他错误
		w.Logger.Error("Error acquiring distributed lock, job skipped", "key", lockName, "error", err)
		return
	}

	// w.Logger.Info("Successfully acquired distributed lock, running job", "key", lockName)
	defer func() {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			w.Logger.Error("Error releasing distributed lock", "key", lockName, "error", err, "ok", ok)
		} else {
			// w.Logger.Info("Successfully released distributed lock", "key", lockName)
		}
	}()

	w.UnderlyingJob.Run()
}
