package errors

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate protoc -I. --go_out=paths=source_relative:. errors.proto

type Error struct {
	Status
	cause error
}

func New(code int, message string) *Error {
	return &Error{
		Status: Status{Code: uint32(code), Message: message},
	}
}

func Newf(code int, format string, a ...interface{}) *Error {
	return New(code, fmt.Sprintf(format, a...))
}

func Errorf(code int, format string, a ...interface{}) error {
	return New(code, fmt.Sprintf(format, a...))
}

// 实现错误接口
func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s data = %v cause = %v", e.Code, e.Message, e.Metadata, e.cause)
}

// Unwrap 获取错误的根本原因
func (e *Error) Unwrap() error {
	return e.cause
}

// Is 判断是否是相同的错误
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code
	}
	return false
}

// WithCause 错误的根本原因
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

// WithMetadata 设置错误元信息
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := Clone(e)
	err.Metadata = md
	return err
}

// GetGrpcStatus 获取 grpc 的错误状态
func (e *Error) GetGrpcStatus() *status.Status {
	s, _ := status.New(codes.Code(e.Code), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Metadata: e.Metadata,
		})
	return s
}

// Clone 克隆一样的错误
func Clone(err *Error) *Error {
	// 如果错误为 nil
	if err == nil {
		return nil
	}

	// 复制 metadata
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}

	// 返回同样的错误
	return &Error{
		cause: err.cause,
		Status: Status{
			Code:     err.Code,
			Message:  err.Message,
			Metadata: metadata,
		},
	}
}

// FromError 将错误转为当前类型
func FromError(err error) *Error {
	// 如果错误为 nil
	if err == nil {
		return nil
	}

	// 如果是当前错误
	if se := new(Error); errors.As(err, &se) {
		return se
	}

	// 如果是 grpc 的错误码
	ge, ok := status.FromError(err)
	if !ok {
		return New(int(codes.Unknown), err.Error())
	}

	// 设置为当前类型的错误
	e := New(int(ge.Code()), ge.Message())

	// 获取 grpc 错误中的 details
	for _, detail := range ge.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			return e.WithMetadata(d.Metadata)
		}
	}

	return e
}
