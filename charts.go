package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-echarts/go-echarts/v2/opts"
)

// generate random data for bar chart
func generateLineItems(data []ConnectionsPerMonth) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i].Count})
	}
	return items
}

func getXaxis(ConnectionDatas []ConnectionsPerMonth) []string {
	xais := make([]string, 0)
	for _, v := range ConnectionDatas {
		xais = append(xais, v.Id)
	}
	return xais
}

func populateMissingConnectionsData(all []ConnectionsPerMonth, current []ConnectionsPerMonth) []ConnectionsPerMonth {
	result := make([]ConnectionsPerMonth, len(all))
	for i := 0; i < len(result); i++ {
		for ry := 0; ry < len(current); ry++ {
			if current[ry].Id == all[i].Id && current[ry].Count > 10000 {
				result[i] = current[ry]
			}
		}
	}

	return result
}

func createConnectionsChart(allConnections []ConnectionsPerMonth, connectionDatas []ConnectionData) {
	// create a new line instance
	line := charts.NewLine()

	// Put data into instance
	line.SetXAxis(getXaxis(allConnections))
	line.AddSeries("all", generateLineItems(allConnections))
	selected := make(map[string]bool)
	for _, cd := range connectionDatas {
		data := populateMissingConnectionsData(allConnections, cd.data)
		line.AddSeries(cd.name, generateLineItems(data))
		selected[cd.name] = false
	}
	line.SetGlobalOptions(charts.WithLegendOpts(opts.Legend{
		Selected: selected,
	}))
	f, _ := os.Create("public/connections.html")
	line.Render(f)
}

func createTop10KillerChart(chartName string, topKiller []TopKiller) {
	// create a new bar instance
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: chartName,
	}))

	items := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	for _, killer := range topKiller {
		xAxis = append(xAxis, killer.CharacterName)
		items = append(items, opts.BarData{Value: killer.Count})
	}
	bar.SetXAxis(xAxis)
	bar.AddSeries("Killer", items)
	// Where the magic happens
	f, error := os.Create("public/" + chartName + ".html")
	if error != nil {
		panic(error)
	}
	bar.Render(f)
}

func createCountPerServerCharts(db *mongo.Database, mongoCtx context.Context, serverList []Server, collectionName string, chartName string) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: chartName,
	}))
	items := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	for i := 0; i < len(serverList); i++ {
		v := serverList[i]
		xAxis = append(xAxis, fmt.Sprintf("%s (%s)", v.Name, v.Region))
		count := getCountPerServer(db, mongoCtx, v.ServerId, collectionName)
		items = append(items, opts.BarData{Value: count})

	}
	bar.SetXAxis(xAxis)
	bar.AddSeries(collectionName, items)
	f, error := os.Create("public/" + chartName + ".html")
	if error != nil {
		panic(error)
	}
	bar.Render(f)
}

func genCharts(db *mongo.Database, mongoCtx context.Context) {
	servers := getServers(db, mongoCtx)
	officialServers := []Server{}
	for _, v := range servers {
		if v.IsOfficial && !v.IsDisabled {
			officialServers = append(officialServers, v)
		}
	}
	createCountPerServerCharts(db, mongoCtx, officialServers, CONSTRUCTIONS_COLLECTION_NAME, "Constructions per server")
	createCountPerServerCharts(db, mongoCtx, officialServers, CROPS_COLLECTION_NAME, "Crops per server")
	top := getTopKiller(db, mongoCtx, 12, "player")
	createTop10KillerChart("Top Killer Main EU 1", top)
	top = getTopKiller(db, mongoCtx, 11, "player")
	createTop10KillerChart("Top Killer Main US 1", top)
	top = getTopKiller(db, mongoCtx, 61, "zombie")
	createTop10KillerChart("Top Zombie Killer Help", top)
	connectionsDatas := make([]ConnectionData, len(servers))
	allConnections := getAllConnections(db, mongoCtx)
	for _, v := range servers {
		println("doing serverid :", v.ServerId)
		data := getConnectionsToServer(db, mongoCtx, v.ServerId)
		println(len(data))
		if len(data) == 0 {
			continue
		}
		connectionsDatas = append(connectionsDatas, ConnectionData{name: v.Name, data: data})
	}
	createConnectionsChart(allConnections, connectionsDatas)
}
