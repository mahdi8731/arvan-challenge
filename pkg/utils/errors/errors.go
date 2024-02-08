package errors

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	sentry.CaptureException(err)

	switch t := err.(type) {
	case *NotFoundError:
		return c.Status(t.code).JSON(map[string]any{
			"message": t.s,
		})
	case *ForbiddenError:
		return c.Status(t.code).JSON(map[string]any{
			"message": t.s,
		})
	case *InternalServerError:
		return c.Status(t.code).JSON(map[string]any{
			"message": t.s,
		})
	case *BadRequestError:
		return c.Status(t.code).JSON(map[string]any{
			"message": t.s,
		})
	}

	// Default 500 statuscode
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return c.Status(code).SendString(err.Error())
}

// Not found error with 404 status code and custom message
type NotFoundError struct {
	s    any
	code int
}

// Forbidden error with 403 status code and custom message
type ForbiddenError struct {
	s    any
	code int
}

// Internal server error with 500 status code and custom message
type InternalServerError struct {
	s    any
	code int
}

// Bad request error with 400 status code and custom message
type BadRequestError struct {
	s    any
	code int
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%v", n.s)
}

func (f *ForbiddenError) Error() string {
	return fmt.Sprintf("%v", f.s)
}

func (i *InternalServerError) Error() string {
	return fmt.Sprintf("%v", i.s)
}

func (b *BadRequestError) Error() string {
	return fmt.Sprintf("%v", b.s)
}

func NewNotFoundError(s any) error {
	return &NotFoundError{
		s:    s,
		code: 404,
	}
}

func NewForbiddenError(s any) error {
	return &ForbiddenError{
		s:    s,
		code: 403,
	}
}

func NewInternalServerError(s any) error {
	return &InternalServerError{
		s:    s,
		code: 500,
	}
}

func NewBadRequestError(s any) error {
	return &BadRequestError{
		s:    s,
		code: 400,
	}
}
