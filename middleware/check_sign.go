package middleware

import (
	"ant/utils/constant"
	"bytes"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

func CheckApiSign() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			params, err := ioutil.ReadAll(ctx.Request().Body)
			if err != nil {
				return constant.SignatureErr
			}

			ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(params))
			return next(ctx)
		}

	}
}
