package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type MessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func newMessageResponse(ctx *fiber.Ctx, code int, message string) error {
	return ctx.Status(code).JSON(MessageResponse{
		Code:    code,
		Message: message,
	})
}

type ResponseError struct {
	StatusCode int
	Code       int
	Message    string
} // @name ResponseError

func (e *ResponseError) Wrap(err error) error {
	if err != nil {
		return &ResponseError{
			StatusCode: e.StatusCode,
			Code:       e.Code,
			Message:    fmt.Sprintf("%s: %s", e.Message, err.Error()),
		}
	}

	return e
}

func (e *ResponseError) WithErr(err error) *ResponseError {
	if err != nil {
		return &ResponseError{
			StatusCode: e.StatusCode,
			Code:       e.Code,
			Message:    fmt.Sprintf("%s: %s", e.Message, err.Error()),
		}
	}

	return e
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("http error: %d %d %s", e.StatusCode, e.Code, e.Message)
}
