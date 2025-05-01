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
	from := req.Header.Get("X-Source-Slug")
	to := req.Header.Get("X-Destination-Slug")
	rapidUrl := req.Header.Get("X-Rapid-Url")

	ctx = context.WithValue(ctx, constants.RequestId, requestID)
	ctx = context.WithValue(ctx, constants.From, from)
	ctx = context.WithValue(ctx, constants.To, to)
	ctx = context.WithValue(ctx, constants.RapidUrl, rapidUrl)

	return ctx
}
