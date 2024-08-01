package main

import "fmt"

type ConnectionData struct {
	name string
	data []ConnectionsPerMonth
}

func main() {
	mongoCtx, mongoCancel := getMongoCtx()
	defer mongoCancel()
	db := getDb(mongoCtx)
	top := getTopKiller(db, mongoCtx, 12, "player")
	fmt.Println("Top 10 Killers:")
	for _, killer := range top {
		fmt.Printf("Player: %s, Count: %d\n", killer.CharacterName, killer.Count)
	}
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
