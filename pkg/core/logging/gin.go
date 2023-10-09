package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HttpLoggingMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	d := time.Since(start)
	logrus.WithContext(c.Request.Context()).
		WithField("duration", d.Nanoseconds()).
		WithField("http_host", c.Request.Host).
		WithField("http_path", c.Request.URL.Path).
		WithField("http_query_string", c.Request.URL.RawQuery).
		WithField("http_status_code", c.Writer.Status()).
		WithField("http_method", c.Request.Method).
		WithField("http_referer", c.Request.Header.Get("Referer")).
		WithField("http_request_id", c.Request.Header.Get("Request_ID")).
		WithField("http_useragent", c.Request.Header.Get("User-Agent")).
		WithField("http_version", c.Request.Proto).
		WithField("http_x_forwarded_for", c.Request.Header.Get("X-Forwarded-For")).
		Info("end handling HTTP request")
}
