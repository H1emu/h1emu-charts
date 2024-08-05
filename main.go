package main

import (
	"os"
)

type ConnectionData struct {
	name string
	data []ConnectionsPerMonth
}

func main() {
	mongoCtx, mongoCancel := getMongoCtx()
	defer mongoCancel()
	db := getDb(mongoCtx)
	os.Mkdir("public", 0755)
	genCharts(db, mongoCtx)
	genHtml()
	serveHtml()
}
