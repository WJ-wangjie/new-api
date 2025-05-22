package common

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

// CustomResponseWriter 自定义的响应写入器
type CustomResponseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// Write 重写 Write 方法，捕获写入的数据
func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	// 将写入的数据追加到 body 缓冲区
	n, err := crw.Body.Write(b)
	if err != nil {
		return n, err
	}
	// 调用原始的 Write 方法将数据写入客户端
	return crw.ResponseWriter.Write(b)
}
