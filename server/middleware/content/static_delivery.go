package content

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo"
)

// Static delivers the static content taking into account the `/api` endpoints
// as well as the `/` and the `indec.html`. This middleware also sends the
// Content-Type needed for the browser to render the data acodingly
func Static(fs embed.FS) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			uri := c.Request().RequestURI

			if strings.HasPrefix(uri, "/api") {
				return next(c)
			}

			if uri == "/" {
				uri = "index.html"
			}

			uri = filepath.Join("assets", uri)

			data, err := fs.ReadFile(uri)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: reading file from embeded FS: %s\n", err.Error())
				return next(c)
			}

			mime := mimetype.Detect(data).String()
			if strings.HasSuffix(uri, "wasm") {
				mime = "application/wasm"
			}
			if strings.HasSuffix(uri, "html") {
				mime = "text/html; charset=utf-8"
			}
			if strings.HasSuffix(uri, "js") {
				mime = "application/js; charset=utf-8"
			}
			if strings.HasSuffix(uri, "png") {
				mime = "image/png"
			}

			return c.Blob(http.StatusOK, mime, data)
		}
	}
}
