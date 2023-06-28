package contract

import (
	"fmt"
)

// Bar ...
type Bar struct {
	percent int64  // progress percentage
	cur     int64  // current progress
	total   int64  // total value for progress
	rate    string // the actual progress bar to be printed
	graph   string // the fill value for progress bar
	indent  string //Indentation
}

func (bar *Bar) NewOption(start, total int64, in string) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "*"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph // initial progress position
	}
	bar.indent = in
}

func (bar *Bar) getPercent() int64 {
	return int64((float32(bar.cur) / float32(bar.total)) * 100)
}

func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	n := last
	for n < bar.percent {
		bar.rate += bar.graph
		n += 2
	}
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r%s[%-50s]%3d%% %8d/%d", bar.indent, bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Finish() {
	fmt.Println()
}
