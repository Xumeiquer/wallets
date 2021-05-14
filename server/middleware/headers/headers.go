package headers

import "github.com/labstack/echo"

type Headers map[string]string

func Set(H Headers) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			for k, v := range H {
				c.Response().Header().Set(k, v)
			}
			return next(c)
		}
	}
}
