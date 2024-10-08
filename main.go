package main

import (
	"os"
	"time"
)

const REFRESH_TIME = 600

func main() {
	mongoCtx, mongoCancel := getMongoCtx()
	defer mongoCancel()
	db := getDb(mongoCtx)
	// rm everything at each reload
	os.RemoveAll("public")
	os.Mkdir("public", 0755)
	go serveHtml()
	for {
		genCharts(db, mongoCtx)
		genHtml()
		time.Sleep(REFRESH_TIME * time.Second)
	}
}
