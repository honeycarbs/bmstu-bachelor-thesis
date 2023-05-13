package heatmap

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"ml/internal/entity"
	"ml/internal/pkg/cm"
	"ml/pkg/logging"
	"os"
)

var (
	labels = []entity.Label{entity.NEUTRAL, entity.SAD, entity.POSITIVE, entity.ANGRY}
)

func genHeatMapData(actual, predicted []entity.Label) []opts.HeatMapData {
	items := make([]opts.HeatMapData, 0)

	//confusionMatrix := make(map[entity.Label]map[entity.Label]int)
	//for _, label := range labels {
	//	confusionMatrix[label] = make(map[entity.Label]int)
	//	for _, label2 := range labels {
	//		confusionMatrix[label][label2] = 0
	//	}
	//}
	//
	//for i := 0; i < len(actual); i++ {
	//	confusionMatrix[actual[i]][predicted[i]]++
	//}
	confusionMatrix := cm.NewConfusionMatrix(labels, actual, predicted)
	logging.GetLogger().Info(confusionMatrix.Values)
	confusionMatrix.Normalize()
	logging.GetLogger().Info(confusionMatrix.Values)

	for _, label1 := range labels {
		for _, label2 := range labels {
			if confusionMatrix.Values[label1][label2] == 0 {
				items = append(items, opts.HeatMapData{Value: [3]interface{}{label1, label2, "-"}})
			} else {
				items = append(items, opts.HeatMapData{Value: [3]interface{}{label1, label2, confusionMatrix.Values[label1][label2]}})
			}
		}
	}

	return items
}

func heatMapBase(actual, predicted []entity.Label) *charts.HeatMap {
	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Матрица ошибок",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      labels,
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Show:       true,
			//Max:        1,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#F0D99C", "#BF444C"},
			},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: false,
		}),
	)

	//fmt.Println(len(actual), float32(len(predicted)))
	hm.SetXAxis(labels).AddSeries("heatmap", genHeatMapData(actual, predicted))
	return hm
}

func GetHeatmap(actual, predicted []entity.Label) {
	page := components.NewPage()
	page.AddCharts(
		heatMapBase(actual, predicted),
	)

	f, err := os.Create("etc/static/heatmap.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
