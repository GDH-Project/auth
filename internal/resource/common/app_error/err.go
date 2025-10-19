package apperror

import "fmt"

type Code string

const (
	UserNotFound         Code = "USER_NOT_FOUND"
	UserAlreadyUsed      Code = "USER_ALREADY_USED"
	UserCreateError      Code = "USER_CREATE_ERROR"
	UserPasswordNotMatch Code = "USER_PASSWORD_NOT_MATCH"
	UserUpdateFailed     Code = "USER_UPDATE_FAILED"
	UserRoleBindFailed   Code = "USER_BINDING_ROLE_FAILED"
	InternalServerError  Code = "INTERNAL_SERVER_ERROR"
	ParameterNotMatch    Code = "PARAMETER_NOT_MATCH"
	BadRequest           Code = "BAD_REQUEST"
)

type Error struct {
	Code        Code   `json:"code"`        // 오류에 대한 고유 식별자 (ex: "USER_NOT_FOUND")
	UserMessage string `json:"message"`     // 요청 응답 메시지
	DevMessage  string `json:"dev_message"` // 디버그 메시지
	StatusCode  int    `json:"-"`           // http Status 코드
	Cause       error  `json:"-"`           // 원본 오류
}

func New(code Code, userMessage string, devMessage string, statusCode int, cause error) *Error {
	return &Error{
		Code:        code,
		UserMessage: userMessage,
		DevMessage:  devMessage,
		StatusCode:  statusCode,
		Cause:       cause,
	}
}

func (e *Error) Error() string {
	if e.DevMessage != "" {
		return fmt.Sprintf("[%s] %s (Dev: %s, err: %v)", e.Code, e.UserMessage, e.DevMessage, e.Cause)
	} else {
		return fmt.Sprintf("[%s] %s (err: %v)", e.Code, e.UserMessage, e.Cause)
	}
}

func (e *Error) Unwrap() error {
	return e.Cause
}
