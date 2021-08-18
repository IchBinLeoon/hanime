package utils

import "fmt"

type Bar struct {
	percent int64
	current int64
	total   int64
	rate    string
	graph   string
}

func NewProgressBar(start int64, total int64, graph string) *Bar {
	bar := &Bar{
		current: start,
		total:   total,
		graph:   graph,
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph
	}
	return bar
}

func (bar *Bar) Add(value int64) {
	bar.current += value
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last {
		var i int64 = 0
		for ; i < bar.percent-last; i++ {
			bar.rate += bar.graph
		}
		fmt.Printf("\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent*2, bar.current, bar.total)
	}
}

func (bar *Bar) Finish() {
	fmt.Print("\n")
}

func (bar *Bar) getPercent() int64 {
	return int64((float32(bar.current) / float32(bar.total)) * 50)
}
