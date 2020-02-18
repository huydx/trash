package trash

import (
	"github.com/gin-gonic/gin"
	"github.com/huydx/trash"
	"log"
)

func main() {
	g := gin.New()
	s := trash.BadgerStorage{}
	g.GET("/api/v2/trace/{traceId}", trash.GetTraceHandler(s))
	g.POST("/api/v2/spans", trash.PostSpansHandler(s))
	if err := g.Run(":9411"); err != nil {
		log.Fatal(err)
	}
}
