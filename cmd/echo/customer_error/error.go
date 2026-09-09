package customer_error

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func CustomerErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return // already sent by a handler/middleware
		}
	}
	code := http.StatusInternalServerError
	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) {
		if tmp := sc.StatusCode(); tmp != 0 {
			code = tmp
		}
	}
	var cErr error
	if c.Request().Method == http.MethodHead {
		cErr = c.NoContent(code)
	} else {
		cErr = c.File(fmt.Sprintf("%d.html", code)) // e.g. 404.html, 500.html
	}
	if cErr != nil {
		c.Logger().Error("failed to send error page", "error", errors.Join(err, cErr))
	}
}
