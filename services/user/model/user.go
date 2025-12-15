package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	pb "github.com/yb2020/odoc-proto/gen/go/user"
	"github.com/yb2020/odoc/pkg/model"
	"github.com/yb2020/odoc/pkg/utils"
)

// UserRoles 用户角色列表类型
type UserRoles []pb.UserRole

// Value 实现driver.Valuer接口，用于将UserRoles保存到数据库
func (r UserRoles) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan 实现sql.Scanner接口，用于从数据库读取UserRoles
func (r *UserRoles) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// User 表示用户实体
type User struct {
	model.BaseModel                   // 嵌入基础模型，继承ID、CreatedAt、UpdatedAt字段和钩子方法
	Username            string        `json:"username" gorm:"uniqueIndex:idx_unique_t_user_username;not null"`
	Email               string        `json:"email" gorm:"uniqueIndex:idx_unique_t_user_email;not null"`
	Password            string        `json:"-" gorm:"not null"` // 不在 JSON 中暴露密码
	Nickname            string        `json:"nickname" gorm:"type:varchar(200)"`
	Avatar              string        `json:"avatar" gorm:"type:varchar(200)"`
	Roles               UserRoleSlice `json:"roles" gorm:"type:json"`
	GoogleOpenId        string        `json:"google_open_id" gorm:"type:varchar(200);uniqueIndex:idx_unique_t_user_google_open_id"`
	Status              pb.UserStatus `json:"status" gorm:"type:int;default:1"` // 默认为激活状态
	AccessToken         string        `json:"access_token" gorm:"-"`
	AccessTokenExpires  uint64        `json:"access_token_expires" gorm:"-"`
	RefreshToken        string        `json:"refresh_token" gorm:"-"`
	RefreshTokenExpires uint64        `json:"refresh_token_expires" gorm:"-"`
	LastLogin           uint64        `json:"last_login" gorm:"-"`
	UserAgent           string        `json:"user_agent" gorm:"-"`
	IPAddress           string        `json:"ip_address" gorm:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "t_user"
}

type UserRoleSlice []pb.UserRole

// ToStringSlice converts a slice of UserRole to a slice of strings.
func (s UserRoleSlice) ToStringSlice() []string {
	result := make([]string, len(s))
	for i, r := range s {
		result[i] = r.String()
	}
	return result
}

// Value 实现 driver.Valuer 接口
func (r UserRoleSlice) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan 实现 sql.Scanner 接口
func (r *UserRoleSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// ToProto 将User模型转换为protobuf消息
// TODO: proto 文件需要更新为 string 类型的 ID
func (u *User) ToProto() *pb.User {
	if u == nil {
		return nil
	}
	u.Email = utils.Around(u.Email, 2, strings.Index(u.Email, "@"))
	u.Username = utils.Left(u.Username, 4)
	pbUser := &pb.User{
		Id:                  u.Id,
		Username:            u.Username,
		Email:               u.Email,
		Nickname:            u.Nickname,
		Avatar:              u.Avatar,
		Roles:               []pb.UserRole(u.Roles),
		AccessToken:         u.AccessToken,
		AccessTokenExpires:  u.AccessTokenExpires,
		RefreshToken:        u.RefreshToken,
		RefreshTokenExpires: u.RefreshTokenExpires,
		LastLogin:           u.LastLogin,
		GoogleOpenId:        u.GoogleOpenId,
		Status:              u.Status,
	}

	return pbUser
}

// FromProto 从protobuf消息更新User模型
// TODO: proto 文件需要更新为 string 类型的 ID
func (u *User) FromProto(pbUser *pb.User) {
	if pbUser == nil {
		return
	}

	// 将uint64类型的Id转换为int64
	u.Id = pbUser.Id
	u.Username = pbUser.Username
	u.Email = pbUser.Email
	u.Password = pbUser.Password
	u.Nickname = pbUser.Nickname
	u.Avatar = pbUser.Avatar
	u.Roles = UserRoleSlice(pbUser.Roles)
	u.AccessToken = pbUser.AccessToken
	u.AccessTokenExpires = pbUser.AccessTokenExpires
	u.RefreshToken = pbUser.RefreshToken
	u.RefreshTokenExpires = pbUser.RefreshTokenExpires
	u.LastLogin = pbUser.LastLogin
	u.GoogleOpenId = pbUser.GoogleOpenId
	if pbUser.Status != pb.UserStatus_STATUS_UNKNOWN {
		u.Status = pbUser.Status
	}
}
