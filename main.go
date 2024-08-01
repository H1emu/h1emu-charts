package main

import (
	"net/http"
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
	// http.HandleFunc("/static", handler)
	fs := http.FileServer(http.Dir("public"))

	// Handle all requests by serving files from the "public" directory
	http.Handle("/", fs)
	http.ListenAndServe(":8080", nil)
}
