package slicex

import "errors"

// 定义预定义错误，用于在特定条件下标识错误类型。
// 这样做可以提供更丰富的错误信息，便于调用者根据错误类型进行处理。
var (
	// LengthError 表示长度必须大于0的错误情况。
	// 当遇到长度不符合要求的情况时，可以返回此错误。
	LengthError = errors.New("length must be greater than 0")

	// FnIsNil 表示函数指针为空的错误情况。
	// 当函数参数预期为非空而实际为空时，可以返回此错误。
	FnIsNil = errors.New("fn is nil")
)
