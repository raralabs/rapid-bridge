package middleware

import (
	"fmt"
	"net/http"
	"rapid-bridge/constants"
	"rapid-bridge/pkg/util"

	"github.com/labstack/echo/v4"
)

func APIContractMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqCtx := util.GetReqCtxFromEchoCtx(c)

			from, _ := reqCtx.Value(constants.From).(string)
			if from == "" {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s not found in header", constants.From))
			}
			to, _ := reqCtx.Value(constants.To).(string)
			if to == "" {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s not found in header", constants.To))
			}
			keyVersion, _ := reqCtx.Value(constants.KeyVersion).(string)
			if keyVersion == "" {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s not found in header", constants.KeyVersion))
			}

			err := next(c)

			return err
		}
	}
}
