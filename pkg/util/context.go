package util

import (
	"context"

	"rapid-bridge/constants"

	"github.com/labstack/echo/v4"
)

func GetReqCtxFromEchoCtx(c echo.Context) context.Context {
	ctx := c.Request().Context()

	req := c.Request()

	requestID := req.Header.Get("X-Request-ID")
	from := req.Header.Get(constants.From)
	to := req.Header.Get(constants.To)
	keyVersion := req.Header.Get(constants.KeyVersion)

	ctx = context.WithValue(ctx, constants.RequestId, requestID)
	ctx = context.WithValue(ctx, constants.From, from)
	ctx = context.WithValue(ctx, constants.To, to)
	ctx = context.WithValue(ctx, constants.KeyVersion, keyVersion)

	return ctx
}
