package main

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"

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
	// create a new bar instance
	bar := charts.NewLine()

	// Put data into instance
	bar.SetXAxis(getXaxis(allConnections))
	bar.AddSeries("all", generateLineItems(allConnections))
	for _, cd := range connectionDatas {
		data := populateMissingConnectionsData(allConnections, cd.data)
		bar.AddSeries(cd.name, generateLineItems(data))
	}
	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)
}
