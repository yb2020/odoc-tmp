package errors

import (
	"fmt"

	"github.com/yb2020/odoc/pkg/response"
)

// TODO: 这个error里面应该要加一个类似fmt.errorf格式化相关的

// BizError 业务错误
type BizError struct {
	Status int32  // 状态码
	MsgID  string // 消息ID
	Err    error  // 原始错误
}

func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.MsgID, e.Err)
	}
	return e.MsgID
}

func (e *BizError) Unwrap() error {
	return e.Err
}

// Biz 创建普通业务错误，代表返回的status为0
func Biz(msgID string) error {
	return &BizError{Status: response.Status_Fail, MsgID: msgID}
}

// BizWithStatus 有特色返回码的业务错误
func BizWithStatus(status int32, msgID string) error {
	return &BizError{Status: status, MsgID: msgID}
}

// BizWrap 包装错误，代表返回的status为0
func BizWrap(msgID string, err error) error {
	return &BizError{Status: response.Status_Fail, MsgID: msgID, Err: err}
}

// BizWrapWithStatus 包装错误,有特色返回码的业务错误
func BizWrapWithStatus(status int32, msgID string, err error) error {
	return &BizError{Status: status, MsgID: msgID, Err: err}
}
