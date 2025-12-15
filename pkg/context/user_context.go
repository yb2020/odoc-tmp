package context

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
)

// 定义上下文键
type contextKey string

const (
	// UserIDKey 用户ID的上下文键
	UserIDKey contextKey = "user_id"
	// AccessTokenKey 访问令牌的上下文键
	AccessTokenKey contextKey = "access_token"
	// RolesKey 用户角色的上下文键
	RolesKey contextKey = "roles"
	// DeviceKey 设备信息的上下文键
	DeviceKey contextKey = "device"
	// UsernameKey 用户名的上下文键
	UsernameKey contextKey = "username"
	// AuthenticatedKey 认证状态的上下文键
	AuthenticatedKey contextKey = "authenticated"
	// ClaimsKey JWT声明的上下文键
	ClaimsKey contextKey = "claims"
	// UserContextKey 用户上下文对象的上下文键
	UserContextKey contextKey = "userContext"
	// ServiceIdKey 服务ID的上下文键
	ServiceIdKey contextKey = "service_id"
	// ServiceNameKey 服务名称的上下文键
	ServiceNameKey contextKey = "service_name"
)

// UserContext 用户上下文，包含用户相关的信息
type UserContext struct {
	UserId        string
	AccessToken   string
	Roles         []string
	Device        string
	Username      string
	Authenticated bool
	Claims        any
	ServiceId     string
	ServiceName   string
	mu            sync.RWMutex
	extraData     map[string]interface{}
}

// NewUserContext 创建新的用户上下文
func NewUserContext() *UserContext {
	return &UserContext{
		extraData: make(map[string]interface{}),
	}
}

// 链式设置方法
func (uc *UserContext) SetUserID(userID string) *UserContext      { uc.UserId = userID; return uc }
func (uc *UserContext) SetRoles(roles []string) *UserContext      { uc.Roles = roles; return uc }
func (uc *UserContext) SetDevice(device string) *UserContext      { uc.Device = device; return uc }
func (uc *UserContext) SetUsername(username string) *UserContext  { uc.Username = username; return uc }
func (uc *UserContext) SetAccessToken(token string) *UserContext  { uc.AccessToken = token; return uc }
func (uc *UserContext) SetAuthenticated(auth bool) *UserContext   { uc.Authenticated = auth; return uc }
func (uc *UserContext) SetClaims(claims interface{}) *UserContext { uc.Claims = claims; return uc }
func (uc *UserContext) SetServiceId(serviceId string) *UserContext {
	uc.ServiceId = serviceId
	return uc
}
func (uc *UserContext) SetServiceName(serviceName string) *UserContext {
	uc.ServiceName = serviceName
	return uc
}

// 额外数据操作
func (uc *UserContext) Set(key string, value interface{}) {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.extraData[key] = value
}

func (uc *UserContext) Get(key string) (interface{}, bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	val, ok := uc.extraData[key]
	return val, ok
}

// ToContext 将用户上下文转换为标准库的context.Context
func (uc *UserContext) ToContext(ctx context.Context) context.Context {
	addToContext := func(key contextKey, value interface{}) {
		if value != nil {
			ctx = context.WithValue(ctx, key, value)
		}
	}

	addToContext(UserIDKey, uc.UserId)
	addToContext(RolesKey, uc.Roles)
	addToContext(DeviceKey, uc.Device)
	addToContext(UsernameKey, uc.Username)
	addToContext(AuthenticatedKey, uc.Authenticated)
	addToContext(AccessTokenKey, uc.AccessToken)
	addToContext(ClaimsKey, uc.Claims)
	addToContext(ServiceIdKey, uc.ServiceId)
	addToContext(ServiceNameKey, uc.ServiceName)

	return ctx
}

// Clone 创建用户上下文的副本
func (uc *UserContext) Clone() *UserContext {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	newUC := &UserContext{
		UserId:        uc.UserId,
		AccessToken:   uc.AccessToken,
		Roles:         make([]string, len(uc.Roles)),
		Device:        uc.Device,
		Username:      uc.Username,
		Authenticated: uc.Authenticated,
		Claims:        uc.Claims,
		ServiceId:     uc.ServiceId,
		ServiceName:   uc.ServiceName,
		extraData:     make(map[string]interface{}),
	}

	copy(newUC.Roles, uc.Roles)
	for k, v := range uc.extraData {
		newUC.extraData[k] = v
	}

	return newUC
}

// FromGinContext 从Gin上下文中提取用户上下文
func FromGinContext(c *gin.Context) *UserContext {
	uc := NewUserContext()

	getValue := func(key contextKey, setter func(interface{})) {
		if value, exists := c.Get(string(key)); exists && value != nil {
			setter(value)
		}
	}

	getValue(UserIDKey, func(v interface{}) {
		if id, ok := v.(string); ok {
			uc.UserId = id
		}
	})

	getValue(RolesKey, func(v interface{}) {
		if roles, ok := v.([]string); ok {
			uc.Roles = roles
		}
	})

	getValue(AccessTokenKey, func(v interface{}) {
		if token, ok := v.(string); ok {
			uc.AccessToken = token
		}
	})

	getValue(DeviceKey, func(v interface{}) {
		if device, ok := v.(string); ok {
			uc.Device = device
		}
	})

	getValue(UsernameKey, func(v interface{}) {
		if username, ok := v.(string); ok {
			uc.Username = username
		}
	})

	getValue(AuthenticatedKey, func(v interface{}) {
		if auth, ok := v.(bool); ok {
			uc.Authenticated = auth
		}
	})

	getValue(ClaimsKey, func(v interface{}) {
		uc.Claims = v
	})

	getValue(ServiceIdKey, func(v interface{}) {
		if serviceId, ok := v.(string); ok {
			uc.ServiceId = serviceId
		}
	})

	getValue(ServiceNameKey, func(v interface{}) {
		if serviceName, ok := v.(string); ok {
			uc.ServiceName = serviceName
		}
	})

	return uc
}

// ToGinContext 将用户上下文存储到Gin上下文中
func (uc *UserContext) ToGinContext(c *gin.Context) {
	c.Set(string(UserIDKey), uc.UserId)
	c.Set(string(RolesKey), uc.Roles)
	c.Set(string(DeviceKey), uc.Device)
	c.Set(string(UsernameKey), uc.Username)
	c.Set(string(AuthenticatedKey), uc.Authenticated)
	c.Set(string(AccessTokenKey), uc.AccessToken)
	c.Set(string(ClaimsKey), uc.Claims)
	c.Set(string(UserContextKey), uc)
	c.Set(string(ServiceIdKey), uc.ServiceId)
	c.Set(string(ServiceNameKey), uc.ServiceName)

	c.Request = c.Request.WithContext(uc.ToContext(c.Request.Context()))
}

// GetUserContext 从标准上下文中获取用户信息
func GetUserContext(ctx context.Context) *UserContext {
	uc := NewUserContext()

	getValue := func(key contextKey, setter func(interface{})) {
		if value := ctx.Value(key); value != nil {
			setter(value)
		}
	}

	getValue(UserIDKey, func(v interface{}) {
		if id, ok := v.(string); ok {
			uc.UserId = id
		}
	})

	getValue(RolesKey, func(v interface{}) {
		if roles, ok := v.([]string); ok {
			uc.Roles = roles
		}
	})

	getValue(DeviceKey, func(v interface{}) {
		if device, ok := v.(string); ok {
			uc.Device = device
		}
	})

	getValue(UsernameKey, func(v interface{}) {
		if username, ok := v.(string); ok {
			uc.Username = username
		}
	})

	getValue(AccessTokenKey, func(v interface{}) {
		if token, ok := v.(string); ok {
			uc.AccessToken = token
		}
	})

	getValue(AuthenticatedKey, func(v interface{}) {
		if auth, ok := v.(bool); ok {
			uc.Authenticated = auth
		}
	})

	getValue(ServiceIdKey, func(v interface{}) {
		if id, ok := v.(string); ok {
			uc.ServiceId = id
		}
	})

	getValue(ServiceNameKey, func(v interface{}) {
		if name, ok := v.(string); ok {
			uc.ServiceName = name
		}
	})

	getValue(ClaimsKey, func(v interface{}) {
		uc.Claims = v
	})

	return uc
}

// GetServiceId 从标准上下文中获取服务ID
func GetServiceId(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ServiceIdKey).(string)
	return id, ok
}

// GetServiceName 从标准上下文中获取服务名称
func GetServiceName(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(ServiceNameKey).(string)
	return name, ok
}

// 辅助函数
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

func GetRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(RolesKey).([]string)
	return roles, ok
}

func HasRole(ctx context.Context, role string) bool {
	roles, ok := GetRoles(ctx)
	if !ok {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func IsAuthenticated(ctx context.Context) bool {
	authenticated, ok := ctx.Value(AuthenticatedKey).(bool)
	return ok && authenticated
}

// FromGinToStdContext 将Gin上下文转换为标准库的context.Context
func FromGinToStdContext(c *gin.Context) context.Context {
	return FromGinContext(c).ToContext(c.Request.Context())
}
