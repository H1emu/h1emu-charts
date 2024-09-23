package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-echarts/go-echarts/v2/opts"
)

// generate random data for bar chart
func generateLineItems(data []CountPerTime) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i].Count})
	}
	return items
}

func getXaxis(ConnectionDatas []CountPerTime) []string {
	xais := make([]string, 0)
	for _, v := range ConnectionDatas {
		xais = append(xais, v.Id)
	}
	return xais
}

// Since servers were not all created a the same time
func populateMissingConnectionsData(all []CountPerTime, current []CountPerTime) []CountPerTime {
	result := make([]CountPerTime, len(all))
	for i := 0; i < len(result); i++ {
		for ry := 0; ry < len(current); ry++ {
			if current[ry].Id == all[i].Id {
				result[i] = current[ry]
			}
		}
	}

	return result
}

func createLineCountChart(name string, allConnections []CountPerTime, connectionDatas []CountData) {
	// create a new line instance
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: name,
	}))

	// Put data into instance
	line.SetXAxis(getXaxis(allConnections))
	line.AddSeries("All", generateLineItems(allConnections))
	selected := make(map[string]bool)
	// disable "all" line per default
	selected["All"] = false
	for _, cd := range connectionDatas {
		data := populateMissingConnectionsData(allConnections, cd.data)
		line.AddSeries(cd.name, generateLineItems(data))
		selected[cd.name] = true
	}
	line.SetGlobalOptions(charts.WithLegendOpts(opts.Legend{
		Selected: selected,
		// Show:     opts.Bool(false),
	}))
	f, _ := os.Create("public/" + name + ".html")
	line.Render(f)
}

func createTop10KillerChart(chartName string, topKiller []TopKiller) {
	// create a new bar instance
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: chartName,
	}))
	bar.SetGlobalOptions(charts.WithLegendOpts(opts.Legend{
		Show: opts.Bool(false),
	}))

	items := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	for _, killer := range topKiller {
		xAxis = append(xAxis, killer.CharacterName)
		items = append(items, opts.BarData{Value: killer.Count})
	}
	bar.SetXAxis(xAxis)
	bar.AddSeries("Kills", items)
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
	bar.SetGlobalOptions(charts.WithLegendOpts(opts.Legend{
		Show: opts.Bool(false),
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

func createPlayTimePerServer(db *mongo.Database, mongoCtx context.Context, serverList []Server) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "PlayTime per server this wipe",
	}))
	bar.SetGlobalOptions(charts.WithLegendOpts(opts.Legend{
		Show: opts.Bool(false),
	}))
	items := make([]opts.BarData, 0)
	xAxis := make([]string, 0)
	total := 0
	for _, v := range serverList {
		chars := getCharacters(db, mongoCtx, v.ServerId)
		xAxis = append(xAxis, fmt.Sprintf("%s (%s)", v.Name, v.Region))
		for _, v := range chars {
			total += int(v.PlayTime)
		}
		items = append(items, opts.BarData{Value: math.Floor(float64(total) / 60.0)})
	}
	bar.SetXAxis(xAxis)
	bar.AddSeries("Cumulative time in hours", items)
	f, error := os.Create("public/" + "playtime" + ".html")
	if error != nil {
		panic(error)
	}
	bar.Render(f)
}

func genCharts(db *mongo.Database, mongoCtx context.Context) {
	servers := getServers(db, mongoCtx)
	officialServers := []Server{}
	enabledServers := []Server{}
	for _, v := range servers {
		if v.IsOfficial && !v.IsDisabled {
			officialServers = append(officialServers, v)
		}
		if !v.IsDisabled {
			enabledServers = append(enabledServers, v)
		}
	}
	createPlayTimePerServer(db, mongoCtx, officialServers)
	createCountPerServerCharts(db, mongoCtx, officialServers, CHARACTERS_COLLECTION_NAME, "Characters per server")
	createCountPerServerCharts(db, mongoCtx, officialServers, CONSTRUCTIONS_COLLECTION_NAME, "Constructions per server")
	createCountPerServerCharts(db, mongoCtx, officialServers, CROPS_COLLECTION_NAME, "Crops per server")
	top := getTopKiller(db, mongoCtx, 12, "player")
	createTop10KillerChart("Top Killer Main EU 1", top)
	top = getTopKiller(db, mongoCtx, 11, "player")
	createTop10KillerChart("Top Killer Main US 1", top)
	top = getTopKiller(db, mongoCtx, 61, "zombie")
	createTop10KillerChart("Top Zombie Killer Help", top)
	allKills := getAllKills(db, mongoCtx)

	killsDatas := make([]CountData, 0)
	for _, v := range officialServers {
		// TODO: shitty it's to avoid help server but it should ignore pve servers in general
		if v.ServerId == 61 {
			continue
		}
		data := getKillsToServer(db, mongoCtx, v.ServerId)
		if len(data) == 0 {
			continue
		}
		killsDatas = append(killsDatas, CountData{name: v.Name + " " + v.Region, data: data})
	}
	createLineCountChart("Kill activity", allKills, killsDatas)

	connectionsDatas := make([]CountData, 0)
	allConnections := getAllConnections(db, mongoCtx)
	for _, v := range officialServers {
		data := getConnectionsToServer(db, mongoCtx, v.ServerId)
		if len(data) == 0 {
			continue
		}
		connectionsDatas = append(connectionsDatas, CountData{name: v.Name + " " + v.Region, data: data})
	}
	createLineCountChart("connections", allConnections, connectionsDatas)
	lastMonthConnectionsDatas := make([]CountData, 0)
	allConnectionsLastMonth := getAllConnectionsLastMonth(db, mongoCtx)
	for _, v := range officialServers {
		data := getConnectionsLastMonthToServer(db, mongoCtx, v.ServerId)
		if len(data) == 0 {
			continue
		}
		lastMonthConnectionsDatas = append(lastMonthConnectionsDatas, CountData{name: v.Name + " " + v.Region, data: data})
	}
	createLineCountChart("Last month connections", allConnectionsLastMonth, lastMonthConnectionsDatas)
}
