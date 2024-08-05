package main

import (
	"os"
	"time"
)

type ConnectionData struct {
	name string
	data []ConnectionsPerMonth
}

const REFRESH_TIME = 600

func main() {
	mongoCtx, mongoCancel := getMongoCtx()
	defer mongoCancel()
	db := getDb(mongoCtx)
	os.Mkdir("public", 0755)
	go serveHtml()
	for {
		genCharts(db, mongoCtx)
		genHtml()
		time.Sleep(REFRESH_TIME * time.Second)
	}
}
