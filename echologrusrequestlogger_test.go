package echologrusrequestlogger

import (
    "testing"
    "github.com/labstack/echo"
    "github.com/labstack/echo/test"
    "net/http"
    "errors"
    "github.com/Sirupsen/logrus"
)

func TestLogrusRequestLogger(t *testing.T) {
    logger := logrus.New()
    // Note: Just for the test coverage, not a real test.
    e := echo.New()
    req := test.NewRequest(echo.GET, "/", nil)
    rec := test.NewResponseRecorder()
    c := e.NewContext(req, rec)
    h := LogrusRequestLogger(logger)(func(c echo.Context) error {
        return c.String(http.StatusOK, "test")
    })

    // Status 2xx
    h(c)

    // Status 3xx
    rec = test.NewResponseRecorder()
    c = e.NewContext(req, rec)
    h = LogrusRequestLogger(logger)(func(c echo.Context) error {
        return c.String(http.StatusTemporaryRedirect, "test")
    })
    h(c)

    // Status 4xx
    rec = test.NewResponseRecorder()
    c = e.NewContext(req, rec)
    h = LogrusRequestLogger(logger)(func(c echo.Context) error {
        return c.String(http.StatusNotFound, "test")
    })
    h(c)

    // Status 5xx with empty path
    req = test.NewRequest(echo.GET, "", nil)
    rec = test.NewResponseRecorder()
    c = e.NewContext(req, rec)
    h = LogrusRequestLogger(logger)(func(c echo.Context) error {
        return errors.New("error")
    })
    h(c)
}
