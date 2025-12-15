package context

import (
	"context"

	userContext "github.com/yb2020/odoc/pkg/context"
	userModel "github.com/yb2020/odoc/services/user/model"
)

func SetUserContext(ctx context.Context, user *userModel.User) context.Context {
	if user == nil {
		return ctx
	}
	uc := userContext.NewUserContext()
	uc.SetUserID(user.Id)
	uc.SetUsername(user.Username)
	uc.SetAuthenticated(true)
	// 转换Roles为[]string
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.String()
	}
	uc.SetRoles(roles)
	uc.SetDevice(user.UserAgent)
	ctx = uc.ToContext(ctx)
	return ctx
}
