package main

import (
	"fmt"
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
		start := time.Now()
		go genHtml()
		genCharts(db, mongoCtx)
		elapsed := time.Since(start)
		fmt.Printf("Elapsed time: %s\n", elapsed)
		time.Sleep(REFRESH_TIME * time.Second)
	}
}
