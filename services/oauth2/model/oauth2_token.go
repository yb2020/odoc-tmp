package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yb2020/odoc/pkg/model"
	pb "github.com/yb2020/odoc/proto/gen/go/oauth2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// OAuth2Token 令牌信息
type OAuth2Token struct {
	model.BaseModel                   // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	UserId          string            `json:"user_id" gorm:"column:user_id;index"`                 // 用户ID
	TokenId         string            `json:"token_id" gorm:"column:token_id;uniqueIndex"`         // 令牌ID
	AccessToken     string            `json:"access_token" gorm:"column:access_token;type:text"`   // 访问令牌
	RefreshToken    string            `json:"refresh_token" gorm:"column:refresh_token;type:text"` // 刷新令牌
	Roles           model.StringSlice `json:"roles" gorm:"column:roles;type:json"`                 // 角色
	Device          string            `json:"device" gorm:"column:device;type:varchar(255)"`       // 设备信息
	ExpiresAt       time.Time         `json:"expires_at"`                                          // 过期时间
	Revoked         bool              `json:"revoked" gorm:"column:revoked;default:false"`         // 是否已撤销
}

// TableName 指定表名
func (OAuth2Token) TableName() string {
	return "t_oauth2_token"
}

// TokenRequest 生成令牌请求
type TokenRequest struct {
	Username string // 用户名
	Password string // 密码
	Device   string // 设备信息
}

// GetUsername 获取用户名
func (x *TokenRequest) GetUsername() string {
	return x.Username
}

// GetPassword 获取密码
func (x *TokenRequest) GetPassword() string {
	return x.Password
}

// GetDevice 获取设备信息
func (x *TokenRequest) GetDevice() string {
	return x.Device
}

// RevokeRequest 撤销令牌请求
type RevokeRequest struct {
	AccessToken string // 访问令牌
}

// GetAccessToken 获取访问令牌
func (x *RevokeRequest) GetAccessToken() string {
	return x.AccessToken
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string // 访问令牌
	RefreshToken string // 刷新令牌
	TokenType    string // 令牌类型
	ExpiresAt    uint64 // 过期时间
	ExpiresIn    uint64 // 过期时间（秒）
	UserId       string // 用户ID
	Scope        string // 作用域
}

// SetExpiresAt 设置过期时间
func (x *TokenResponse) SetExpiresAt(t time.Time) {
	x.ExpiresAt = uint64(t.Unix())
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	ServiceId   string `json:"service_id"`   // 服务ID
	ServiceName string `json:"service_name"` // 服务名称
}

// Claims JWT声明
type Claims struct {
	UserId      string       `json:"user_id"`
	Username    string       `json:"username"`
	Roles       []string     `json:"roles"`
	Device      string       `json:"device"`
	ServiceInfo *ServiceInfo `json:"service_info,omitempty"` // 服务信息，仅服务令牌使用
	jwt.StandardClaims
}

// GetUserID 获取用户ID
func (c *Claims) GetUserID() string {
	return c.UserId
}

// GetUsername 获取用户名
func (c *Claims) GetUsername() string {
	return c.Username
}

// GetRoles 获取用户角色
func (c *Claims) GetRoles() []string {
	return c.Roles
}

// GetDevice 获取设备信息
func (c *Claims) GetDevice() string {
	return c.Device
}

// GetServiceId 获取服务ID
func (c *Claims) GetServiceId() string {
	return c.ServiceInfo.ServiceId
}

// GetServiceName 获取服务名称
func (c *Claims) GetServiceName() string {
	return c.ServiceInfo.ServiceName
}

// GetAuthCodeRequest 获取授权码请求
type GetAuthCodeRequest struct {
	pb.GetAuthCodeRequest
}

// ProtoReflect 实现 proto.Message 接口
func (x *GetAuthCodeRequest) ProtoReflect() protoreflect.Message {
	return x.GetAuthCodeRequest.ProtoReflect()
}

// GetAuthCodeResponse 获取授权码响应
type GetAuthCodeResponse struct {
	pb.GetAuthCodeResponse
}

// ProtoReflect 实现 proto.Message 接口
func (x *GetAuthCodeResponse) ProtoReflect() protoreflect.Message {
	return x.GetAuthCodeResponse.ProtoReflect()
}
