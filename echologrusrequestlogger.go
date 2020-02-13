package echologrusrequestlogger

import (
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"net"
	"net/http"
	"strconv"
	"time"
)

// func LogrusLogger(l *logrus.Logger) echo.MiddlewareFunc
// RequestLoggerを
func LogrusRequestLogger(l *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()

			// entry fields
			ef := l.WithFields(logrus.Fields{
				"host":      req.Host(),
				"uri":       req.URI(),
				"method":    req.Method(),
				"remote_ip": remoteAddr(req),
				"referer":   req.Referer(),
				"ua":        req.UserAgent(),
			})

			if rid := c.Request().Header().Get("X-Request-Id"); rid != "" {
				ef = ef.WithField("x-request-id", rid)
			}

			if err := next(c); err != nil {
				c.Error(err)
			}
			res := c.Response()

			latency := time.Since(start)

			ef.WithFields(logrus.Fields{
				"status_code": c.Response().Status(),
				"text_status": http.StatusText(c.Response().Status()),
				"took":        latency,
				"bytes":       strconv.FormatInt(res.Size(), 10),
			}).Info("completed handling request")

			return nil
		}
	}
}

// リモートアドレスを判定
func remoteAddr(req engine.Request) string {
	ra := req.RemoteAddress()
	if ip := req.Header().Get(echo.HeaderXRealIP); ip != "" {
		ra = ip
	} else if ip = req.Header().Get(echo.HeaderXForwardedFor); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}
