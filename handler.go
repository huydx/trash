package trash

import (
	"github.com/gin-gonic/gin"
)

func PostSpansHandler(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans := make([]Span, 0)
		if err := c.BindJSON(&spans); err != nil {
			c.Status(400)
			return
		}
		if err := storage.Index(spans); err != nil {
			c.Status(400)
			return
		}
	}
}

func GetTraceHandler(storage Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id = c.Param("traceId")
		spans, err := storage.GetTrace([]byte(id))
		if err != nil {
			c.Status(400)
			return
		}
		c.JSON(200, spans)
		return
	}
}
