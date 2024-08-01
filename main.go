package main

type ConnectionData struct {
	name string
	data []ConnectionsPerMonth
}

func main() {
	mongoCtx, mongoCancel := getMongoCtx()
	defer mongoCancel()
	db := getDb(mongoCtx)
	count := getCountPerServer(db, mongoCtx, 11, CONNECTIONS_COLLECTION_NAME)
	println(count)
	// servers := getServers(db, mongoCtx)
	// connectionsDatas := make([]ConnectionData, len(servers))
	// allConnections := getAllConnections(db, mongoCtx)
	// for _, v := range servers {
	// 	println("doing serverid :", v.ServerId)
	// 	data := getConnectionsToServer(db, mongoCtx, v.ServerId)
	// 	println(len(data))
	// 	if len(data) == 0 {
	// 		continue
	// 	}
	// 	connectionsDatas = append(connectionsDatas, ConnectionData{name: v.Name, data: data})
	// }
	// createConnectionsChart(allConnections, connectionsDatas)
}
