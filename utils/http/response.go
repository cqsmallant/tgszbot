package http

import (
	"ant/utils/constant"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Resp struct{}

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	RequestId  string      `json:"request_id"`
}

func (r *Resp) View(e echo.Context, code int, html string) error {
	return e.HTML(code, html)
}

// SucView 成功返回
func (r *Resp) SucView(e echo.Context, html string) error {
	return r.View(e, http.StatusOK, html)
}

func (r *Resp) Json(e echo.Context, code int, data interface{}) error {
	return e.JSON(code, data)
}

// SucJson 成功返回
func (r *Resp) SucJson(e echo.Context, data interface{}, message ...string) error {
	rp := new(Response)
	rp.StatusCode = http.StatusOK
	if len(message) == 0 {
		rp.Message = "Success"
	} else {
		for _, m := range message {
			rp.Message += "," + m
		}
	}
	rp.Data = data
	rp.RequestId = e.Request().Header.Get(echo.HeaderXRequestID)
	return r.Json(e, http.StatusOK, rp)
}

func (r *Resp) ErrJson(e echo.Context, err error) error {
	rr := new(Response)
	switch err.(type) {
	case *constant.RspError:
		rr.StatusCode, rr.Message = err.(*constant.RspError).Render()
	default:
		rr.StatusCode = 400
		rr.Message = err.Error()
	}
	rr.RequestId = e.Request().Header.Get(echo.HeaderXRequestID)
	return r.Json(e, http.StatusOK, &rr)
}
