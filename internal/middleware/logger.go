package middleware

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// Custom response writer to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (writer bodyLogWriter) Write(b []byte) (int, error) {
	writer.body.Write(b)
	return writer.ResponseWriter.Write(b)
}

func ResponseLogger(ctx *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
	ctx.Writer = blw

	// Process request
	ctx.Next()

	responseBody := blw.body.String()
	println(responseBody)
}
