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

// MembershipExpiredJob 会员过期任务
type MembershipExpiredJob struct {
	logger            logging.Logger
	spec              string                 // 任务的cron表达式，6字段标准cron表达式
	key               string                 // 任务的锁key，必须是唯一的unique-job-key
	expiry            time.Duration          // 任务的锁过期时间
	lockOpts          *scheduler.LockOptions // 任务的锁选项
	membershipService interfaces.IMembershipService
	userService       *userService.UserService
}

func NewMembershipExpiredJob(logger logging.Logger, cfg *config.Config, membershipService interfaces.IMembershipService, userService *userService.UserService) *MembershipExpiredJob {
	spec := cfg.Scheduler.Jobs.MembershipExpiredJob.Spec
	key := cfg.Scheduler.Jobs.MembershipExpiredJob.Key
	expiry := time.Duration(cfg.Scheduler.Jobs.MembershipExpiredJob.Expiry) * time.Second
	lockOpts := &scheduler.LockOptions{
		Key:    key,
		Expiry: expiry,
	}
	return &MembershipExpiredJob{logger: logger, spec: spec, key: key, expiry: expiry, lockOpts: lockOpts, membershipService: membershipService, userService: userService}
}

// NewUserContext 创建用户上下文
func (j *MembershipExpiredJob) NewUserContext(ctx context.Context, userId string) context.Context {
	user, err := j.userService.GetUserByID(ctx, userId)
	if err != nil {
		j.logger.Error("msg", "Get user by id failed", "userId", userId, "error", err)
		return ctx
	}
	j.logger.Info("msg", "Get user by id success", "userId", userId, "user", user)
	ctx = userContextUtil.SetUserContext(ctx, user)
	return ctx
}

// Spec 获取任务的cron表达式
func (j *MembershipExpiredJob) Spec() string {
	return j.spec
}

// LockOpts 获取任务的锁选项
func (j *MembershipExpiredJob) LockOpts() *scheduler.LockOptions {
	return j.lockOpts
}

// Run 执行任务，在执行任务前会获取锁，执行任务后会释放锁
// 这是必须实现的接口，否则实现ScheduledJob接口的调度器无法正确调度任务
func (j *MembershipExpiredJob) Run() {
	//j.logger.Info("msg", "Running membership expired job")
	ctx := context.Background()
	list, err := j.membershipService.GetFreeAccountExpiredList(ctx, 10) //只扫描Free会员
	if err != nil {
		j.logger.Error("msg", "Get account expired list failed", "error", err)
		return
	}
	// j.logger.Info("msg", "Get account expired list success", "list", list)
	for _, account := range list {
		j.logger.Info("msg", "Expired account", "account", account)
		ctx = j.NewUserContext(ctx, account.UserId)
		err := j.membershipService.HandleAccountExpired(ctx, account.Id, account.UserId)
		if err != nil {
			j.logger.Error("msg", "Handle account expired failed", "error", err)
			continue
		}
		j.logger.Info("msg", "Handle account expired success", "accountId", account.Id)
	}
}
