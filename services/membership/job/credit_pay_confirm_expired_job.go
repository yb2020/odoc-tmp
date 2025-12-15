package job

import (
	"context"
	"time"

	"github.com/yb2020/odoc/config"
	userContextUtil "github.com/yb2020/odoc/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/scheduler"
	"github.com/yb2020/odoc/services/membership/interfaces"
	userService "github.com/yb2020/odoc/services/user/service"
)

// CreditPayConfirmExpiredJob 会员积分确认支付过期任务
type CreditPayConfirmExpiredJob struct {
	logger               logging.Logger
	spec                 string                 // 任务的cron表达式，6字段标准cron表达式
	key                  string                 // 任务的锁key，必须是唯一的unique-job-key
	expiry               time.Duration          // 任务的锁过期时间
	lockOpts             *scheduler.LockOptions // 任务的锁选项
	membershipService    interfaces.IMembershipService
	creditPaymentService interfaces.ICreditPaymentService
	userService          *userService.UserService
}

func NewCreditPayConfirmExpiredJob(logger logging.Logger, cfg *config.Config, membershipService interfaces.IMembershipService, creditPaymentService interfaces.ICreditPaymentService, userService *userService.UserService) *CreditPayConfirmExpiredJob {
	spec := cfg.Scheduler.Jobs.CreditPayConfirmExpiredJob.Spec
	key := cfg.Scheduler.Jobs.CreditPayConfirmExpiredJob.Key
	expiry := time.Duration(cfg.Scheduler.Jobs.CreditPayConfirmExpiredJob.Expiry) * time.Second
	lockOpts := &scheduler.LockOptions{
		Key:    key,
		Expiry: expiry,
	}
	return &CreditPayConfirmExpiredJob{logger: logger, spec: spec, key: key, expiry: expiry, lockOpts: lockOpts, membershipService: membershipService, creditPaymentService: creditPaymentService, userService: userService}
}

// Spec 获取任务的cron表达式
func (j *CreditPayConfirmExpiredJob) Spec() string {
	return j.spec
}

// LockOpts 获取任务的锁选项
func (j *CreditPayConfirmExpiredJob) LockOpts() *scheduler.LockOptions {
	return j.lockOpts
}

// NewUserContext 创建用户上下文
func (j *CreditPayConfirmExpiredJob) NewUserContext(ctx context.Context, userId string) context.Context {
	user, err := j.userService.GetUserByID(ctx, userId)
	if err != nil {
		j.logger.Error("msg", "Get user by id failed", "userId", userId, "error", err)
		return ctx
	}
	j.logger.Info("msg", "Get user by id success", "userId", userId, "user", user)
	ctx = userContextUtil.SetUserContext(ctx, user)
	return ctx
}

// Run 执行任务，在执行任务前会获取锁，执行任务后会释放锁
// 这是必须实现的接口，否则实现ScheduledJob接口的调度器无法正确调度任务
func (j *CreditPayConfirmExpiredJob) Run() {
	// j.logger.Info("msg", "Running credit pay confirm expired job")
	ctx := context.Background()
	list, err := j.creditPaymentService.GetConfirmExpiredList(ctx, 10)
	if err != nil {
		j.logger.Error("msg", "Get confirm expired list failed", "error", err)
		return
	}
	// j.logger.Info("msg", "Get confirm expired list success", "list", list)
	for _, record := range list {
		//1.获取包含用户信息的上下文
		ctx = j.NewUserContext(ctx, record.UserId)
		//2.确认积分服务成功
		err = j.membershipService.ConfirmCreditFun(ctx, record.Id)
		if err != nil {
			j.logger.Error("msg", "Confirm credit pay confirm expired job failed", "error", err)
			continue
		}
		j.logger.Info("msg", "Confirm credit pay confirm expired job success", "recordId", record.Id)
	}
}
