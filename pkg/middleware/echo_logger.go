package middleware

import (
	"rapid-bridge/domain/port"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateEchoLogger(logger port.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogHost:         true,
		LogMethod:       true,
		LogRemoteIP:     true,
		LogResponseSize: true,
		LogLatency:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				// "request_id", c.Get(constants.RequestId),
				"method", v.Method,
				"host", v.Host,
				"uri", v.URI,
				"remote_ip", v.RemoteIP,
				"start_time", v.StartTime,
				"status", v.Status,
				"response_size", v.ResponseSize,
				"latency", v.Latency,
			)

			return nil
		},
	})
}
