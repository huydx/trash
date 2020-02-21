package trash

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/huydx/trash"

	"github.com/dgraph-io/badger"
)

func main() {
	g := gin.New()
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	s := trash.BadgerStorage{
		DB: db,
	}
	g.GET("/api/v2/trace/{traceId}", trash.GetTraceHandler(s))
	g.POST("/api/v2/spans", trash.PostSpansHandler(s))
	if err := g.Run(":9411"); err != nil {
		log.Fatal(err)
	}
}
